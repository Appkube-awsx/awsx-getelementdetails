package S3

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	 "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxErrorsCmd = &cobra.Command{
	Use:   "errors_panel",
	Short: "get errors metrics data for s3",
	Long:  `command to get errors metrics data for s3`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command..")
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
			jsonResp, cloudwatchMetricResp, err := GetS3ErrorsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting errors: ", err)
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

func GetS3ErrorsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	//bucketName := "abdulweb.com"

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
	fourxxErrorsData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "4xxErrors", startTime, endTime, "Average", "bucketName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting s3 4xx errors data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["4xxErrorsData"] = fourxxErrorsData

	fivexxErrorsData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "5xxErrors", startTime, endTime, "Average", "bucketName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting s3 4xx errors data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["5xxErrorsData"] = fivexxErrorsData
	return "", cloudwatchMetricData, nil

}
func init() {
	comman_function.InitAwsCmdFlags(AwsxErrorsCmd)
}
