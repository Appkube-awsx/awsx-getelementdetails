package RDS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type CPUUtilizationResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU Utilization"`
// }

var AwsxRDSCpuUtilizationGraphCmd = &cobra.Command{
	Use:   "cpu_utilization_graph_panel",
	Short: "get cpu utilization graph metrics data",
	Long:  `command to get cpu utilization graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSCPUUtilizationGraphPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization graph data : ", err)
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

func GetRDSCPUUtilizationGraphPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricDatabaseData(clientAuth, instanceId, "AWS/RDS", "CPUUtilization", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu utilization data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU Utilization"] = rawData

	// // Process the raw data if needed
	// result := processCPURawData(rawData)

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Println("Error in marshalling json in string: ", err)
	// 	return "", nil, err
	// }

	return "", cloudwatchMetricData, nil
}

// func processCPURawData(result *cloudwatch.GetMetricDataOutput) CPUUtilizationResult {
// 	var rawData CPUUtilizationResult
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
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSCpuUtilizationGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
