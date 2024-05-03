package comman_function

import (
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"log"
	"time"
)

func GetCloudWatchAlarms(clientAuth *model.Auth, startTime, endTime *time.Time) ([]*cloudwatch.MetricAlarm, error) {
	svc := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	// Call DescribeAlarms to get all alarms
	resp, err := svc.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{})
	if err != nil {
		log.Println("Error describing alarms:", err)
		return nil, err
	}

	// Filter alarms based on their last updated time within the specified range
	filteredAlarms := make([]*cloudwatch.MetricAlarm, 0)
	for _, alarm := range resp.MetricAlarms {
		stateUpdatedTime := alarm.StateUpdatedTimestamp
		if stateUpdatedTime != nil && stateUpdatedTime.After(*startTime) && stateUpdatedTime.Before(*endTime) {
			filteredAlarms = append(filteredAlarms, alarm)
		}
	}

	return filteredAlarms, nil
}
