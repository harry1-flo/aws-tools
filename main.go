package main

import (
	"fmt"

	"github.com/ldg940804-aws-tools/core/aws"
	"github.com/ldg940804-aws-tools/jobs"
)

/*
	var DreamsuProfileARNS = map[string]string{
	"dev":           "972467631093",
	"dreamusaws001": "294438099013",
	"data":          "822479450351",
	"deploy":        "803135791119",
	"prd":           "865306278000",
	"sec":           "953265915023",
}
*/

var acs = []string{"dev", "dreamusaws001", "data", "deploy", "prd", "sec"}

// var acs = []string{"deploy"}

func main() {

	for _, ac := range acs {
		fmt.Println("account : ", ac)
		cfg, account := aws.NewAWSUseProfile(ac)

		// resources
		ec2 := aws.NewEC2(*cfg, account)
		// route53 := aws.NewRoute53Config(*cfg, account)

		// jobs
		// jobs.ListEC2NotAutoScaling(*ec2)
		jobs.ListEC2(*ec2)
		// jobs.ListEC2AutoScalingButNotBackupTag(*ec2)
		// jobs.ListingDNSRecord(*route53)

		// jobs.ISMSEC2(*ec2)
	}
}
