package EKS

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NetworkResult struct {
	InboundTraffic  float64 `json:"inboundTraffic"`
	OutboundTraffic float64 `json:"outboundTraffic"`
	DataTransferred float64 `json:"dataTransferred"`
}

func GetNetworkUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
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
	inboundTraffic, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "pod_network_rx_bytes")
	if err != nil {
		log.Println("Error in getting inbound traffic: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["InboundTraffic"] = inboundTraffic

	// Get Outbound Traffic
	outboundTraffic, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "pod_network_tx_bytes")
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

	// Calculate Data Transferred (sum of inbound and outbound)
	dataTransferred := *inboundTraffic.MetricDataResults[0].Values[0] + *outboundTraffic.MetricDataResults[0].Values[0]
	cloudwatchMetricData["DataTransferred"] = createMetricDataOutput(dataTransferred)

	jsonOutput := NetworkResult{
		InboundTraffic:  *inboundTraffic.MetricDataResults[0].Values[0],
		OutboundTraffic: *outboundTraffic.MetricDataResults[0].Values[0],
		DataTransferred: dataTransferred,
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNetworkMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time, metricName string) (*cloudwatch.GetMetricDataOutput, error) {
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
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Sum"),
				},
			},
		},
	}
	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func extractMetricValue(result *cloudwatch.GetMetricDataOutput, index int) float64 {
	if len(result.MetricDataResults) > 0 && len(result.MetricDataResults[0].Values) > index {
		return *result.MetricDataResults[0].Values[index]
	}
	return 0
}

func createMetricDataOutput(value float64) *cloudwatch.GetMetricDataOutput {
	return &cloudwatch.GetMetricDataOutput{
		MetricDataResults: []*cloudwatch.MetricDataResult{
			{
				Values: []*float64{&value},
			},
		},
	}
}

// package EKS

// import (
// 	"encoding/json"
// 	"log"
// 	"time"

// 	"github.com/Appkube-awsx/awsx-common/awsclient"
// 	"github.com/Appkube-awsx/awsx-common/model"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/cloudwatch"
// 	"github.com/spf13/cobra"
// )

// type NetworkResult struct {
// 	InboundTraffic  float64 `json:"inboundTraffic"`
// 	OutboundTraffic float64 `json:"outboundTraffic"`
// 	DataTransferred float64 `json:"dataTransferred"`
// }

// func GetNetworkUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
// 	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
// 	namespace, _ := cmd.PersistentFlags().GetString("elementType")
// 	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
// 	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

// 	var startTime, endTime *time.Time

// 	// Parse start time if provided
// 	if startTimeStr != "" {
// 		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing start time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return "", nil, err
// 			}
// 			return "", nil, err
// 		}
// 		startTime = &parsedStartTime
// 	} else {
// 		defaultStartTime := time.Now().Add(-5 * time.Minute)
// 		startTime = &defaultStartTime
// 	}

// 	if endTimeStr != "" {
// 		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing end time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return "", nil, err
// 			}
// 			return "", nil, err
// 		}
// 		endTime = &parsedEndTime
// 	} else {
// 		defaultEndTime := time.Now()
// 		endTime = &defaultEndTime
// 	}

// 	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

// 	// Get Inbound Traffic
// 	inboundTraffic, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "NetworkPacketsIn")
// 	if err != nil {
// 		log.Println("Error in getting inbound traffic: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["InboundTraffic"] = inboundTraffic

// 	// Get Outbound Traffic
// 	outboundTraffic, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "NetworkPacketsOut")
// 	if err != nil {
// 		log.Println("Error in getting outbound traffic: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

// 	// Get Data Transferred
// 	dataTransferred, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "NetworkOut")
// 	if err != nil {
// 		log.Println("Error in getting data transferred: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["DataTransferred"] = dataTransferred

// 	jsonOutput := NetworkResult{
// 		InboundTraffic:  *inboundTraffic.MetricDataResults[0].Values[0],
// 		OutboundTraffic: *outboundTraffic.MetricDataResults[0].Values[0],
// 		DataTransferred: *dataTransferred.MetricDataResults[0].Values[0],
// 	}

// 	jsonString, err := json.Marshal(jsonOutput)
// 	if err != nil {
// 		log.Println("Error in marshalling json in string: ", err)
// 		return "", nil, err
// 	}

// 	return string(jsonString), cloudwatchMetricData, nil
// }

// func GetNetworkMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time, metricName string) (*cloudwatch.GetMetricDataOutput, error) {
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
// 								Value: aws.String(clusterName),
// 							},
// 						},
// 						MetricName: aws.String(metricName),
// 						Namespace:  aws.String(namespace),
// 					},
// 					Period: aws.Int64(300),
// 					Stat:   aws.String("Sum"),
// 				},
// 			},
// 		},
// 	}
// 	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
// 	result, err := cloudWatchClient.GetMetricData(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// package EKS

// import (
// 	"encoding/json"
// 	"log"
// 	"time"

// 	"github.com/Appkube-awsx/awsx-common/awsclient"
// 	"github.com/Appkube-awsx/awsx-common/model"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/cloudwatch"
// 	"github.com/spf13/cobra"
// )

// type NetworkResult struct {
// 	InboundTraffic  float64 `json:"inboundTraffic"`
// 	OutboundTraffic float64 `json:"outboundTraffic"`
// 	DataTransferred float64 `json:"dataTransferred"`
// }

// func GetNetworkPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

// 	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
// 	namespace, _ := cmd.PersistentFlags().GetString("elementType")
// 	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
// 	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

// 	var startTime, endTime *time.Time

// 	// Parse start time if provided
// 	if startTimeStr != "" {
// 		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing start time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return "", nil, err
// 			}
// 			return "", nil, err
// 		}
// 		startTime = &parsedStartTime
// 	} else {
// 		defaultStartTime := time.Now().Add(-5 * time.Minute)
// 		startTime = &defaultStartTime
// 	}

// 	if endTimeStr != "" {
// 		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing end time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return "", nil, err
// 			}
// 			return "", nil, err
// 		}
// 		endTime = &parsedEndTime
// 	} else {
// 		defaultEndTime := time.Now()
// 		endTime = &defaultEndTime
// 	}
// 	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

// 	// Fetch Inbound Traffic
// 	inboundTraffic, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "InboundTraffic")
// 	if err != nil {
// 		log.Println("Error in getting Inbound Traffic: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["InboundTraffic"] = inboundTraffic

// 	// Fetch Outbound Traffic
// 	outboundTraffic, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "OutboundTraffic")
// 	if err != nil {
// 		log.Println("Error in getting Outbound Traffic: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

// 	// Fetch Data Transferred
// 	dataTransferred, err := GetNetworkMetricData(clientAuth, clusterName, namespace, startTime, endTime, "DataTransferred")
// 	if err != nil {
// 		log.Println("Error in getting Data Transferred: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["DataTransferred"] = dataTransferred

// 	jsonOutput := NetworkResult{
// 		InboundTraffic:  *inboundTraffic.MetricDataResults[0].Values[0],
// 		OutboundTraffic: *outboundTraffic.MetricDataResults[0].Values[0],
// 		DataTransferred: *dataTransferred.MetricDataResults[0].Values[0],
// 	}

// 	jsonString, err := json.Marshal(jsonOutput)
// 	if err != nil {
// 		log.Println("Error in marshalling json in string: ", err)
// 		return "", nil, err
// 	}

// 	return string(jsonString), cloudwatchMetricData, nil
// }

// func GetNetworkMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time, metricName string) (*cloudwatch.GetMetricDataOutput, error) {
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
// 								Value: aws.String(clusterName),
// 							},
// 						},
// 						MetricName: aws.String("node_network_total_bytes"),
// 						Namespace:  aws.String(namespace),
// 					},
// 					Period: aws.Int64(300),
// 					Stat:   aws.String("Average"), // Assuming "Average" as the statistic, adjust if needed
// 				},
// 			},
// 		},
// 	}
// 	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
// 	result, err := cloudWatchClient.GetMetricData(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
