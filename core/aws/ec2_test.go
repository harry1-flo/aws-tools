package aws

import (
	"fmt"
	"testing"
)

func Test_ListInstance(t *testing.T) {
	cfg, account := NewAWSUseProfile("prd")
	ec2 := NewEC2(*cfg, account)

	instancId := "i-02cbfdfcb99557859"

	volumeDetail, _ := ec2.GetEC2Details(instancId)
	fmt.Println(volumeDetail)

}
