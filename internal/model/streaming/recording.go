package streaming

import "fmt"

type Recording struct {
	ID    uint64
	Title string
}

func NewRecording(title string) *Recording {
	return &Recording{
		Title: title,
	}
}

func (r *Recording) String() string {
	return fmt.Sprintf("Recording{ID: %d, Title: %s}", r.ID, r.Title)
}

func (r *Recording) Copy() Recording {
	return Recording{
		ID:    r.ID,
		Title: r.Title,
	}
}
