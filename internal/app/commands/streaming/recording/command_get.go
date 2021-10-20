package recording

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/model/streaming"
	recordingService "github.com/ozonmp/omp-bot/internal/service/streaming/recording"
	"log"
)

func (c *RecordingCommander) handleGetRecording(arg string) (*streaming.Recording, error) {
	recordingID, err := c.getRecordingIDArgument(arg)
	if err != nil {
		return nil, err
	}

	recording, err := c.recordingService.Describe(recordingID)
	if err != nil {
		if err == recordingService.ErrRecordingNotFound {
			return nil, fmt.Errorf("did not found recording for ID = %d", recordingID)
		}

		log.Printf("failed to get recording with id (%d): %v", recordingID, err)

		return nil, errors.New("failed to get recording")
	}

	return recording, nil
}

func (c *RecordingCommander) CommandGet(inputMessage *tgbotapi.Message) {
	var getMessage string

	recording, err := c.handleGetRecording(inputMessage.CommandArguments())
	if err != nil {
		getMessage = err.Error()
	} else {
		getMessage = fmt.Sprintf("found recording: %s", recording)
	}

	outputMessage := tgbotapi.NewMessage(inputMessage.Chat.ID, getMessage)

	_, err = c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
