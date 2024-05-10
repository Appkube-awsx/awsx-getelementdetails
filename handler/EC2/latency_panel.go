package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type Ec2Latency struct {
	InboundTraffic  float64 `json:"inboundTraffic"`
	OutboundTraffic float64 `json:"outboundTraffic"`
	DataTransferred float64 `json:"dataTransferred"`
	Latency         float64 `json:"latency"`
}

var AwsxEc2LatencyCmd = &cobra.Command{
	Use:   "latency_panel",
	Short: "get latency metrics data",
	Long:  `command to get latency metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLatencyPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network utilization: ", err)
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

func GetLatencyPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", nil, err
		}
		instanceId = cmdbData.InstanceId

	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Inbound Traffic
	inboundTraffic, err := GetLatencyMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "NetworkIn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting inbound traffic: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["InboundTraffic"] = inboundTraffic
	// Get Outbound Traffic
	outboundTraffic, err := GetLatencyMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "NetworkOut", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

	// Calculate Data Transferred (sum of inbound and outbound)
	dataTransferred := *inboundTraffic.MetricDataResults[0].Values[0] + *outboundTraffic.MetricDataResults[0].Values[0]
	latency := dataTransferred / 2
	cloudwatchMetricData["DataTransferred"] = createMetricData(dataTransferred)

	jsonOutput := Ec2Latency{
		InboundTraffic:  *inboundTraffic.MetricDataResults[0].Values[0],
		OutboundTraffic: *outboundTraffic.MetricDataResults[0].Values[0],
		DataTransferred: dataTransferred,
		Latency:         latency,
	}

	jsonString, err := json.Marshal(struct{ Latency float64 }{Latency: jsonOutput.Latency})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLatencyMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)
	elmType := "AWS/EC2"
	if elementType == "EC2" {
		elmType = "AWS/" + elementType
	}
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String(instanceID),
							},
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
				},
			},
		},
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func createMetricData(value float64) *cloudwatch.GetMetricDataOutput {
	return &cloudwatch.GetMetricDataOutput{
		MetricDataResults: []*cloudwatch.MetricDataResult{
			{
				Values: []*float64{&value},
			},
		},
	}
}
func init() {
	AwsxEc2LatencyCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2LatencyCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2LatencyCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2LatencyCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2LatencyCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2LatencyCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2LatencyCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2LatencyCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2LatencyCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2LatencyCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2LatencyCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2LatencyCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2LatencyCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2LatencyCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2LatencyCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2LatencyCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

//package EC2

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/Appkube-awsx/awsx-common/authenticate"
// 	"github.com/Appkube-awsx/awsx-common/model"
// 	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
// 	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
// 	"github.com/aws/aws-sdk-go/service/cloudwatch"
// 	"github.com/spf13/cobra"
// )

// type Ec2Latency struct {
// 	InboundTraffic  float64 `json:"InboundTraffic"`
// 	OutboundTraffic float64 `json:"OutboundTraffic"`
// 	DataTransferred float64 `json:"DataTransferred"`
// 	Latency         float64 `json:"latency"`
// }

// var AwsxEc2LatencyCmd = &cobra.Command{
// 	Use:   "latency_panel",
// 	Short: "get latency metrics data",
// 	Long:  `command to get latency metrics data`,

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
// 			jsonResp, cloudwatchMetricResp, err := GetLatencyPanel(cmd, clientAuth, nil)
// 			if err != nil {
// 				log.Println("Error getting network utilization: ", err)
// 				return
// 			}
// 			if responseType == "frame" {
// 				fmt.Println(cloudwatchMetricResp)
// 			} else {
// 				// default case. it prints json
// 				fmt.Println(jsonResp)
// 			}
// 		}

// 	},
// }

// func GetLatencyPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

// 	elementType, _ := cmd.PersistentFlags().GetString("elementType")
// 	fmt.Println(elementType)
// 	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

// 	startTime, endTime, err := comman_function.ParseTimes(cmd)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("error parsing time: %v", err)
// 	}

// 	instanceId, err = comman_function.GetCmdbData(cmd)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
// 	}
// 	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

// 	// Get Inbound Traffic
// 	inboundTraffic, err := metricData.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, startTime, endTime, "Average", "NetworkIn", cloudWatchClient)
// 	if err != nil {
// 		log.Println("Error in getting inbound traffic: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["InboundTraffic"] = inboundTraffic
// 	// Get Outbound Traffic
// 	outboundTraffic, err := metricData.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, startTime, endTime, "Maximum", "NetworkOut", cloudWatchClient)
// 	if err != nil {
// 		log.Println("Error in getting outbound traffic: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

// 	// Calculate Data Transferred (sum of inbound and outbound)
// 	dataTransferred := *inboundTraffic.MetricDataResults[0].Values[0] + *outboundTraffic.MetricDataResults[0].Values[0]
// 	latency := dataTransferred / 2
// 	cloudwatchMetricData["DataTransferred"] = createMetricData(dataTransferred)

// 	jsonOutput := Ec2Latency{
// 		InboundTraffic:  *inboundTraffic.MetricDataResults[0].Values[0],
// 		OutboundTraffic: *outboundTraffic.MetricDataResults[0].Values[0],
// 		DataTransferred: dataTransferred,
// 		Latency:         latency,
// 	}

// 	jsonString, err := json.Marshal(struct{ Latency float64 }{Latency: jsonOutput.Latency})
// 	if err != nil {
// 		log.Println("Error in marshalling json in string: ", err)
// 		return "", nil, err
// 	}

// 	return string(jsonString), cloudwatchMetricData, nil
// }

// // func GetLatencyMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
// // 	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)
// // 	elmType := "AWS/EC2"
// // 	if elementType == "EC2" {
// // 		elmType = "AWS/" + elementType
// // 	}
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
// // 								Name:  aws.String("InstanceId"),
// // 								Value: aws.String(instanceID),
// // 							},
// // 						},
// // 						MetricName: aws.String(metricName),
// // 						Namespace:  aws.String(elmType),
// // 					},
// // 					Period: aws.Int64(300),
// // 					Stat:   aws.String(statistic),
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

// // 	return result, nil
// // }

// func createMetricData(value float64) *cloudwatch.GetMetricDataOutput {
// 	return &cloudwatch.GetMetricDataOutput{
// 		MetricDataResults: []*cloudwatch.MetricDataResult{
// 			{
// 				Values: []*float64{&value},
// 			},
// 		},
// 	}
// }
// func init() {
// 	AwsxEc2LatencyCmd.PersistentFlags().String("elementId", "", "element id")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("elementType", "", "element type")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("query", "", "query")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("vaultToken", "", "vault token")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("zone", "", "aws region")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("accessKey", "", "aws access key")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("secretKey", "", "aws secret key")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("externalId", "", "aws external id")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("instanceId", "", "instance id")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("startTime", "", "start time")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("endTime", "", "endcl time")
// 	AwsxEc2LatencyCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
// }
