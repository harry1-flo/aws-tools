package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ldg940804-aws-tools/core/aws"
	"github.com/ldg940804-aws-tools/fs"
)

// var acs = []string{"dev", "dreamusaws001", "data", "deploy", "prd", "sec"}
var acs = "prd"

// asg의 설정 따오기
// https://docs.google.com/spreadsheets/d/168ZK0pfKtR5-PqNO7ZDXpMK3gQaNqmtzNibzPV4zqz4/edit?gid=2080214049#gid=2080214049
func main() {
	cfg, account := aws.NewAWSUseProfile(acs)
	ec2 := aws.NewEC2(*cfg, account)

	csv := fs.NewCSV("ec2")
	defer csv.End()

	// EC2 리스트 추출
	// fmt.Println("EC2 리스트 추출 중 ...")
	// ec2List := ec2.ListInstance(true)
	// txt, _ := os.OpenFile("./list.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// for _, ec2 := range ec2List {

	// 	if ec2.AutoScalingGroupName == "" {
	// 		continue
	// 	}

	// 	if strings.Contains(ec2.AutoScalingGroupName, "eks") {
	// 		continue
	// 	}

	// 	if strings.Contains(ec2.AutoScalingGroupName, "alpha") {
	// 		continue
	// 	}

	// 	txt.WriteString(fmt.Sprintf("%s,%s\n", ec2.Name, ec2.AutoScalingGroupName))
	// }

	// ASG 설정 추출
	fmt.Println("ASG List 추출 중...")
	pwd, _ := os.Getwd()
	b, _ := os.ReadFile(pwd + "/list.txt")
	for _, line := range strings.Split(string(b), "\n") {
		split := strings.Split(line, ",")

		if len(split) != 2 {
			fmt.Println("split error : ", line)
			continue
		}

		ec2Name := split[0]
		asgName := split[1]

		list := ec2.GetEC2DetailsByInstanceId(asgName)

		schedules := []string{}
		for key, schedule := range list.Schdules {
			schedules = append(schedules, fmt.Sprintf("%s : %s", key, schedule))
		}

		/*
			서비스 영향도
			서비스
			ASG 이름
			타입
			주간 최소대수
			주간 최대개수
			스케쥴링 관련정책
			야간 최소대수
			야간 최대대수
			scale-out 기준
			scale-out 대수
			scale-int 기준
			scale-int 대수
		*/

		// fmt.Println(ec2Name, asgName, list.InstanceType, list.MinSize, list.MaxSize, strings.Join(schedules, " || "), list.NightMinSize, list.NightMaxSize, list.ScaleOutCondition, list.ScaleOutConditionValue, list.ScaleInConditionValue, list.ScaleInCondition)
		// os.Exit(0)

		csv.OneFileWrite(
			serviceImportance[ec2Name],
			ec2Name,
			asgName,
			list.InstanceType,
			list.MinSize,
			list.MaxSize,
			strings.Join(schedules, " || "),
			list.NightMinSize,
			list.NightMaxSize,
			list.ScaleOutCondition,
			list.ScaleOutConditionValue,
			list.ScaleInCondition,
			list.ScaleInConditionValue)
	}
}

// 2025.12.05 기준
var serviceImportance = map[string]string{
	"apigw.prod.music-flo.com":                          "높음",
	"chartandchannel.prod.music-flo.com":                "높음",
	"curation.prod.music-flo.com":                       "높음",
	"display.prod.music-flo.com":                        "높음",
	"external.prod.music-flo.com":                       "높음",
	"floadmin.prod.music-flo.com":                       "높음",
	"general.prod.music-flo.com":                        "높음",
	"history-async-api.prod.music-flo.com":              "높음",
	"mcp-key.prod.music-flo.com":                        "높음",
	"mcp-stream.prod.music-flo.com":                     "높음",
	"member-command.prod.music-flo.com":                 "높음",
	"member-query.prod.music-flo.com":                   "높음",
	"meta-mgo.prod.music-flo.com":                       "높음",
	"partner-apigw.prod.music-flo.com":                  "높음",
	"partner-coordinator.prod.music-flo.com":            "높음",
	"party-apigw.prod.music-flo.com":                    "높음",
	"party.prod.music-flo.com":                          "높음",
	"payment.prod.music-flo.com":                        "높음",
	"personal-mgo.prod.music-flo.com":                   "높음",
	"play.prod.music-flo.com":                           "높음",
	"product.prod.music-flo.com":                        "높음",
	"purchase.prod.music-flo.com":                       "높음",
	"search.prod.music-flo.com":                         "높음",
	"spring-config.prod.music-flo.com":                  "높음",
	"stream.prod.music-flo.com":                         "높음",
	"support.prod.music-flo.com":                        "높음",
	"www3.prod.music-flo.com":                           "높음",
	"m.prod.music-flo.com":                              "높음",
	"www.prod.music-flo.com":                            "높음",
	"cast-receiver.prod.music-flo.com":                  "높음",
	"digital-card.prod.music-flo.com":                   "높음",
	"mds.dreamuscompany.com":                            "높음",
	"[PROD]-[EC2]-ocr-prd":                              "높음",
	"outif-api.prod.music-flo.com":                      "높음",
	"picka.prod.music-flo.com":                          "높음",
	"playwith.prod.music-flo.com":                       "높음",
	"pp-external.prod.music-flo.com":                    "높음",
	"pp-party.prod.music-flo.com":                       "높음",
	"pri-mds.dreamuscompany.com":                        "높음",
	"reco-1seed.prod.music-flo.com":                     "높음",
	"reco.api.music-flo.io":                             "높음",
	"reco-model.api.music-flo.io":                       "높음",
	"settlement-api-v2.prod.music-flo.com":              "높음",
	"settlement-api-v3.prod.music-flo.com":              "높음",
	"share.prod.music-flo.com":                          "높음",
	"voucher-inbound.prod.music-flo.com":                "높음",
	"voucher-worker.prod.music-flo.com":                 "높음",
	"collection.prod.log.infra.music-flo.io":            "보통",
	"ingestion.prod.log.infra.music-flo.io":             "보통",
	"reco-kafka-client.prod.music-flo.com":              "보통",
	"reco-targeting-inventory.music-flo.com":            "보통",
	"vue3-cmsflo.prod.music-flo.com":                    "보통",
	"admin-kms.music-flo.com":                           "낮음",
	"cmsflo.prod.music-flo.com":                         "낮음",
	"kms-key.music-flo.com":                             "낮음",
	"mcp-api.prod.music-flo.com":                        "낮음",
	"mcp-approver.prod.music-flo.com":                   "낮음",
	"mcp-archive.prod.music-flo.com":                    "낮음",
	"mcp-cdn-cleanser.prod.music-flo.com":               "낮음",
	"mcp-cms.prod.music-flo.com":                        "낮음",
	"mcp-di.prod.music-flo.com":                         "낮음",
	"mcp-downloader.prod.music-flo.com":                 "낮음",
	"mcp-heartqueen.prod.music-flo.com":                 "낮음",
	"mcp-logstash.prod.music-flo.com":                   "낮음",
	"mcp-mds-ddex-generator.prod.music-flo.com":         "낮음",
	"mcp-mds-event-listener.prod.music-flo.com":         "낮음",
	"mcp-mds-platform-sender.prod.music-flo.com":        "낮음",
	"mcp-mds.prod.music-flo.com":                        "낮음",
	"mcp-mediaconvert-event-handler.prod.music-flo.com": "낮음",
	"mcp-mediainfo-api.prod.music-flo.com":              "낮음",
	"mcp-meta.prod.music-flo.com":                       "낮음",
	"mcp-observer.prod.music-flo.com":                   "낮음",
	"mcp-parser.prod.music-flo.com":                     "낮음",
	"mcp-vcoloring.prod.music-flo.com":                  "낮음",
	"mcp-video-transcoder-v2.prod.music-flo.com":        "낮음",
	"mcp3.prod.music-flo.com":                           "낮음",
}
