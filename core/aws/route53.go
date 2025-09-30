package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

type Route53Config struct {
	Account       string
	route53Client *route53.Client
}

func NewRoute53Config(cfg aws.Config, account string) *Route53Config {
	return &Route53Config{
		route53Client: route53.NewFromConfig(cfg),
		Account:       NamingConvert[account],
	}
}

/*
@response hostedId : hostName
*/
func (r *Route53Config) ListHostingLayer() (map[string]string, error) {
	resp, err := r.route53Client.ListHostedZones(context.Background(), &route53.ListHostedZonesInput{})

	if err != nil {
		return nil, err
	}

	hostingLayer := make(map[string]string)
	for _, hz := range resp.HostedZones {
		hostingLayer[*hz.Id] = *hz.Name
	}

	return hostingLayer, nil
}

type RecordParams struct {
	Name  string
	Type  string
	Value string
}

func (r *Route53Config) ListDNSRecord(hostedId string) ([]RecordParams, error) {
	resp, err := r.route53Client.ListResourceRecordSets(context.Background(), &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedId),
	})
	if err != nil {
		return nil, err
	}

	recordParams := make([]RecordParams, 0, len(resp.ResourceRecordSets))
	for _, re := range resp.ResourceRecordSets {
		// ResourceRecords가 비어있지 않은 경우만 처리
		if len(re.ResourceRecords) > 0 {

			// Alias 일 경우
			if re.AliasTarget != nil {
				recordParams = append(recordParams, RecordParams{
					Name:  *re.Name,
					Type:  string(re.Type),
					Value: *re.AliasTarget.DNSName,
				})
			} else {
				recordParams = append(recordParams, RecordParams{
					Name:  *re.Name,
					Type:  string(re.Type),
					Value: *re.ResourceRecords[0].Value,
				})
			}
		}

	}

	return recordParams, err
}
