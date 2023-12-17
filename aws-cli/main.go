package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2API interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func NewAPI() (EC2API, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("sandbox"),
		config.WithRegion("ap-northeast-1"),
		config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.TokenProvider = func() (string, error) {
				return stscreds.StdinTokenProvider()
			}
		}),
	)

	if err != nil {
		return nil, err
	}

	return ec2.NewFromConfig(cfg), nil
}

func GetPublicHostname(api EC2API, instanceId string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}

	result, err := api.DescribeInstances(context.TODO(), input)
	if err != nil {
		return "", err
	}

	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			hostname := i.PublicDnsName
			if hostname != nil && *hostname != "" {
				return *hostname, nil
			}
			publicIP := i.PublicIpAddress
			if publicIP != nil && *publicIP != "" {
				return *publicIP, nil
			}
		}
	}
	return "", errors.New("no public DNS name or public IP address found")
}

func main() {
	api, err := NewAPI()
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		os.Exit(1)
	}

	myInstanceId := "i-01c54db43645eadad"
	res, err := GetPublicHostname(api, myInstanceId)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("public hostname: %v\n", res)
}
