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

// type DBResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Database_Connections"`
// }

var AwsxRDSDatabaseConnectionsCmd = &cobra.Command{
	Use:   "database_connections_panel",
	Short: "get database connections metrics data",
	Long:  `command to get database connections metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDatabaseConnectionsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting database connections: ", err)
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

func GetDatabaseConnectionsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := commanFunction.GetMetricDatabaseData(clientAuth, instanceId, "AWS/RDS", "DatabaseConnections", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting database connection data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Database_Connections"] = rawData

	return "", cloudwatchMetricData, nil

}

// func processRawData(result *cloudwatch.GetMetricDataOutput) DBResult {
// 	var rawData DBResult
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
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDatabaseConnectionsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
