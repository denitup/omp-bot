package recording

import "github.com/ozonmp/omp-bot/internal/model/streaming"

const recordingCommanderDefaultStateId = 0
const recordingCommanderNewRecordingStateId = 1
const recordingCommanderEditRecordingStateId = 2

type stateByChat map[int64]int
type recordingByChat map[int64]*streaming.Recording
type callbackDataByChat map[int64][]byte

type RecordingCommanderState struct {
	stateId      stateByChat
	recording    recordingByChat
	callbackData callbackDataByChat
}

func newRecordingCommanderState() *RecordingCommanderState {
	return &RecordingCommanderState{
		stateId:      stateByChat{},
		recording:    recordingByChat{},
		callbackData: callbackDataByChat{},
	}
}

func (s *RecordingCommanderState) Reset(chatID int64) {
	s.stateId[chatID] = recordingCommanderDefaultStateId
	delete(s.recording, chatID)
	delete(s.callbackData, chatID)
}

func (s *RecordingCommanderState) SetState(chatID int64, stateId int) {
	s.stateId[chatID] = stateId
}

func (s *RecordingCommanderState) GetState(chatID int64) (int, bool) {
	stateId, hasState := s.stateId[chatID]

	return stateId, hasState
}

func (s *RecordingCommanderState) SetRecording(chatID int64, data *streaming.Recording) {
	s.recording[chatID] = data
}

func (s *RecordingCommanderState) GetRecording(chatID int64) (*streaming.Recording, bool) {
	recording, hasRecording := s.recording[chatID]

	return recording, hasRecording
}

func (s *RecordingCommanderState) SetCallbackData(chatID int64, data []byte) {
	s.callbackData[chatID] = data
}

func (s *RecordingCommanderState) GetCallbackData(chatID int64) ([]byte, bool) {
	callbackData, hasCallbackData := s.callbackData[chatID]

	return callbackData, hasCallbackData
}
