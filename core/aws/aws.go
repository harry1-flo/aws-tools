package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var DreamsuProfileARNS = map[string]string{
	"dev":           "972467631093",
	"dreamusaws001": "294438099013",
	"data":          "822479450351",
	"deploy":        "803135791119",
	"prd":           "865306278000",
	"sec":           "953265915023",
}

var NamingConvert = map[string]string{
	"dev":           "dev-flo",
	"dreamusaws001": "dreamusaws001",
	"data":          "flo-data",
	"deploy":        "flo-deploy",
	"prd":           "flo-production",
	"sec":           "flo-security",
}

// useProfile
func NewAWSUseProfile(profile string) (*aws.Config, string) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		panic(err)
	}

	return &aws.Config{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	}, profile
}
