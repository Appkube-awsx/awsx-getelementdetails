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

// type IndexSizeResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Index_Size"`
// }

var AwsxRDSIndexSizeCmd = &cobra.Command{
	Use:   "index_size_panel",
	Short: "Get index size metrics data",
	Long:  `Command to get index size metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetIndexSizePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting index size: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case. It prints JSON
				fmt.Println(jsonResp)
			}
		}
	},
}

func GetIndexSizePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "FreeStorageSpace", startTime, endTime, "Average", "DBInstanceIdentifier", cloudWatchClient)

	if err != nil {
		log.Println("Error in getting index size data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Index_Size"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processsRawData(result *cloudwatch.GetMetricDataOutput) IndexSizeResult {
// 	var rawData IndexSizeResult
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
	AwsxRDSIndexSizeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSIndexSizeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
