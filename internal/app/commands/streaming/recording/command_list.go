package recording

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/path"
	"github.com/ozonmp/omp-bot/internal/service/streaming/recording"
	"log"
	"strings"
)

const listPageLimit = 5

type listResponse struct {
	StartRecordingId uint64
	Limit            uint64
	PageLimit        uint64
	Output           *string
}

type firstPageResponse struct {
	Response *listResponse
	JSON     []byte
}

func (c *RecordingCommander) getListResponse(startRecordingId, limit, pageLimit uint64) (*listResponse, error) {
	getLimit := pageLimit
	if getLimit > limit {
		getLimit = limit
	}

	recordings, err := c.recordingService.List(startRecordingId, getLimit)
	if err != nil {
		if err == recording.ErrInvalidStartRecordingID {
			return nil, err
		}

		log.Printf(
			"failed to list recordings with start_recording_id (%d) and limit (%d): %v",
			startRecordingId,
			limit,
			err,
		)

		return nil, errors.New("failed to list recordings")
	}

	recordingsCount := len(recordings)

	if recordingsCount == 0 {
		return nil, errors.New("no recordings found")
	}

	buf := strings.Builder{}

	for idx, listedRecording := range recordings {
		buf.WriteString(listedRecording.String())

		if idx+1 != recordingsCount {
			buf.WriteRune('\n')
		}
	}

	output := buf.String()

	leftLimit := limit
	if uint64(recordingsCount) < getLimit {
		leftLimit = 0
	} else {
		leftLimit -= getLimit
	}

	return &listResponse{
		StartRecordingId: startRecordingId + getLimit,
		Limit:            leftLimit,
		PageLimit:        pageLimit,
		Output:           &output,
	}, nil
}

func (c *RecordingCommander) handleListRecordings(args []string) (*firstPageResponse, error) {
	argsCount := len(args)
	if argsCount < 2 {
		return nil, errors.New("not enough arguments for command")
	} else if argsCount > 3 {
		return nil, errors.New("too many arguments for command")
	}

	startRecordingId, err := c.getSingleUInt64Arg(args[0])
	if err != nil {
		log.Printf("invalid value for {start_recording_id} argument: %v", err)

		return nil, errors.New("invalid value for {start_recording_id} argument")
	}

	limit, err := c.getSingleUInt64Arg(args[1])
	if err != nil {
		log.Printf("invalid value for {limit} argument: %v", err)

		return nil, errors.New("invalid value for {limit} argument")
	}

	var pageLimit uint64 = listPageLimit

	if argsCount > 2 {
		pageLimit, err = c.getSingleUInt64Arg(args[2])
		if err != nil {
			log.Printf("invalid value for {page_limit} argument: %v", err)

			return nil, errors.New("invalid value for {page_limit} argument")
		}
	}

	response, err := c.getListResponse(startRecordingId, limit, pageLimit)
	if err != nil {
		return nil, err
	}

	var responseJson []byte

	if response != nil && response.Limit > 0 {
		responseJson, err = json.Marshal(response)
		if err != nil {
			log.Printf(
				"failed to list recordings: failed marshalling callback data (%v): %v",
				response,
				err,
			)

			return nil, errors.New("failed to list recordings")
		}
	}

	return &firstPageResponse{
		Response: response,
		JSON:     responseJson,
	}, nil
}

func (c *RecordingCommander) CommandList(inputMessage *tgbotapi.Message) {
	var listMessage *string
	var callbackData *path.CallbackPath

	firstResponse, err := c.handleListRecordings(strings.Split(inputMessage.CommandArguments(), " "))
	if err != nil {
		errMessage := err.Error()
		listMessage = &errMessage
	} else {
		listMessage = firstResponse.Response.Output
		if firstResponse.JSON != nil {
			// according to the telegram docs:
			// https://core.telegram.org/bots/api#inlinekeyboardbutton
			// callbackData is limited to 64 bytes max
			// three uint64 stored as text will take up to 60 bytes,
			// and we also need space for domain, subdomain, callbackName and separators
			// so in order to fit in we need  more compact string method marshaller for callbackData,
			// but unfortunately we're not supposed to change it
			// therefore the only thing we can do is to rely on the custom state data

			c.state.SetCallbackData(inputMessage.Chat.ID, firstResponse.JSON)

			callbackData = &path.CallbackPath{
				Domain:       "streaming",
				Subdomain:    "recording",
				CallbackName: "list",
				CallbackData: "",
			}
		}
	}

	outputMessage := tgbotapi.NewMessage(inputMessage.Chat.ID, *listMessage)

	if callbackData != nil {
		outputMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next page", callbackData.String()),
			),
		)
	}

	_, err = c.bot.Send(outputMessage)
	if err != nil {
		log.Printf("error sending reply message (%v) to chat: %v", outputMessage, err)
	}
}
