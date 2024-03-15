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

type RestAPIResult struct {
	Value float64 `json:"Value"`
}

var AwsxApiGatewayRestAPICmd = &cobra.Command{
	Use:   "rest_api_panel",
	Short: "get rest API metrics data",
	Long:  `Command to get rest API metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiGatewayRestAPIData(clientAuth, nil)
			if err != nil {
				log.Println("Error getting rest API data : ", err)
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

func GetApiGatewayRestAPIData(clientAuth *model.Auth, apiGatewayClient *apigateway.APIGateway) (string, map[string]float64, error) {
	cloudwatchMetricData := map[string]float64{}

	restAPIs, err := GetRestAPIs(clientAuth, apiGatewayClient)
	if err != nil {
		log.Println("Error in getting rest APIs: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RestAPIs"] = float64(restAPIs)

	log.Printf("Rest APIs: %d", restAPIs)

	jsonString, err := json.Marshal(RestAPIResult{Value: float64(restAPIs)})
	if err != nil {
		log.Println("Error in marshalling JSON in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetRestAPIs(clientAuth *model.Auth, apiGatewayClient *apigateway.APIGateway) (int, error) {
	if apiGatewayClient == nil {
		apiGatewayClient = awsclient.GetClient(*clientAuth, awsclient.APIGATEWAY_CLIENT).(*apigateway.APIGateway)
	}

	restAPIs := 0

	input := &apigateway.GetRestApisInput{}

	for {
		resp, err := apiGatewayClient.GetRestApis(input)
		if err != nil {
			return 0, err
		}

		restAPIs += len(resp.Items)

		// If there are more APIs to fetch, update the position
		if resp.Position != nil {
			input.Position = resp.Position
		} else {
			break
		}
	}

	return restAPIs, nil
}

func init() {
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("elementId", "", "element ID")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("query", "", "query")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("cmdbApiUrl", "", "CMDB API")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("vaultUrl", "", "Vault endpoint")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("vaultToken", "", "Vault token")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("zone", "", "AWS region")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("accessKey", "", "AWS access key")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("secretKey", "", "AWS secret key")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("crossAccountRoleArn", "", "AWS cross account role ARN")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("externalId", "", "AWS external ID")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("cloudWatchQueries", "", "AWS CloudWatch metric queries")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("instanceId", "", "Instance ID")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("endTime", "", "End time")
	AwsxApiGatewayRestAPICmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
