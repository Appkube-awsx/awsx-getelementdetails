package ApiGateway

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/spf13/cobra"
)

type TotalApiResult struct {
	Value float64 `json:"Value"`
}

var AwsxTotalApiCmd = &cobra.Command{
	Use:   "total_api_panel",
	Short: "get total api metrics data",
	Long:  `command to get total api metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetTotalApiData(clientAuth, nil)
			if err != nil {
				log.Println("Error getting total api data : ", err)
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

func GetTotalApiData(clientAuth *model.Auth, apiClient *apigateway.APIGateway) (string, map[string]float64, error) {
	cloudwatchMetricData := map[string]float64{}

	totalApis, err := GetTotalApi(clientAuth, apiClient)
	if err != nil {
		log.Println("Error in getting total functions: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TotalAPIs"] = float64(totalApis)

	log.Printf("Total APIs: %d", totalApis)

	jsonString, err := json.Marshal(TotalApiResult{Value: float64(totalApis)})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetTotalApi(clientAuth *model.Auth, apiClient *apigateway.APIGateway) (int, error) {
	if apiClient == nil {
		apiClient = awsclient.GetClient(*clientAuth, awsclient.APIGATEWAY_CLIENT).(*apigateway.APIGateway)
	}

	totalApis := 0

	// Get total count of REST APIs
	restAPIs, err := GetRestAPIs(clientAuth, apiClient)
	if err != nil {
		return 0, err
	}
	totalApis += restAPIs

	// Get total count of HTTP APIs (API Gateway V2)
	httpAPIs, err := GetHttpAPIs(clientAuth, nil)
	if err != nil {
		return 0, err
	}
	totalApis += httpAPIs

	return totalApis, nil
}

func init() {
	AwsxTotalApiCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxTotalApiCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxTotalApiCmd.PersistentFlags().String("query", "", "query")
	AwsxTotalApiCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxTotalApiCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxTotalApiCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxTotalApiCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxTotalApiCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxTotalApiCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxTotalApiCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxTotalApiCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxTotalApiCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxTotalApiCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxTotalApiCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxTotalApiCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxTotalApiCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
