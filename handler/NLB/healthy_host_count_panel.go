package NLB

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxNLBHealthyHostCountCmd = &cobra.Command{
	Use:   "nlb_healthy_host_count_panel",
	Short: "Get NLB healthy host count metrics data",
	Long:  `Command to get NLB healthy host count metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command..")
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
			jsonResp, cloudwatchMetricResp, err := GetNLBHealthyHostCountPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB healthy host count: ", err)
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

func GetNLBHealthyHostCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/NetworkELB", "HealthyHostCount", startTime, endTime, "Average", "LoadBalancer", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB healthy host count data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["HealthyHostCount"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("nlbArn", "", "NLB ARN")
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
