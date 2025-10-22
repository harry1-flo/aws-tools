package main

import (
	"fmt"

	"github.com/ldg940804-aws-tools/core/aws"
	"github.com/ldg940804-aws-tools/jobs"
)

var acs = []string{"dev", "dreamusaws001", "data", "deploy", "prd", "sec"}

func main() {

	for _, ac := range acs {
		fmt.Println("account : ", ac)
		cfg, account := aws.NewAWSUseProfile(ac)

		// resources
		ec2 := aws.NewEC2(*cfg, account)
		jobs.ListEC2(*ec2)
	}
}
