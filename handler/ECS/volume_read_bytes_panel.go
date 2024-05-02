package ECS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type ReadBytes struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"RawData"`
// }

var AwsxECSReadBytesCmd = &cobra.Command{
	Use:   "volume_readbytes_panel",
	Short: "get volume read bytes metrics data",
	Long:  `command to get volume read bytes metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetECSReadBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting volume read bytes metrics data: ", err)
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

func GetECSReadBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	startTime, endTime, err := commanFunction.ParseTimes(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}
	instanceId, err = commanFunction.GetCmdbData(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ECS/ContainerInsights", "StorageReadBytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting volume read bytes data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Volume_read_bytes"] = rawData
	return "", cloudwatchMetricData, nil

}

// func processECSReadBytesRawdata(result *cloudwatch.GetMetricDataOutput) ReadBytes {
// 	var rawData ReadBytes
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
	AwsxECSReadBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSReadBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSReadBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxECSReadBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSReadBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSReadBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSReadBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSReadBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSReadBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSReadBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSReadBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSReadBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSReadBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSReadBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSReadBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSReadBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
