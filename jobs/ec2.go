package jobs

import (
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
			csv.Write(id, item.Tags["Name"], item.Tags["Service"], item.Tags["Environment"], item.Tags["aws:autoscaling:groupName"], item.Tags["Backup"])
		}
	}
}
