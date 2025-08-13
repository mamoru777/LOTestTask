package getTasksByStatusPkg

import (
	"LOTestTask/pkg"
	"errors"
	"time"
)

var (
	ErrEmptyID = errors.New("empty id")
)

type GetTasksByStatusV0Request struct {
	Statuses []pkg.Status `json:"statuses"`
}

type GetTasksByStatusV0Response struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	ID          string     `json:"id"`
	Discription string     `json:"discription"`
	AuthorID    string     `json:"authorID"`
	Status      pkg.Status `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (g *GetTasksByStatusV0Request) Validate() error {
	if len(g.Statuses) > 0 {
		for _, s := range g.Statuses {
			err := s.Validate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
