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

// type CPUReservedResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU_Reservation"`
// }

var AwsxCpuReservedCmd = &cobra.Command{
	Use:   "cpu_reserved_panel",
	Short: "get cpu reserved metrics data",
	Long:  `command to get cpu reserved metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUReservationData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu reserved data : ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetCPUReservationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ECS/ContainerInsights", "CpuReserved", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu reservation raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_Reservation"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processCPUReservedRawData(result *cloudwatch.GetMetricDataOutput) CPUReservedResult {
// 	var rawData CPUReservedResult
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
	AwsxCpuReservedCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxCpuReservedCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxCpuReservedCmd.PersistentFlags().String("query", "", "query")
	AwsxCpuReservedCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxCpuReservedCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxCpuReservedCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxCpuReservedCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxCpuReservedCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxCpuReservedCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxCpuReservedCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxCpuReservedCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxCpuReservedCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxCpuReservedCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxCpuReservedCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxCpuReservedCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxCpuReservedCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
