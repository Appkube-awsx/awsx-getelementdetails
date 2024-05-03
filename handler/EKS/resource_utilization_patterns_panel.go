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

// type ResourceUtilizationResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU Utilization"`
// }

var AwsxResourceUtilizationCmd = &cobra.Command{
	Use:   "resource_utilization_patterns_panel",
	Short: "get resource utilization metrics data",
	Long:  `command to get resource utilization  metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetResourceUtilizationData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting resource utilization  data : ", err)
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

func GetResourceUtilizationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_cpu_utilization", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Resource Utilization Patterns"] = rawData

	// result := processResourceUtilizationRawData(rawData)

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Println("Error in marshalling json in string: ", err)
	// 	return "", nil, err
	// }

	return "", cloudwatchMetricData, nil
}

// func processResourceUtilizationRawData(result *cloudwatch.GetMetricDataOutput) ResourceUtilizationResult {
// 	var rawData ResourceUtilizationResult
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
	AwsxResourceUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxResourceUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxResourceUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxResourceUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxResourceUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxResourceUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxResourceUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxResourceUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxResourceUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxResourceUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxResourceUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxResourceUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxResourceUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxResourceUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxResourceUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxResourceUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
