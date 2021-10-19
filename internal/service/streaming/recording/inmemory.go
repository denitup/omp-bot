package recording

import (
	"errors"
	"github.com/ozonmp/omp-bot/internal/model/streaming"
	"log"
)

type recordingIdToIdxMap map[uint64]int

type InMemoryService struct {
	recordings       []*streaming.Recording
	recordingIdToIdx recordingIdToIdxMap
	nextId           uint64
}

const firstId = 1

var ErrRecordingNotFound = errors.New("recording not found")
var ErrNoMoreRecordingIDs = errors.New("all recording ids are already used")
var ErrInvalidStartRecordingID = errors.New("invalid value for starting id")

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		recordings:       []*streaming.Recording{},
		recordingIdToIdx: recordingIdToIdxMap{},
		nextId:           firstId,
	}
}

func NewFilledInMemoryService() *InMemoryService {
	service := &InMemoryService{
		recordings:       []*streaming.Recording{},
		recordingIdToIdx: recordingIdToIdxMap{},
		nextId:           firstId,
	}

	recordingTitles := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
	}

	for _, title := range recordingTitles {
		_, err := service.Create(*streaming.NewRecording(title))
		if err != nil {
			log.Fatal(err)
		}
	}

	return service
}

func (s *InMemoryService) getRecording(recordingID uint64) (*streaming.Recording, int) {
	if recordingIdx, hasRecording := s.recordingIdToIdx[recordingID]; hasRecording && s.recordings[recordingIdx] != nil {
		return s.recordings[recordingIdx], recordingIdx
	}

	return nil, 0
}

func (s *InMemoryService) Describe(recordingID uint64) (*streaming.Recording, error) {
	recording, _ := s.getRecording(recordingID)
	if recording == nil {
		return nil, ErrRecordingNotFound
	}

	recordingCopy := recording.Copy()

	return &recordingCopy, nil
}

func (s *InMemoryService) List(cursor uint64, limit uint64) ([]streaming.Recording, error) {
	if cursor < firstId {
		return nil, ErrInvalidStartRecordingID
	}

	if limit == 0 {
		return []streaming.Recording{}, nil
	}

	recordings := make([]streaming.Recording, 0, limit)
	var recordingsCount uint64

	for _, recording := range s.recordings {
		if recording == nil {
			continue
		}

		if recording.ID >= cursor {
			recordings = append(recordings, recording.Copy())
			recordingsCount++

			if recordingsCount >= limit {
				break
			}
		}
	}

	return recordings, nil
}

func (s *InMemoryService) Create(newRecording streaming.Recording) (uint64, error) {
	if s.nextId+1 == 0 {
		return 0, ErrNoMoreRecordingIDs
	}

	newRecording.ID = s.nextId
	s.nextId++

	s.recordings = append(s.recordings, &newRecording)
	s.recordingIdToIdx[newRecording.ID] = len(s.recordings) - 1

	return newRecording.ID, nil
}

func (s *InMemoryService) Update(recordingID uint64, update streaming.Recording) error {
	recording, recordingIdx := s.getRecording(recordingID)
	if recording == nil {
		return ErrRecordingNotFound
	}

	update.ID = recording.ID
	s.recordings[recordingIdx] = &update

	return nil
}

func (s *InMemoryService) Remove(recordingID uint64) (bool, error) {
	recording, recordingIdx := s.getRecording(recordingID)
	if recording == nil {
		return false, nil
	}

	s.recordings[recordingIdx] = nil
	delete(s.recordingIdToIdx, recordingID)

	return true, nil
}
