package Lambda

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"
)

var AwsxLambdaFullConcurrencyCommmand = &cobra.Command{
	Use:   "full_concurrency_panel",
	Short: "get full concurrency metrics data",
	Long:  `command to get full concurrency metrics data`,
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
			jsonResp, resp, err := GetLambdaFullConcurrencyData(cmd, clientAuth)
			if err != nil {
				log.Println("Error getting full concurrency data : ", err)
				return
			}
			if responseType == "json" {
				fmt.Println(jsonResp)
			} else {
				fmt.Println(resp)
			}
		}
	},
}

func GetLambdaFullConcurrencyData(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]int, error) {
	lambdaClient := awsclient.GetClient(*clientAuth, awsclient.LAMBDA_CLIENT).(*lambda.Lambda)
    input := lambda.GetAccountSettingsInput{}
	result, err := lambdaClient.GetAccountSettings(&input)
	if err != nil {
		log.Printf("Error getting full concurrency of lambda")
	}
    fullConcurrency :=  int(*result.AccountLimit.ConcurrentExecutions)
    data := make(map[string]int)
    data["full_concurrency"] = fullConcurrency
    jsonData,err := json.Marshal(data)    
    if err != nil {
		log.Printf("error parsing data: %s", err)
        return "", nil, err
	}
    return string(jsonData), data, nil  
}

