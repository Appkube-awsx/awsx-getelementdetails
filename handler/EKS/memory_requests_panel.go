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

// type memoryResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Memory requests"`
// }

var AwsxEKSMemoryRequestsCmd = &cobra.Command{
	Use:   "memory_requests_panel",
	Short: "get memory_requests metrics data",
	Long:  `command to get memory_requests metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetMemoryRequestData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory_requests: ", err)
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

func GetMemoryRequestData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	// Fetch raw data
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "pod_memory_request", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory requests"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSMemoryRequestsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
