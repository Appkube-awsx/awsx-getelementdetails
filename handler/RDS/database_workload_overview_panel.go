package RDS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type DatabaseWorkloadOverview struct {
// 	Timestamp time.Time
// 	Value     float64
// }

var AwsxRDSDBLoadCmd = &cobra.Command{
	Use:   "db_load_panel",
	Short: "get database workload overview metrics data",
	Long:  `Command to get database workload overview metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			jsonResp, cloudwatchMetricResp, err := GetRDSDBLoadPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting database workload overview data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case: print JSON
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetRDSDBLoadPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "DBLoad", startTime, endTime, "Average", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error getting database workload overview data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DBLoad"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processedRawDatabaseWorkloadOverviewData(result *cloudwatch.GetMetricDataOutput) []DatabaseWorkloadOverview {
// 	var processedData []DatabaseWorkloadOverview

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		value := *result.MetricDataResults[0].Values[i]
// 		processedData = append(processedData, DatabaseWorkloadOverview{
// 			Timestamp: *timestamp,
// 			Value:     value,
// 		})
// 	}

// 	return processedData
// }

func init() {
	AwsxRDSDBLoadCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDBLoadCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDBLoadCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDBLoadCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDBLoadCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDBLoadCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDBLoadCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDBLoadCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDBLoadCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDBLoadCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDBLoadCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDBLoadCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDBLoadCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDBLoadCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDBLoadCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDBLoadCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
