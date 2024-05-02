package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type MemoryUsageResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Memory Usage"`
// }

var AwsxEKSMemoryUsageCmd = &cobra.Command{
	Use:   "memory_usage_panel",
	Short: "get memory_usage metrics data",
	Long:  `command to get memory_usage metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetMemoryUsageData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory_usage: ", err)
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

func GetMemoryUsageData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_memory_utilization", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory Usage"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSMemoryUsageCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
