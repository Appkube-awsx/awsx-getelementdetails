package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

//type CpuUtilizationsResult struct {
//	RawData []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"cpu utilization graph"`
//}

var AwsxEc2CpuUtilizationGraphsCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GetMemoryUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization graph: ", err)
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

func GetCpuUtilizationGraphPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
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

	// Get average utilization
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting rawdata: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU Utilization"] = rawData

	//result := processCpuUtilizationGraphRawData(rawData)

	//jsonString, err := json.Marshal(result)
	//if err != nil {
	//	log.Println("Error in marshalling json in string: ", err)
	//	return "", nil, err
	//}

	return "", cloudwatchMetricData, nil

}

//
//func processCpuUtilizationGraphRawData(result *cloudwatch.GetMetricDataOutput) CpuUtilizationsResult {
//	var rawData CpuUtilizationsResult
//	rawData.RawData = make([]struct {
//		Timestamp time.Time
//		Value     float64
//	}, len(result.MetricDataResults[0].Timestamps))
//
//	for i, timestamp := range result.MetricDataResults[0].Timestamps {
//		rawData.RawData[i].Timestamp = *timestamp
//		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
//	}
//
//	return rawData
//}

func init() {
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUtilizationGraphsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
