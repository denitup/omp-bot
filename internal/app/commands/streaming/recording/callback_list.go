package recording

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/path"
	"log"
)

type nextPageResponse struct {
	Output      *string
	JSON        []byte
	HasMoreData bool
}

func (c *RecordingCommander) handleGetNextPage(callbackData []byte) (*nextPageResponse, error) {
	response := &listResponse{}

	err := json.Unmarshal(callbackData, &response)
	if err != nil {
		log.Printf("failed to list recording's next page: %v", err)

		return nil, errors.New("failed to list recording's next page")
	}

	response, err = c.getListResponse(response.StartRecordingId, response.Limit, response.PageLimit)
	if err != nil {
		return nil, err
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to list recording's next page: %v", err)

		return nil, errors.New("failed to list recording's next page")
	}

	hasMoreData := response.Limit > 0

	return &nextPageResponse{
		Output:      response.Output,
		JSON:        responseJson,
		HasMoreData: hasMoreData,
	}, nil
}

func (c *RecordingCommander) CallbackList(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	var callbackMessage *string
	chatID := callback.Message.Chat.ID
	hasMoreData := false

	callbackJsonData, hasCallbackData := c.state.GetCallbackData(callback.Message.Chat.ID)
	if !hasCallbackData {
		errMessage := "request expired, repeat original command"
		callbackMessage = &errMessage
	} else {
		nextResponse, err := c.handleGetNextPage(callbackJsonData)
		if err != nil {
			errMessage := err.Error()
			callbackMessage = &errMessage
		} else {
			c.state.SetCallbackData(chatID, nextResponse.JSON)

			callbackMessage = nextResponse.Output

			hasMoreData = nextResponse.HasMoreData

			if !hasMoreData {
				c.state.Reset(chatID)
			}
		}
	}

	outputMessage := tgbotapi.NewMessage(chatID, *callbackMessage)

	if hasMoreData {
		outputMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next page", callbackPath.String()),
			),
		)
	}

	_, err := c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
