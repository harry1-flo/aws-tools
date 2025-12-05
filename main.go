package main

// var acs = []string{"dev", "dreamusaws001", "data", "deploy", "prd", "sec"}
var acs = []string{"prd"}

// func main() {

// 	csv := fs.NewCSV("ec2")
// 	defer csv.End()

// 	for _, ac := range acs {
// 		fmt.Println("account : ", ac)
// 		cfg, account := aws.NewAWSUseProfile(ac)

// 		// resources
// 		ec2 := aws.NewEC2(*cfg, account)

// 		asgList := ec2.ListInstance(true)

// 		for ec2Name, attr := range asgList {

// 			if !strings.Contains(ec2Name, "alpha") {
// 				csv.OneFileWrite(ec2Name, attr.AutoScalingGroupName, string(attr.InstanceType), attr.AsgMin, attr.AsgMax)
// 			}
// 		}
// 	}
// }

// asg의 설정 따오기
// https://docs.google.com/spreadsheets/d/168ZK0pfKtR5-PqNO7ZDXpMK3gQaNqmtzNibzPV4zqz4/edit?gid=2080214049#gid=2080214049
func main() {

}
