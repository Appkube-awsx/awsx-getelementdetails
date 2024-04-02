package RDS

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/Appkube-awsx/awsx-common/config"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
// 	"github.com/spf13/cobra"
// )

// var AwsxEKSRecentEventsCmd = &cobra.Command{

// 	Use:   "recent_events_panel",
// 	Short: "get recent events panel data for Amazon EKS",
// 	Long:  `Command to retrieve recent events panel data for Amazon EKS`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Running from child command")
// 		var authFlag bool
// 		var clientAuth *model.Auth
// 		var err error
// 		authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)

// 		if err != nil {

// 			log.Printf("Error during authentication: %v\n", err)

// 			err := cmd.Help()

// 			if err != nil {

// 				return
// 			}

// 			return
// 		}
// 		if authFlag {
// 		  err := GetEKSRecentEventsPanel(cmd)
// 		  if err != nil {
// 			return
// 			log.Println("Error:", err)
// 		  }
// 		}
// 	},
// }

// func GetEKSRecentEventsPanel(cmd *cobra.Command) error {
// 	elementId, _ := cmd.PersistentFlags().GetString("elementId")
// 	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
// 	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
// 	logGroupName := fmt.Sprintf("/aws/eks/%s/cluster", clusterName)
// 	if elementId != "" {
// 		log.Println("getting cloud-element data from cmdb")
// 		apiUrl := cmdbApiUrl
// 		if cmdbApiUrl == "" {
// 			log.Println("using default cmdb url")
// 			apiUrl = config.CmdbUrl
// 		}
// 		log.Println("cmdb url: " + apiUrl)
// 		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		logGroupName = cmdbData.LogGroup

// 	}

// 	endTime := time.Now()
// 	startTime := endTime.Add(-1 * time.Hour) // Retrieve events for the past hour
// 	// Parse start time if provided
// 	if startTimeStr != "" {
// 		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing start time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {

// 			}
// 		}
// 		startTime = &parsedStartTime
// 	} else {
// 		defaultStartTime := time.Now().Add(-5 * time.Minute)
// 		startTime = &defaultStartTime
// 	}

// 	if endTimeStr != "" {
// 		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing end time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				// handle error
// 			}
// 		}
// 		endTime = &parsedEndTime
// 	} else {
// 		defaultEndTime := time.Now()
// 		endTime = &defaultEndTime
// 	}
// 	results, err := filterCloudWatchLogs(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
// 	if err != nil {
// 		return nil, nil
// 	}
// 	processedResults := processQueryResultss(results)

// 	return processedResults, nil
// }

// 	// Construct input parameters
// 	params := &cloudwatchlogs.FilterLogEventsInput{
// 		LogGroupName: aws.String(logGroupName),
// 		StartTime:    aws.Int64(startTime.Unix() * 1000),
// 		EndTime:      aws.Int64(endTime.Unix() * 1000),
// 	}

// 	resp, err := cloudWatchLogs.FilterLogEvents(params)
// 	if err != nil {
// 		return fmt.Errorf("failed to retrieve log events: %v", err)
// 	}

// 	// Process and display the log events
// 	for _, event := range resp.Events {
// 		log.Printf("Timestamp: %v, Message: %s\n", *event.Timestamp, *event.Message)
// 	}

// 	return nil
// }

// func init() {
// 	AwsxEKSRecentEventsCmd.PersistentFlags().String("clusterName", "", "Amazon EKS cluster name")
// }
