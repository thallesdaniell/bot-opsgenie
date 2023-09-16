package cmd

import (
	"encoding/json"
	"fmt"
	"opslink/internal/environments"
	"opslink/internal/slack"
	"strings"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Recipient struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Period struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Type      string    `json:"type"`
	Recipient Recipient `json:"recipient"`
}

type Rotation struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Order   float64  `json:"order"`
	Periods []Period `json:"periods"`
}

type ScheduleResponse struct {
	Rotations []Rotation `json:"rotations"`
}

func abv(t time.Time) string {
	diaSemana := t.Weekday().String()
	switch diaSemana {
	case "Sunday":
		return "Dom"
	case "Monday":
		return "Seg"
	case "Tuesday":
		return "Ter"
	case "Wednesday":
		return "Qua"
	case "Thursday":
		return "Qui"
	case "Friday":
		return "Sex"
	case "Saturday":
		return "Sáb"
	default:
		return ""
	}
}

func nextOncall() {

	config := environments.Key()

	clientSlack, err := slack.Client(config.SlackApiKey)
	if err != nil {
		return
	}
	scheduleClient, err := schedule.NewClient(&client.Config{
		//LogLevel: logrus.DebugLevel,
		ApiKey: config.OpsgenieApiKey,
	})

	if err != nil {
		logrus.WithField("service", "slack").Error(err)
	}

	loc, err := time.LoadLocation("America/Recife")
	if err != nil {
		logrus.WithField("cmd", "nextOncall").Errorln("Erro ao carregar o fuso horário de Recife:", err)
		return
	}

	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)
	//date := time.Date(2023, 8, 7, 12, 0, 0, 0, time.UTC)

	scheduleTimeline, err := scheduleClient.GetTimeline(nil, &schedule.GetTimelineRequest{
		IdentifierType:  schedule.Id,
		IdentifierValue: config.OpsGenieOncallPrimary,
		Expands:         []schedule.ExpandType{schedule.Base, schedule.Forwarding, schedule.Override},
		Interval:        7,
		IntervalUnit:    schedule.Days,
		Date:            &date,
	})

	jsonData, err := json.Marshal(scheduleTimeline.FinalTimeline)
	if err != nil {
		logrus.WithField("cmd", "nextOncall").Error("Erro ao gerar JSON:", err)
		return
	}
	var response ScheduleResponse
	err = json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		logrus.WithField("cmd", "nextOncall").Error("Erro ao decodificar JSON:", err)
		return
	}
	total := len(response.Rotations)

	var message strings.Builder
	message.WriteString("```")
	for i, rotation := range response.Rotations {
		message.WriteString(fmt.Sprintf("Rotação: %s\n", rotation.Name))
		for _, period := range rotation.Periods {
			start := period.StartDate.In(loc)
			end := period.EndDate.In(loc)
			message.WriteString(fmt.Sprintf(" %-3s | %-3s | %-3s\n",
				fmt.Sprintf("%-3s %s", abv(start), start.Format("15:04")),
				fmt.Sprintf("%-3s %s", abv(end), end.Format("15:04")),
				period.Recipient.Name))
		}
		if i < total-1 {
			message.WriteString("\n")
		}
	}
	message.WriteString("```\n")

	text := "*O <https://rdstation.app.opsgenie.com/schedule/whoIsOnCall|link> de onde foram retiradas as informações acima.*\n" +
		"*Caso precise fazer override, combine antecipadamente com seu pair.*"

	title := slack.CreateSectionBlock("Rotação | Início | Fim | Nome", "header", "plain_text")
	body := slack.CreateSectionBlock(message.String(), "section", "mrkdwn")
	footer := slack.CreateSectionBlock(text, "context", "mrkdwn")

	blocks := slack.ConvertToSlackBlocks([]slack.MySlackBlock{title, body, footer})
	_, _, _, err = clientSlack.Client.SendMessage(config.SlackGroupIdNextOncall, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		logrus.WithField("cmd", "nextOncall").Errorf("Erro ao enviar a mensagem para o Slack: %s\n", err)
		return
	}
}

var nextOncallCmd = &cobra.Command{
	Use:   "nextOncall",
	Short: "Retorna usuários da escala da semana",
	Long: `Exibe em uma timeline usuários que está na escala de plantão. 
Essa timeline também ira exibir os overrides que exisitem, entçao pode 
ocorrer de apacerer mais de um linha do mesmo dia.`,
	Example: "opslink nextOncall",
	Run:     runNextOncall,
}

func runNextOncall(cmd *cobra.Command, args []string) {
	nextOncall()
}

func init() {
	rootCmd.AddCommand(nextOncallCmd)
}
