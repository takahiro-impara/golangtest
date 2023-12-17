package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type MockAPI struct {
	Output *ec2.DescribeInstancesOutput
	Error  error
}

func (m *MockAPI) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	// Implement your mock behavior here
	return m.Output, m.Error
}

func TestGetPublicHostname(t *testing.T) {
	// モック用の値
	var (
		instanceId    = "i-1234567890"
		publicDNSName = "public.example.com"
		publicIP      = "1.2.3.4"
	)

	// ケース1 : PublicDnsName, PublicIpAddress ともに存在するケース
	var mock = &MockAPI{
		Output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{{Instances: []types.Instance{{InstanceId: &instanceId, PublicDnsName: &publicDNSName, PublicIpAddress: &publicIP}}}},
		},
		Error: nil,
	}
	var result, err = GetPublicHostname(mock, instanceId)
	if err != nil {
		t.Errorf("GetPublicHostName関数でエラーが発生 (%v)", err.Error())
	}
	if result != publicDNSName {
		t.Errorf("戻り値が%vでない", publicDNSName)
	}

	// ケース2 : PublicIpAddress だけ存在するケース
	mock = &MockAPI{
		Output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{{Instances: []types.Instance{{InstanceId: &instanceId, PublicIpAddress: &publicIP}}}},
		},
		Error: nil,
	}
	// APIをmockに差し替えて実行 : エラー無く "1.2.3.4" が戻り値になる
	result, err = GetPublicHostname(mock, instanceId)
	if err != nil {
		t.Errorf("GetPublicHostName関数でエラーが発生 (%v)", err.Error())
	}
	if result != publicIP {
		t.Errorf("戻り値が%vでない", publicIP)
	}

	// ケース3 : PublicDnsName, PublicIpAddress どちらも無いケース (インスタンス停止中に相当)
	mock = &MockAPI{
		Output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{{Instances: []types.Instance{{InstanceId: &instanceId}}}},
		},
		Error: nil,
	}
	// APIをmockに差し替えて実行 : エラーになる
	_, err = GetPublicHostname(mock, instanceId)
	if err == nil {
		t.Error("GetPublicHostName関数でエラーが発生しない")
	}
}
