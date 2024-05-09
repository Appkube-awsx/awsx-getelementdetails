package Lambda

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type LatencyGraph struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"latency_graph_panel"`
// }

var AwsxLambdaLatencyGraphCmd = &cobra.Command{
	Use:   "Latency_graph_panel",
	Short: "get Latency count graph metrics data",
	Long:  `command to get Latency count graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaLatencyGraphData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda Latency response data: ", err)
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

func GetLambdaLatencyGraphData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	//elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	//cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
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
	LatencyCount, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/Lambda", "Duration", startTime, endTime, "Average", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting lambda latency count data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Latency"] = LatencyCount

	//result := ProcessLambdaLatencyRawData(LatencyCount)

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Println("Error in marshalling json in string: ", err)
	// 	return "", nil, err
	// }
	// fmt.Println(jsonString)

	return "", cloudwatchMetricData, nil
}

// func ProcessLambdaLatencyRawData(result *cloudwatch.GetMetricDataOutput) LatencyGraph {
// 	var rawData LatencyGraph
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
	comman_function.InitAwsCmdFlags(AwsxLambdaLatencyGraphCmd)
}
