package recording

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/model/streaming"
	"github.com/ozonmp/omp-bot/internal/service/streaming/recording"
	"log"
	"strings"
)

func (c *RecordingCommander) handleNewRecording(title string) (uint64, error) {
	if strings.TrimSpace(title) == "" {
		return 0, errors.New("enter non empty title")
	}

	recordingID, err := c.recordingService.Create(streaming.Recording{
		Title: title,
	})
	if err != nil {
		log.Printf("failed to create new recording: %v", err)

		if err == recording.ErrNoMoreRecordingIDs {
			return 0, errors.New("no more recordings can be added because no ids are available")
		}

		return 0, errors.New("failed to create new recording")
	}

	return recordingID, nil
}

func (c *RecordingCommander) CommandNew(inputMessage *tgbotapi.Message) {
	var newMessage string
	chatID := inputMessage.Chat.ID

	stateID, hasState := c.state.GetState(chatID)
	if hasState && stateID == recordingCommanderNewRecordingStateId {
		recordingID, err := c.handleNewRecording(inputMessage.CommandArguments())
		if err != nil {
			newMessage = err.Error()
		} else {
			c.state.Reset(chatID)

			newMessage = fmt.Sprintf("successfully created new recording with ID = %d", recordingID)
		}
	} else {
		c.state.SetState(chatID, recordingCommanderNewRecordingStateId)

		newMessage = "enter recording title field as an argument for this command"
	}

	outputMessage := tgbotapi.NewMessage(inputMessage.Chat.ID, newMessage)

	_, err := c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
