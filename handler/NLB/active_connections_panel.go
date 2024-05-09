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

// type ActiveConnectionsData struct {
// 	ActiveConnections []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"ActiveConnections"`
// }

var AwsxNLBActiveConnectionsCmd = &cobra.Command{
	Use:   "nlb_active_connections_panel",
	Short: "Get NLB active connections metrics data",
	Long:  `Command to get NLB active connections metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNLBActiveConnectionsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB active connections: ", err)
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

func GetNLBActiveConnectionsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	
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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "ActiveFlowCount", startTime, endTime, "Sum", "LoadBalancer", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB active connections data: ", err)
		return "", nil, err
	
	}
	cloudwatchMetricData["ActiveConnections"] = rawData
	return "", cloudwatchMetricData, nil
}
	



// func processNLBActiveConnectionsRawData(result *cloudwatch.GetMetricDataOutput) ActiveConnectionsData {
// 	var rawData ActiveConnectionsData
// 	rawData.ActiveConnections = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.ActiveConnections[i].Timestamp = *timestamp
// 		rawData.ActiveConnections[i].Value = *result.MetricDataResults[0].Values[i]
// 	}

// 	return rawData
// }

func init() {
	comman_function.InitAwsCmdFlags(AwsxNLBActiveConnectionsCmd)
}
