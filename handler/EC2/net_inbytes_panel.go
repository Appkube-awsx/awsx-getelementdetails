package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type NetworkInBytes struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Net_Inbytes"`
// }

var AwsxEc2NetworkInBytesCmd = &cobra.Command{
	Use:   "network_inbytes_utilization_panel",
	Short: "get network inbytes metrics data",
	Long:  `command to get network inbytes metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")
		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, cloudwatchMetricResp, err := GetNetworkInBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network inbytes metrics data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetNetworkInBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	startTime, endTime, err := comman_function.ParseTimes(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}
	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/EC2", "NetworkIn", startTime, endTime, "Sum", "InstanceId", cloudWatchClient)

	if err != nil {
		log.Println("Error in getting network inbytes data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Net_Inbytes"] = rawData
	return "", cloudwatchMetricData, nil
}

// func processInbytesRawdata(result *cloudwatch.GetMetricDataOutput) NetworkInBytes {
// 	var rawData NetworkInBytes
// 	rawData.RawData = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.RawData[i].Timestamp = *timestamp
// 		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
// 	}

// 	return rawData
// }

func init() {
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2NetworkInBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
