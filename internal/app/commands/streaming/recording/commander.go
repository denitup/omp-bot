package recording

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/path"
	"github.com/ozonmp/omp-bot/internal/model/streaming"
	"log"
)

type RecordingService interface {
	Describe(recordingID uint64) (*streaming.Recording, error)
	List(cursor uint64, limit uint64) ([]streaming.Recording, error)
	Create(streaming.Recording) (uint64, error)
	Update(recordingID uint64, recording streaming.Recording) error
	Remove(recordingID uint64) (bool, error)
}

type RecordingCommander struct {
	bot              *tgbotapi.BotAPI
	recordingService RecordingService
	state            *RecordingCommanderState
}

func NewRecordingCommander(bot *tgbotapi.BotAPI, service RecordingService) *RecordingCommander {
	return &RecordingCommander{
		bot:              bot,
		recordingService: service,
		state:            newRecordingCommanderState(),
	}
}

func (c *RecordingCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case "list":
		c.CallbackList(callback, callbackPath)
	default:
		log.Printf("unknown callback name: %s", callbackPath.CallbackName)
	}
}

func (c *RecordingCommander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	chatID := msg.Chat.ID

	switch commandPath.CommandName {
	case "help":
		c.CommandHelp(msg)
		c.state.Reset(chatID)
	case "get":
		c.CommandGet(msg)
		c.state.Reset(chatID)
	case "list":
		c.CommandList(msg)
	case "delete":
		c.CommandDelete(msg)
		c.state.Reset(chatID)
	case "new":
		c.CommandNew(msg)
	case "edit":
		c.CommandEdit(msg)
	default:
		c.Default(msg)
	}
}

func (c *RecordingCommander) Default(inputMessage *tgbotapi.Message) {
	outputMessage := tgbotapi.NewMessage(inputMessage.Chat.ID, "unknown command")

	_, err := c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
