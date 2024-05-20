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

var AwsxNLBTcpProcesedBytesCmd = &cobra.Command{
	Use:   "tcp_procesed_bytes_panel",
	Short: "Get NLB tcp procesed bytes",
	Long:  `Command to get NLB procesed bytes`,

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
			jsonResp, cloudwatchMetricResp, err := GetNLBTcpProcesedBytes(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB tcp procesed bytes data: ", err)
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

func GetNLBTcpProcesedBytes(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := comman_function.ParseTimes(cmd)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/NetworkELB", "ProcessedBytes_TCP", startTime, endTime, "Average", "LoadBalancer", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB Tcp processed bytes data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["ProcessedBytes_TCP"] = rawData

	var totalSum float64
	for _, value := range rawData.MetricDataResults {
		for _, datum := range value.Values {
			totalSum += *datum
		}
	}
	totalSumStr := fmt.Sprintf("{ProcesedBytes count: %f}", totalSum)
	return totalSumStr, cloudwatchMetricData, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxNLBTcpProcesedBytesCmd)
}
