package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
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
	rawData, err := metricData.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "pod_memory_request", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory requests"] = rawData

	// Debug prints
	// log.Printf("RawData Result: %+v", rawData)

	// Process the raw data if needed
	// result := ProcessMemoryRequestRawData(rawData)

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Println("Error in marshalling json in string: ", err)
	// 	return "", nil, err
	// }

	return "", cloudwatchMetricData, nil
}

// func GetMemoryRequestMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
// 	elmType := "ContainerInsights"
// 	input := &cloudwatch.GetMetricDataInput{
// 		EndTime:   endTime,
// 		StartTime: startTime,
// 		MetricDataQueries: []*cloudwatch.MetricDataQuery{
// 			{
// 				Id: aws.String("m1"),
// 				MetricStat: &cloudwatch.MetricStat{
// 					Metric: &cloudwatch.Metric{
// 						Dimensions: []*cloudwatch.Dimension{
// 							{
// 								Name:  aws.String("ClusterName"),
// 								Value: aws.String(instanceId),
// 							},
// 						},
// 						MetricName: aws.String("pod_memory_request"),
// 						Namespace:  aws.String(elmType),
// 					},
// 					Period: aws.Int64(60),
// 					Stat:   aws.String("Average"),
// 				},
// 			},
// 		},
// 	}
// 	if cloudWatchClient == nil {
// 		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
// 	}
// 	result, err := cloudWatchClient.GetMetricData(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func ProcessMemoryRequestRawData(result *cloudwatch.GetMetricDataOutput) memoryResult {
// 	var rawData memoryResult
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
