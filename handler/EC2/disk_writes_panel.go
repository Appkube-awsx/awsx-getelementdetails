package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

//type DiskWritePanelData struct {
//	RawData []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"Disk_Writes"`
//}

var AwsxEc2DiskWriteCmd = &cobra.Command{
	Use:   "disk_write_panel",
	Short: "get disk write metrics data",
	Long:  `command to get disk write metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskWritePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetDiskWritePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId,  "CWAgent", "diskio_writes", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk write data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Disk_Writes"] = rawData

	return "", cloudwatchMetricData, nil



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
	comman_function.InitAwsCmdFlags(AwsxEc2DiskWriteCmd)
}
