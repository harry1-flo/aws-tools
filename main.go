package main

import (
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

func main() {
	cfg, account := aws.NewAWSUseProfile("dev")
	// ec2 := aws.NewEC2(*cfg, account)
	route53 := aws.NewRoute53Config(*cfg, account)

	// jobs.ListEC2NotAutoScaling(*ec2)
	jobs.ListingDNSRecord(*route53)

}
