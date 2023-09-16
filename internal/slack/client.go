package slack

import (
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackClient struct {
	Client *slack.Client
}

func Client(ApiKey string) (*SlackClient, error) {

	if ApiKey == "" {
		message := "Empty SLACK_BOT_TOKEN"
		logrus.WithField("service", "slack").Error(message)
		return nil, errors.New(message)
	}

	client := slack.New(ApiKey)
	return &SlackClient{
		Client: client,
	}, nil
}

func (client *SlackClient) GetUserByEmail(userEmail string) (*slack.User, error) {

	user, err := client.Client.GetUserByEmail(userEmail)
	if err != nil {
		logrus.WithField("service", "slack").Error(err)
		return nil, err
	}
	logrus.WithField("service", "slack").Info("Get user:" + user.Name + " " + user.ID)
	return user, nil
}

func (c *SlackClient) UpdateUserGroupMembers(groupId string, userIds []string) {

	additional := []string{os.Getenv("SLACK_ADDITIONAL_USERS_GROUP")}
	userIds = append(userIds, additional...)
	users := strings.Join(userIds, ",")

	group, err := c.Client.UpdateUserGroupMembers(groupId, ""+users+"")
	if err != nil {
		logrus.WithField("service", "slack").Error(err)
	}

	logrus.WithField("service", "slack").Info("Upate members of group:" + group.Name)
}

func MsgOptionBlocks(blocks ...slack.Block) slack.MsgOption {
	return slack.MsgOptionBlocks(blocks...)
}

type MySlackBlock struct {
	BlockType   slack.MessageBlockType
	ElementType slack.MessageObjectType
	Text        string
}

func CreateSectionBlock(text string, blockType slack.MessageBlockType, elementType slack.MessageObjectType) MySlackBlock {
	return MySlackBlock{
		blockType,
		elementType,
		text,
	}
}

func ConvertToSlackBlocks(blocks []MySlackBlock) []slack.Block {
	slackBlocks := []slack.Block{}
	for _, block := range blocks {

		textBlock := slack.NewTextBlockObject(string(block.ElementType), block.Text, false, false)
		switch block.BlockType {
		case slack.MBTSection:
			sectionBlock := slack.NewSectionBlock(textBlock, nil, nil)
			slackBlocks = append(slackBlocks, sectionBlock)
		case slack.MBTHeader:
			headerBlock := slack.NewHeaderBlock(textBlock)
			slackBlocks = append(slackBlocks, headerBlock)
		case slack.MBTContext:
			sectionBlock := slack.NewContextBlock("", textBlock)
			slackBlocks = append(slackBlocks, sectionBlock)
		default:
			// ...
		}
	}
	return slackBlocks
}
