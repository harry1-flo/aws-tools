package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Config struct {
	Account string

	ec2Client         *ec2.Client
	autoScalingClient *autoscaling.Client
}

func NewEC2(cfg aws.Config, account string) *EC2Config {
	return &EC2Config{
		Account:           NamingConvert[account],
		ec2Client:         ec2.NewFromConfig(cfg),
		autoScalingClient: autoscaling.NewFromConfig(cfg),
	}
}

type ListInstanceParmas struct {
	InstanceId   string
	Name         string
	InstanceType types.InstanceType

	IsAutoScaling bool

	VolumeSize int
	Tags       map[string]string

	// Autoscaling
	AsgMin             string
	AsgMax             string
	SchedulingMinCount string // 예약 스케줄링 최소 개수

	// Network
	PublicIp    string
	PrivateIp   string
	VpcName     string
	SubnetNames string

	// Spread Tag
	Team        string
	Environment string
	Backup      string

	// Instance Count (AutoScaling 그룹에 아닌 단일인스턴스로 떠있는 것들)
	SingleInstanceCount int

	// Status
	InstanceState types.InstanceStateName
}

func (c EC2Config) ListInstance(isAutoScaling bool) map[string]ListInstanceParmas {

	resp, err := c.ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		panic(err)
	}

	fmt.Println("ec2 count : ", len(resp.Reservations))

	instances := make(map[string]ListInstanceParmas)

	for _, reserve := range resp.Reservations {
		for _, instance := range reserve.Instances {
			tagMap := spreadTags(instance.Tags)

			// autoscaling group

			// ebs details
			ec2Detail, err := c.GetEC2Details(*instance.InstanceId)
			if err != nil {
				fmt.Println("ec2Detail Error : ", err)
				continue
			}

			// 이미 있다면 개수 증가
			if _, exists := instances[tagMap["Name"]]; exists {
				temp := instances[tagMap["Name"]]
				temp.SingleInstanceCount++
				instances[tagMap["Name"]] = temp
				continue
			}

			var instanceParmas ListInstanceParmas = ListInstanceParmas{}

			instanceParmas = ListInstanceParmas{
				Name:          tagMap["Name"],
				InstanceId:    *instance.InstanceId,
				Team:          tagMap["Team"],
				Environment:   tagMap["Environment"],
				InstanceType:  instance.InstanceType,
				IsAutoScaling: c.isAutoscaling(tagMap),
				Backup:        tagMap["Backup"],
				PublicIp:      aws.ToString(instance.PublicIpAddress),
				PrivateIp:     aws.ToString(instance.PrivateIpAddress),
				VpcName:       aws.ToString(instance.VpcId),
				SubnetNames:   aws.ToString(instance.SubnetId),
				InstanceState: instance.State.Name,

				SingleInstanceCount: 1,

				Tags:       tagMap,
				VolumeSize: ec2Detail.VolumeSize,
			}

			// AutoScaling 표기
			if isAutoScaling && c.isAutoscaling(tagMap) {
				autoScalingGroupName := tagMap["aws:autoscaling:groupName"]

				resp, err := c.autoScalingClient.DescribeAutoScalingGroups(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []string{autoScalingGroupName},
				})
				if err != nil {
					fmt.Println("autoScalingGroup Error : ", err)
					continue
				}

				if len(resp.AutoScalingGroups) > 0 {
					min := *resp.AutoScalingGroups[0].MinSize
					max := *resp.AutoScalingGroups[0].MaxSize

					instanceParmas.AsgMin = strconv.Itoa(int(min))
					instanceParmas.AsgMax = strconv.Itoa(int(max))
				}

				// 예약된 작업에서 최소 개수 조회
				scheduledActionsResp, err := c.autoScalingClient.DescribeScheduledActions(context.Background(), &autoscaling.DescribeScheduledActionsInput{
					AutoScalingGroupName: aws.String(autoScalingGroupName),
				})

				if err != nil {
					fmt.Println("DescribeScheduledActions Error : ", err)
				} else if len(scheduledActionsResp.ScheduledUpdateGroupActions) > 0 {
					// 예약된 작업 중 가장 작은 MinSize 찾기
					var minScheduledSize *int32
					for _, action := range scheduledActionsResp.ScheduledUpdateGroupActions {
						if action.MinSize != nil {
							if minScheduledSize == nil || *action.MinSize < *minScheduledSize {
								minScheduledSize = action.MinSize
							}
						}
					}

					if minScheduledSize != nil {
						instanceParmas.SchedulingMinCount = strconv.Itoa(int(*minScheduledSize))
					}
				}

			}

			instances[tagMap["Name"]] = instanceParmas

		}
	}

	return instances
}

type EC2DetailParmas struct {
	VolumeSize int
}

func (c EC2Config) GetEC2Details(id string) (EC2DetailParmas, error) {

	resp, err := c.ec2Client.DescribeVolumes(context.Background(), &ec2.DescribeVolumesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: []string{id},
			},
		},
	})

	if err != nil {
		return EC2DetailParmas{}, err
	}

	volumeSize := 0
	for _, ebs := range resp.Volumes {
		volumeSize += int(*ebs.Size)
	}

	return EC2DetailParmas{
		VolumeSize: volumeSize,
	}, nil
}

func (c EC2Config) GetAccount() string {
	return c.Account
}

func spreadTags(tags []types.Tag) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[*tag.Key] = *tag.Value
	}
	return tagMap
}

func (c EC2Config) isAutoscaling(tag map[string]string) bool {
	if tag["aws:autoscaling:groupName"] != "" {
		return true
	}
	return false
}
