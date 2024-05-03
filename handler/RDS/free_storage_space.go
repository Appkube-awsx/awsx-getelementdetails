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

// type StorageSpace struct {
// 	Timestamp time.Time
// 	Value     float64
//}

var AwsxRDSFreeStorageSpaceCmd = &cobra.Command{
	Use:   "free_storage_space_panel",
	Short: "get free storage space metrics data",
	Long:  `command to get free storage space metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSFreeStorageSpacePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting free storage space data: ", err)
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

func GetRDSFreeStorageSpacePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
		log.Println("Error in getting for free storage space data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["FreeStorageSpace"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processedRawStorageSpaceData(result *cloudwatch.GetMetricDataOutput) []StorageSpace {
// 	var processedData []StorageSpace

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		value := *result.MetricDataResults[0].Values[i]
// 		processedData = append(processedData, StorageSpace{
// 			Timestamp: *timestamp,
// 			Value:     value,
// 		})
// 	}

// 	return processedData
// }

func init() {
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSFreeStorageSpaceCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
