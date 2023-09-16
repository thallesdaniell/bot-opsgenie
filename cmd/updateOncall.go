package cmd

import (
	"opslink/internal/environments"
	"opslink/internal/gcp"
	"opslink/internal/opsgenie"
	"opslink/internal/slack"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func init() {

	if os.Getenv("APP_ENV") == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&nested.Formatter{
			HideKeys:        true,
			TimestampFormat: time.RFC3339Nano,
		})
	}
}

func updateOncall() {

	config := environments.Key()

	clientOpsgenie, err := opsgenie.Client(config.OpsgenieApiKey)
	if err != nil {
		return
	}
	usersOpsgenie, err := clientOpsgenie.GetOnCalls(config.OpsgenieOnCall)
	if err != nil {
		return
	}
	clientGcp, err := gcp.Service(config.ServiceAccount, config.SubjectEmail)
	if err != nil {
		return
	}
	usersGcp, err := clientGcp.GetMembers(config.GroupKey)
	if err != nil {
		return
	}
	clientSlack, err := slack.Client(config.SlackApiKey)
	if err != nil {
		return
	}
	var usersSlack []string

	for _, user := range usersOpsgenie {
		if !slices.Contains(usersGcp, user) {
			clientGcp.InsertMember(config.GroupKey, user)
		}
		userByEmail, err := clientSlack.GetUserByEmail(user)
		if err != nil {
			return
		}
		usersSlack = append(usersSlack, userByEmail.ID)

	}

	for _, user := range usersGcp {
		if !slices.Contains(usersOpsgenie, user) {
			clientGcp.DeleteMember(config.GroupKey, user)
		}
	}
	clientSlack.UpdateUserGroupMembers(config.SlackGroupIdUpdateOncall, usersSlack)
}

var updateOncallCmd = &cobra.Command{
	Use:   "updateOncall",
	Short: "Atualiza grupos do slack e goole cloud",
	Long: `Executa o comando updateOncall checando os membros atuais do opsgenie
e compara com os usuários do slack e google cloud, caso estejam diferentes serão atualizados na gcp e no slack "`,
	Example: "opslink updateOncall",
	Run:     runUpdateOncall,
}

func runUpdateOncall(cmd *cobra.Command, args []string) {
	updateOncall()
}

func init() {
	rootCmd.AddCommand(updateOncallCmd)
}
