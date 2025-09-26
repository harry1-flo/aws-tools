package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Config struct {
	Account string

	ec2Client *ec2.Client
}

func NewEC2(cfg aws.Config, account string) *EC2Config {
	return &EC2Config{
		ec2Client: ec2.NewFromConfig(cfg),
		Account:   account,
	}
}

type ListInstanceParmas struct {
	InstanceId string
	Name       string
	Tags       map[string]string
}

func (c EC2Config) ListInstance() map[string]ListInstanceParmas {

	resp, err := c.ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil
	}

	fmt.Println("ec2 count : ", len(resp.Reservations))

	instances := make(map[string]ListInstanceParmas)
	for _, reserve := range resp.Reservations {
		for _, instance := range reserve.Instances {

			tagMap := spreadTags(instance.Tags)

			instances[*instance.InstanceId] = ListInstanceParmas{
				InstanceId: *instance.InstanceId,
				Name:       tagMap["Name"],
				Tags:       tagMap,
			}

		}
	}

	return instances
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
