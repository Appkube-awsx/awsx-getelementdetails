package RDS

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

var AwsxRDSDBLoadNonCPUCmd = &cobra.Command{
	Use:   "db_load_non_cpu_panel",
	Short: "get non-cpu load in database operations",
	Long:  `command to get non-cpu load in database operations`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSDBLoadNonCPU(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting non-cpu load data: ", err)
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

func GetRDSDBLoadNonCPU(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := metricData.GetMetricDatabaseData(clientAuth, instanceId, "AWS/RDS", "DBLoadNonCPU", startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting non-cpu load data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DBLoadNonCPU"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processedRawDataa(result *cloudwatch.GetMetricDataOutput) []struct {
// 	Timestamp time.Time
// 	Value     float64
// } {
// 	var processedData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		value := *result.MetricDataResults[0].Values[i]
// 		processedData = append(processedData, struct {
// 			Timestamp time.Time
// 			Value     float64
// 		}{Timestamp: *timestamp, Value: value})
// 	}

// 	return processedData
// }

func init() {
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDBLoadNonCPUCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
