package recording

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (c *RecordingCommander) handleDeleteRecording(arg string) (bool, error) {
	recordingID, err := c.getRecordingIDArgument(arg)
	if err != nil {
		return false, err
	}

	isDeleted, err := c.recordingService.Remove(recordingID)
	if err != nil {
		log.Printf("failed to delete recording with id (%d): %v", recordingID, err)
		return false, errors.New("failed to delete recording")
	}

	return isDeleted, nil
}

func (c *RecordingCommander) CommandDelete(inputMessage *tgbotapi.Message) {
	var deleteMessage string

	isDeleted, err := c.handleDeleteRecording(inputMessage.CommandArguments())
	if err != nil {
		deleteMessage = err.Error()
	} else if !isDeleted {
		deleteMessage = "record is not deleted.\ndid you specify existing recording_id?"
	} else {
		deleteMessage = "record has been deleted"
	}

	outputMessage := tgbotapi.NewMessage(inputMessage.Chat.ID, deleteMessage)

	_, err = c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
