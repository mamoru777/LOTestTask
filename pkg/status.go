package pkg

import (
	"errors"
)

type Status string

const (
	Open           Status = "open"
	InProgress     Status = "in progress"
	CodeReview     Status = "code review"
	InTest         Status = "in test"
	ReadyToRelease Status = "ready to release"
	Done           Status = "done"
)

var ErrStatusUndefined = errors.New("status undefined")

func (s Status) Validate() error {
	switch s {
	case Open, InProgress, CodeReview, InTest, ReadyToRelease, Done:
		return nil
	default:
		return ErrStatusUndefined
	}
}

func ConvertStatusesToStringArr(statuses []Status) []string {
	resp := make([]string, 0, len(statuses))

	for _, status := range statuses {
		resp = append(resp, string(status))
	}

	return resp
}

func ConvertStringArrToStatuses(statuses []string) []Status {
	resp := make([]Status, 0, len(statuses))

	for _, status := range statuses {
		resp = append(resp, Status(status))
	}

	return resp
}
