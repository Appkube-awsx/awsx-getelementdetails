package NLB

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type NlbTargetErrorCountTime struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"target_error_count_panel"`
// }

var AwsxNlbTargetErrorCountCmd = &cobra.Command{
	Use:   "target_error_count_panel",
	Short: "get target error count metrics data",
	Long:  `command to get target count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetTargetErrorCountData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting nlb target response data: ", err)
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

func GetTargetErrorCountData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	startTime, endTime, err := comman_function.ParseTimes(cmd)

	

	
		if err != nil {
			return "", nil, fmt.Errorf("error parsing time: %v", err)
			instanceId, err = comman_function.GetCmdbData(cmd)
		}
		if err != nil {
			return "", nil, fmt.Errorf("error getting instance ID: %v", err)
		
		}
		

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "ClientTLSNegotiationErrorCount", startTime, endTime, "Sum", "LoadBalancer", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB active connections data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TargetTLSErrorCount"] = rawData

	

	return "", cloudwatchMetricData, nil
}


		



func init() {
	comman_function.InitAwsCmdFlags(AwsxNlbTargetErrorCountCmd)

}
