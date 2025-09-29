package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

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
