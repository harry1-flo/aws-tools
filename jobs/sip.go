package jobs

import (
	"fmt"
	"strings"

	"github.com/ldg940804-aws-tools/core/aws"
	"github.com/ldg940804-aws-tools/fs"
)

/*
2025.9.26
AutoScaling이 없는 단독 EC2 인스턴스 리스트
AWS Backup 정책을 활성화 해야 함 (있는것과 없는것 확인해야 함)

Backup plan
- oa
- flo-deploy
- production
*/
func ListEC2NotAutoScaling(ec2 aws.EC2Config) {

	ec2List := ec2.ListInstance()
	csv := fs.NewCSV(ec2.Account)
	defer csv.End()

	csv.Write("InstanceId", "Name", "Service", "Environment", "AutoScalingGroup", "Backup")
	for id, item := range ec2List {

		if item.Tags["aws:autoscaling:groupName"] == "" {
			csv.OneFileWrite(ec2.Account, id, item.Tags["Name"], item.Tags["Service"], item.Tags["Environment"], item.Tags["aws:autoscaling:groupName"], item.Tags["Backup"])
		}
	}
}

func ListEC2AutoScalingButNotBackupTag(ec2 aws.EC2Config) {

	ec2List := ec2.ListInstance()

	csv := fs.NewCSV(fmt.Sprintf("%s-%s", ec2.Account, "ec2-autoscaling-not-backuptag"))
	defer csv.End()

	csv.Write("InstanceId", "Name", "Service", "Environment", "AutoScalingGroup", "Backup")
	for id, item := range ec2List {
		fmt.Println(id, item.Tags["Backup"])
		if item.Tags["aws:autoscaling:groupName"] != "" {
			csv.Write(id, item.Tags["Name"], item.Tags["Service"], item.Tags["Environment"], item.Tags["aws:autoscaling:groupName"], item.Tags["Backup"])
		}
	}
}

/*
2025.9.26
DNS Record에 Alias type 중 AWS ID가 명시되지 않는 레코드 값이 유효한지 확인
*/
func ListingDNSRecord(route53 aws.Route53Config) {
	csv := fs.NewCSV(fmt.Sprintf("%s-%s", route53.Account, "dns"))
	defer csv.End()

	hostingsMap, err := route53.ListHostingLayer()
	if err != nil {
		panic(err)
	}

	csv.Write("AccountId", "HostedZoneId", "Name", "Type", "Value")
	for hostId, hostName := range hostingsMap {

		record, err := route53.ListDNSRecord(hostId)
		if err != nil {
			panic(err)
		}

		for _, r := range record {

			// acm 일단 제외, 근데 안쓰는 acm이 있을 수도 있음...
			if r.Name == "" || r.Type == "SOA" || strings.Contains(r.Value, "acm-validations") || r.Type == "SRV" {
				continue
			}

			csv.OneFileWrite(route53.Account, hostName, r.Name, r.Type, r.Value)
		}
	}

}
