package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type memoryLimitResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Memory limits"`
// }

var AwsxEKSMemoryLimitsCmd = &cobra.Command{
	Use:   "memory_limits_panel",
	Short: "get memory_limits metrics data",
	Long:  `command to get memory_limits metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetMemoryLimitsData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory_limits: ", err)
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

func GetMemoryLimitsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "pod_memory_limit", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory limits"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSMemoryLimitsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
