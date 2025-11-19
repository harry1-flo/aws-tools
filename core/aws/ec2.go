package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	AutoScalingGroupName string
	AsgMin               string
	AsgMax               string
	SchedulingMinCount   string // 예약 스케줄링 최소 개수

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
				instanceParmas.AutoScalingGroupName = autoScalingGroupName

				// 동적 크기 정책 존재 여부
				policiesResp, err := c.autoScalingClient.DescribePolicies(context.Background(), &autoscaling.DescribePoliciesInput{
					AutoScalingGroupName: aws.String(autoScalingGroupName),
				})

				if err != nil {
					fmt.Println("DescribePolicies Error : ", err)
					break
				}

				if len(policiesResp.ScalingPolicies) > 0 {

					resp, err := c.autoScalingClient.DescribeAutoScalingGroups(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
						AutoScalingGroupNames: []string{autoScalingGroupName},
					})

					if err != nil {
						fmt.Println("autoScalingGroup Error : ", err)
						break
					}

					// AutoScalin Group 조회
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

type ASGType struct {
	InstanceType string

	MinSize string
	MaxSize string

	Schdules map[string]string

	// Schdule 걸려있을 경우
	NightMinSize string
	NightMaxSize string

	ScaleOutCondition      string
	ScaleOutConditionValue string

	ScaleInCondition      string
	ScaleInConditionValue string
}

func (c EC2Config) GetEC2DetailsByInstanceId(asgName string) ASGType {

	asgParams := ASGType{}

	output, err := c.autoScalingClient.DescribeAutoScalingGroups(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})

	if err != nil {
		panic(err)
	}

	asg := output.AutoScalingGroups[0]

	// Instance Type
	var instanceType string

	// 현재 실행 중인 인스턴스에서 타입 가져오기
	if len(asg.Instances) > 0 && asg.Instances[0].InstanceId != nil {
		ec2Resp, err := c.ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
			InstanceIds: []string{*asg.Instances[0].InstanceId},
		})
		if err == nil && len(ec2Resp.Reservations) > 0 && len(ec2Resp.Reservations[0].Instances) > 0 {
			instanceType = string(ec2Resp.Reservations[0].Instances[0].InstanceType)
		}
	}

	// 인스턴스가 없으면 Launch Template에서 가져오기
	if instanceType == "" {
		var ltName *string
		var ltVersion *string

		if asg.MixedInstancesPolicy != nil && asg.MixedInstancesPolicy.LaunchTemplate != nil {
			ltName = asg.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.LaunchTemplateName
			ltVersion = asg.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.Version
		} else if asg.LaunchTemplate != nil {
			ltName = asg.LaunchTemplate.LaunchTemplateName
			ltVersion = asg.LaunchTemplate.Version
		}

		if ltName != nil {
			launchTemplateResp, err := c.ec2Client.DescribeLaunchTemplateVersions(context.Background(), &ec2.DescribeLaunchTemplateVersionsInput{
				LaunchTemplateName: ltName,
				Versions:           []string{aws.ToString(ltVersion)},
			})
			if err == nil && len(launchTemplateResp.LaunchTemplateVersions) > 0 {
				ltInstanceType := launchTemplateResp.LaunchTemplateVersions[0].LaunchTemplateData.InstanceType
				if ltInstanceType != "" {
					instanceType = string(ltInstanceType)
				}
			}
		}
	}

	asgParams.InstanceType = instanceType

	asgParams.MinSize = strconv.Itoa(int(*asg.MinSize))
	asgParams.MaxSize = strconv.Itoa(int(*asg.MaxSize))

	asgParams.NightMinSize = strconv.Itoa(int(*asg.MaxSize))

	scheduledActionsResp, err := c.autoScalingClient.DescribeScheduledActions(context.Background(), &autoscaling.DescribeScheduledActionsInput{
		AutoScalingGroupName: aws.String(*asg.AutoScalingGroupName),
	})
	if err != nil {
		panic(err)
	}

	// 스케쥴 정책 있다면
	if len(scheduledActionsResp.ScheduledUpdateGroupActions) > 0 {
		shce := make(map[string]string)

		for _, schedule := range scheduledActionsResp.ScheduledUpdateGroupActions {

			shce[*schedule.ScheduledActionName] = *schedule.Recurrence

			if strings.Contains(*schedule.ScheduledActionName, "night") {
				asgParams.NightMinSize = strconv.Itoa(int(*schedule.MinSize))

				if schedule.MaxSize != nil {
					asgParams.NightMaxSize = strconv.Itoa(int(*schedule.MaxSize))
				}
			}
		}

		asgParams.Schdules = shce
	}

	// 스케쥴링
	policyResp, err := c.autoScalingClient.DescribePolicies(context.Background(), &autoscaling.DescribePoliciesInput{
		AutoScalingGroupName: aws.String(*asg.AutoScalingGroupName),
	})
	if err != nil {
		panic(err)
	}

	for _, policy := range policyResp.ScalingPolicies {
		fmt.Println(*policy.PolicyName)

		if strings.Contains(*policy.PolicyName, "cpu-out") {
			asgParams.ScaleOutCondition = ""
			asgParams.ScaleOutConditionValue = strconv.Itoa(int(*policy.ScalingAdjustment))
		}

		if strings.Contains(*policy.PolicyName, "cpu-in") {
			asgParams.ScaleInCondition = ""
			asgParams.ScaleInConditionValue = strconv.Itoa(int(*policy.ScalingAdjustment))
		}
	}

	return asgParams

}
