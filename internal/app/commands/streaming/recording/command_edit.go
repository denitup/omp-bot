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

func (c *RecordingCommander) handleEditRecording(recordingID uint64, title string) (bool, error) {
	if strings.TrimSpace(title) == "" {
		return false, errors.New("enter non empty title")
	}

	err := c.recordingService.Update(recordingID, streaming.Recording{
		Title: title,
	})
	if err != nil {
		log.Printf("failed to update recording with id (%d): %v", recordingID, err)

		if err == recording.ErrRecordingNotFound {
			return true, errors.New("recording not found")
		}

		return false, errors.New("failed to edit recording")
	}

	return true, nil
}

func (c *RecordingCommander) CommandEdit(inputMessage *tgbotapi.Message) {
	var editMessage string
	chatID := inputMessage.Chat.ID

	stateID, hasState := c.state.GetState(chatID)
	if hasState && stateID == recordingCommanderEditRecordingStateId {
		recording, hasRecording := c.state.GetRecording(chatID)
		if !hasRecording {
			editMessage = "request expired, run edit command again"
		} else {
			isEdited, err := c.handleEditRecording(recording.ID, inputMessage.CommandArguments())
			if isEdited {
				c.state.Reset(chatID)
			}
			if err != nil {
				editMessage = err.Error()
			} else {
				editMessage = fmt.Sprintf("successfully edited recording with ID = %d", recording.ID)
			}
		}
	} else {
		editRecording, err := c.handleGetRecording(inputMessage.CommandArguments())
		if err != nil {
			editMessage = err.Error()
		} else {
			c.state.SetState(chatID, recordingCommanderEditRecordingStateId)
			c.state.SetRecording(chatID, editRecording)

			editMessage = fmt.Sprintf(
				"current title: %s\nenter new value for the \"title\" field as an argument for this command",
				editRecording.Title,
			)
		}
	}

	outputMessage := tgbotapi.NewMessage(inputMessage.Chat.ID, editMessage)

	_, err := c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
