package EC2

import (
	//"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxEc2DiskReadCmd = &cobra.Command{
	Use:   "disk_read_panel",
	Short: "get disk read metrics data",
	Long:  `command to get disk read metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskReadPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk read  utilization: ", err)
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

func GetDiskReadPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "CWAgent", "diskio_reads", startTime, endTime, "Sum", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk read data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Disk_Reads"] = rawData
	//jsonOutput := make(map[string]float64)
	//if len(rawData.MetricDataResults) > 0 && len(rawData.MetricDataResults[0].Values) > 0 {
		//jsonOutput["Disk_Reads"] = *rawData.MetricDataResults[0].Values[0]
		//return "", cloudwatchMetricData, nil
		//jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	//return string(jsonString), cloudwatchMetricData, nil
	var totalSum float64
	for _, value := range rawData.MetricDataResults {
		for _, datum := range value.Values {
			totalSum += *datum
		}
	}
	totalSumStr := fmt.Sprintf("{disk io reads count: %f}", totalSum)
	return totalSumStr, cloudwatchMetricData, nil

}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2DiskReadCmd)
}
