package streaming

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	recordingCommander "github.com/ozonmp/omp-bot/internal/app/commands/streaming/recording"
	"github.com/ozonmp/omp-bot/internal/app/path"
	recordingService "github.com/ozonmp/omp-bot/internal/service/streaming/recording"
	"log"
)

type Commander interface {
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath)
}

type StreamingCommander struct {
	bot                *tgbotapi.BotAPI
	recordingCommander Commander
}

func NewStreamingCommander(
	bot *tgbotapi.BotAPI,
) *StreamingCommander {
	return &StreamingCommander{
		bot:                bot,
		recordingCommander: recordingCommander.NewRecordingCommander(bot, recordingService.NewFilledInMemoryService()),
	}
}

func (c *StreamingCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.Subdomain {
	case "recording":
		c.recordingCommander.HandleCallback(callback, callbackPath)
	default:
		log.Printf("StreamingCommander.HandleCallback: unknown subdomain - %s", callbackPath.Subdomain)
	}
}

func (c *StreamingCommander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.Subdomain {
	case "recording":
		c.recordingCommander.HandleCommand(msg, commandPath)
	default:
		log.Printf("StreamingCommander.HandleCommand: unknown subdomain - %s", commandPath.Subdomain)
	}
}
