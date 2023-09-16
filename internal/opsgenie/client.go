package opsgenie

import (
	"strings"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/sirupsen/logrus"
)

type ScheduleClient struct {
	Client *schedule.Client
}

func Client(ApiKey string) (*ScheduleClient, error) {

	clientSchedule, err := schedule.NewClient(&client.Config{
		ApiKey:   ApiKey,
		LogLevel: logrus.ErrorLevel,
	})
	if err != nil {
		logrus.WithField("service", "opsgenie").Error("Erro criar client", err)
		return nil, err
	}
	return &ScheduleClient{
		Client: clientSchedule,
	}, err
}

func (c *ScheduleClient) GetOnCalls(schedules string) ([]string, error) {

	date := time.Now().AddDate(0, 0, 0)
	flat := true
	users := []string{}
	oncalls := strings.Split(schedules, ",")
	for _, oncall := range oncalls {
		result, err := c.Client.GetOnCalls(nil, &schedule.GetOnCallsRequest{
			Flat:                   &flat,
			Date:                   &date,
			ScheduleIdentifierType: schedule.Id,
			ScheduleIdentifier:     oncall,
		})
		if err != nil {
			logrus.WithField("service", "opsgenie").Error("Erro ao pegar on-calls", err)
			return nil, err
		}
		logrus.WithField("service", "opsgenie").Info(result.Parent.Name, ": ", result.OnCallRecipients[0])
		users = append(users, result.OnCallRecipients[0])
	}
	return users, nil
}
