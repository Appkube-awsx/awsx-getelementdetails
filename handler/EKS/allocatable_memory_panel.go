package EKS

// import (
// 	"fmt"
// 	"log"

// 	"github.com/Appkube-awsx/awsx-common/authenticate"
// 	"github.com/Appkube-awsx/awsx-common/model"
// 	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
// 	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
// 	"github.com/aws/aws-sdk-go/service/cloudwatch"
// 	"github.com/spf13/cobra"
// )

// // type TimeSeriesMemData struct {
// // 	Timestamp      time.Time
// // 	AllocatableMem float64
// // 	// ReservedMem    float64
// // }

// // type AllocateMemResult struct {
// // 	AllocatableMemory []TimeSeriesMemData `json:"AllocatableMemory"`
// // }

// var AwsxEKSAllocatableMemCmd = &cobra.Command{
// 	Use:   "allocatable_mem_panel",
// 	Short: "get allocatable memory metrics data",
// 	Long:  `command to get allocatable memory metrics data`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("running from child command")
// 		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
// 		if err != nil {
// 			log.Printf("Error during authentication: %v\n", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return
// 			}
// 			return
// 		}
// 		if authFlag {
// 			responseType, _ := cmd.PersistentFlags().GetString("responseType")
// 			jsonResp, cloudwatchMetricResp, err := GetAllocatableMemData(cmd, clientAuth, nil)
// 			if err != nil {
// 				log.Println("Error getting allocatable memory: ", err)
// 				return
// 			}
// 			if responseType == "frame" {
// 				fmt.Println(cloudwatchMetricResp)
// 			} else {
// 				fmt.Println(jsonResp)
// 			}
// 		}

// 	},
// }

// func GetAllocatableMemData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

// 	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
// 	elementType, _ := cmd.PersistentFlags().GetString("elementType")
// 	fmt.Println(elementType)

// 	startTime, endTime, err := commanFunction.ParseTimes(cmd)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("error parsing time: %v", err)
// 	}

// 	instanceId, err = commanFunction.GetCmdbData(cmd)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
// 	}

// 	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
// 	// Fetch raw data
// 	rawData, err := metricData.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "node_memory_limit", "Average", "node_memory_reserved_capacity", startTime, endTime, cloudWatchClient)
// 	if err != nil {
// 		log.Println("Error in getting raw data: ", err)
// 		return "", nil, err
// 	}

// 	cloudwatchMetricData["Allocatable_memory"] = rawData

// 	return "", cloudwatchMetricData, nil
// }

// // func GetAllocatableMemMetricData(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
// // 	elmType := "ContainerInsights"
// // 	input := &cloudwatch.GetMetricDataInput{
// // 		EndTime:   endTime,
// // 		StartTime: startTime,
// // 		MetricDataQueries: []*cloudwatch.MetricDataQuery{
// // 			{
// // 				Id: aws.String("m1"),
// // 				MetricStat: &cloudwatch.MetricStat{
// // 					Metric: &cloudwatch.Metric{
// // 						Dimensions: []*cloudwatch.Dimension{
// // 							{
// // 								Name:  aws.String("ClusterName"),
// // 								Value: aws.String(instanceId),
// // 							},
// // 						},
// // 						MetricName: aws.String("node_memory_limit"),
// // 						Namespace:  aws.String(elmType),
// // 					},
// // 					Period: aws.Int64(60),
// // 					Stat:   aws.String("Average"),
// // 				},
// // 			},
// // 			{
// // 				Id: aws.String("m2"),
// // 				MetricStat: &cloudwatch.MetricStat{
// // 					Metric: &cloudwatch.Metric{
// // 						Dimensions: []*cloudwatch.Dimension{
// // 							{
// // 								Name:  aws.String("ClusterName"),
// // 								Value: aws.String(instanceId),
// // 							},
// // 						},
// // 						MetricName: aws.String("node_memory_reserved_capacity"),
// // 						Namespace:  aws.String(elmType),
// // 					},
// // 					Period: aws.Int64(60),
// // 					Stat:   aws.String("Average"),
// // 				},
// // 			},
// // 		},
// // 	}
// // 	if cloudWatchClient == nil {
// // 		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
// // 	}
// // 	result, err := cloudWatchClient.GetMetricData(input)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	// fmt.Println("result",result)
// // 	// fmt.Println("instanceId",instanceId)
// // 	// fmt.Println("elmType",elmType)
// // 	// fmt.Println("input",input)

// // 	return result, nil
// // }

// // func processMemRawData(result *cloudwatch.GetMetricDataOutput) AllocateMemResult {
// // 	var rawData AllocateMemResult
// // 	rawData.AllocatableMemory = make([]TimeSeriesMemData, len(result.MetricDataResults[0].Timestamps))

// // 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// // 		rawData.AllocatableMemory[i].Timestamp = *timestamp
// // 		memLimit := *result.MetricDataResults[0].Values[i]
// // 		reservedCapacity := *result.MetricDataResults[1].Values[i]
// // 		fmt.Println("memlimit",memLimit)
// // 		fmt.Println("reserved capacity",reservedCapacity)
// // 		allocatableMem := memLimit - reservedCapacity

// // 		// Only include the calculated allocatable memory in the result
// // 		rawData.AllocatableMemory[i].AllocatableMem = allocatableMem
// // 	}
// // 	// fmt.Println("raw data",rawData)
// // 	return rawData
// // }

// func init() {
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("elementId", "", "element id")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("elementType", "", "element type")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("query", "", "query")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("vaultToken", "", "vault token")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("zone", "", "aws region")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("accessKey", "", "aws access key")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("secretKey", "", "aws secret key")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("externalId", "", "aws external id")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("instanceId", "", "instance id")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("startTime", "", "start time")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("endTime", "", "endcl time")
// 	AwsxEKSAllocatableMemCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
// }
