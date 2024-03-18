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

type WebSocketAPIResult struct {
	Value float64 `json:"Value"`
}

var AwsxApiGatewayWebSocketCmd = &cobra.Command{
	Use:   "websocket_api_panel",
	Short: "get WebSocket API metrics data",
	Long:  `Command to get WebSocket API metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiGatewayWebSocketAPIData(clientAuth, nil)
			if err != nil {
				log.Println("Error getting WebSocket API data : ", err)
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

func GetApiGatewayWebSocketAPIData(clientAuth *model.Auth, apiGatewayClient *apigatewayv2.ApiGatewayV2) (string, map[string]float64, error) {
	cloudwatchMetricData := map[string]float64{}

	websocketAPIs, err := GetWebSocketAPIs(clientAuth, apiGatewayClient)
	if err != nil {
		log.Println("Error in getting WebSocket APIs: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["WebSocketAPIs"] = float64(websocketAPIs)

	log.Printf("WebSocket APIs: %d", websocketAPIs)

	jsonString, err := json.Marshal(WebSocketAPIResult{Value: float64(websocketAPIs)})
	if err != nil {
		log.Println("Error in marshalling JSON in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetWebSocketAPIs(clientAuth *model.Auth, apiGatewayClient *apigatewayv2.ApiGatewayV2) (int, error) {
	if apiGatewayClient == nil {
		apiGatewayClient = awsclient.GetClient(*clientAuth, awsclient.APIGATEWAYV2_CLIENT).(*apigatewayv2.ApiGatewayV2)
	}

	websocketAPIs := 0

	input := &apigatewayv2.GetApisInput{}

	for {
		resp, err := apiGatewayClient.GetApis(input)
		if err != nil {
			return 0, err
		}

		for _, api := range resp.Items {
			if *api.ProtocolType == "WEBSOCKET" {
				websocketAPIs++
			}
		}

		// If there are more APIs to fetch, update the NextToken
		if resp.NextToken != nil {
			input.NextToken = resp.NextToken
		} else {
			break
		}
	}

	return websocketAPIs, nil
}

func init() {
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("elementId", "", "element ID")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("query", "", "query")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("cmdbApiUrl", "", "CMDB API")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("vaultUrl", "", "Vault endpoint")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("vaultToken", "", "Vault token")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("zone", "", "AWS region")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("accessKey", "", "AWS access key")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("secretKey", "", "AWS secret key")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("crossAccountRoleArn", "", "AWS cross account role ARN")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("externalId", "", "AWS external ID")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("cloudWatchQueries", "", "AWS CloudWatch metric queries")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("instanceId", "", "Instance ID")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxApiGatewayWebSocketCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
