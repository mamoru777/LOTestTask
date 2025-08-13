package getTaskByIDPkg

import (
	"LOTestTask/pkg"
	"errors"
	"time"
)

var (
	ErrEmptyID = errors.New("empty id")
)

type GetTaskByIDV0Request struct {
	ID string `json:"id"`
}

type GetTaskByIDV0Response struct {
	ID          string     `json:"id"`
	Discription string     `json:"discription"`
	AuthorID    string     `json:"authorID"`
	Status      pkg.Status `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (g *GetTaskByIDV0Request) Validate() error {
	if g.ID == "" {
		return ErrEmptyID
	}

	return nil
}
