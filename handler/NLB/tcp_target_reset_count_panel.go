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

// type TCPResetCountData struct {
// 	TCPResetCount []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"TCPResetCount"`
// }

var AwsxNLBTCPResetCountCmd = &cobra.Command{
	Use:   "tcp_target_reset_count_panel",
	Short: "Get NLB TCP target reset count metrics data",
	Long:  `Command to get NLB TCP target reset count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNLBTCPResetCountPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB TCP target reset count: ", err)
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

func GetNLBTCPResetCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Fetch raw data
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/NetworkELB", "TCP_Target_Reset_Count", startTime, endTime, "Sum", "LoadBalancer", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB TCP target reset count data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TCP Target Reset Count"] = rawData

	var totalSum float64
	for _, value := range rawData.MetricDataResults {
		for _, datum := range value.Values {
			totalSum += *datum
		}
	}
	totalSumStr := fmt.Sprintf("{request count: %f}", totalSum)
	return totalSumStr, cloudwatchMetricData, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxNLBTCPResetCountCmd)

}
