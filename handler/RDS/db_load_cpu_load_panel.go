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

var AwsxRDSDBLoadCPUCmd = &cobra.Command{
	Use:   "db_load_cpu_panel",
	Short: "get CPU load in database operations",
	Long:  `command to get CPU load in database operations`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSDBLoadCPU(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting CPU load data: ", err)
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

func GetRDSDBLoadCPU(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := commanFunction.GetMetricDatabaseData(clientAuth, instanceId, "AWS/RDS", "DBLoadCPU", startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting CPU load data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DBLoadCPU"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processedRawdata(result *cloudwatch.GetMetricDataOutput) []struct {
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
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDBLoadCPUCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
