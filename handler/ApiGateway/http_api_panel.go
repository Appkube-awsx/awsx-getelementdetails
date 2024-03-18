package ApiGateway

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/spf13/cobra"
)

type HttpAPIResult struct {
	Value float64 `json:"Value"`
}

var AwsxApiGatewayHTTPCmd = &cobra.Command{
	Use:   "http_api_panel",
	Short: "get HTTP API metrics data",
	Long:  `Command to get HTTP API metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			jsonResp, cloudwatchMetricResp, err := GetApiGatewayHttpApiData(clientAuth, nil)
			if err != nil {
				log.Println("Error getting HTTP API data : ", err)
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

func GetApiGatewayHttpApiData(clientAuth *model.Auth, apiGatewayClient *apigatewayv2.ApiGatewayV2) (string, map[string]float64, error) {
	cloudwatchMetricData := map[string]float64{}

	httpAPIs, err := GetHttpAPIs(clientAuth, apiGatewayClient)
	if err != nil {
		log.Println("Error in getting HTTP APIs: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["HTTPAPIs"] = float64(httpAPIs)

	log.Printf("HTTP APIs: %d", httpAPIs)

	jsonString, err := json.Marshal(HttpAPIResult{Value: float64(httpAPIs)})
	if err != nil {
		log.Println("Error in marshalling JSON in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetHttpAPIs(clientAuth *model.Auth, apiGatewayClient *apigatewayv2.ApiGatewayV2) (int, error) {
	if apiGatewayClient == nil {
		apiGatewayClient = awsclient.GetClient(*clientAuth, awsclient.APIGATEWAYV2_CLIENT).(*apigatewayv2.ApiGatewayV2)
	}

	httpAPIs := 0

	input := &apigatewayv2.GetApisInput{}

	for {
		resp, err := apiGatewayClient.GetApis(input)
		if err != nil {
			return 0, err
		}

		for _, api := range resp.Items {
			if *api.ProtocolType == "HTTP" {
				httpAPIs++
			}
		}

		// If there are more APIs to fetch, update the position
		if resp.NextToken != nil {
			input.NextToken = resp.NextToken
		} else {
			break
		}
	}

	return httpAPIs, nil
}

func init() {
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("elementId", "", "element ID")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("query", "", "query")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("cmdbApiUrl", "", "CMDB API")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("vaultUrl", "", "Vault endpoint")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("vaultToken", "", "Vault token")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("zone", "", "AWS region")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("accessKey", "", "AWS access key")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("secretKey", "", "AWS secret key")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("crossAccountRoleArn", "", "AWS cross account role ARN")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("externalId", "", "AWS external ID")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("cloudWatchQueries", "", "AWS CloudWatch metric queries")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("instanceId", "", "Instance ID")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxApiGatewayHTTPCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
