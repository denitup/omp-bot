package recording

import "github.com/ozonmp/omp-bot/internal/model/streaming"

type DummyRecordingService struct{}

func NewDummyRecordingService() *DummyRecordingService {
	return &DummyRecordingService{}
}

func (s *DummyRecordingService) Describe(recordingID uint64) (*streaming.Recording, error) {
	return nil, nil
}

func (s *DummyRecordingService) List(cursor uint64, limit uint64) ([]streaming.Recording, error) {
	return nil, nil
}

func (s *DummyRecordingService) Create(streaming.Recording) (uint64, error) {
	return 0, nil
}

func (s *DummyRecordingService) Update(recordingID uint64, recording streaming.Recording) error {
	return nil
}

func (s *DummyRecordingService) Remove(recordingID uint64) (bool, error) {
	return false, nil
}
