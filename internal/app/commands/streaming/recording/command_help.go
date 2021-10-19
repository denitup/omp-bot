package recording

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (c *RecordingCommander) CommandHelp(inputMessage *tgbotapi.Message) {
	helpMessage := `supported commands:
/help__streaming__recording - help
/get__streaming__recording {recording_id} - get recording
/list__streaming__recording {start_recording_id} {limit} {page_limit} - list recordings
/delete__streaming__recording {recording_id} - delete recording
/new__streaming__recording - create new recording
/edit__streaming__recording {recording_id} - edit recording`

	outputMessage := tgbotapi.NewMessage(inputMessage.Chat.ID, helpMessage)

	_, err := c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
