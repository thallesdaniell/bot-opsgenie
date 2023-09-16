package environments

import (
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Environments struct {
	OpsgenieApiKey           string `env:"OPSGENIE_API_KEY"`
	OpsgenieOnCall           string `env:"OPSGENIE_ON_CALL_SCHEDULE"`
	OpsGenieOncallPrimary    string `env:"OPSGENIE_ON_CALL_SCHEDULE_PRIMARY"`
	ServiceAccount           string `env:"GOOGLE_APPLICATION_CREDENTIALS"`
	SubjectEmail             string `env:"GOOGLE_SUBJECT_EMAIL"`
	GroupKey                 string `env:"GOOGLE_GROUP_KEY"`
	SlackGroupIdUpdateOncall string `env:"SLACK_GROUP_ID_UPDATE_ONCALL"`
	SlackGroupIdNextOncall   string `env:"SLACK_GROUP_ID_NEXT_ONCALL"`
	SlackApiKey              string `env:"SLACK_BOT_TOKEN"`
}

func Key() Environments {
	envs := Environments{}
	err := env.Parse(&envs)
	if err != nil {
		logrus.Fatalf("Erro ao fazer o parse: %s", err.Error())
	}
	return envs
}
