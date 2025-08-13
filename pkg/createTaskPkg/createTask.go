package createTaskPkg

import (
	"LOTestTask/pkg"
	"errors"
)

var (
	ErrEmptyAuthorID = errors.New("empty author id")
	ErrEmptyID       = errors.New("empty id")
)

type CreateTaskV0Request struct {
	ID          string     `json:"id"`
	Discription string     `json:"discription"`
	AuthorID    string     `json:"authorID"`
	Status      pkg.Status `json:"status"`
}

type CreateTaskV0Response struct {
	ID string `json:"id"`
}

func (c *CreateTaskV0Request) Validate() error {
	if c.ID == "" {
		return ErrEmptyID
	}

	if c.AuthorID == "" {
		return ErrEmptyAuthorID
	}

	return c.Status.Validate()
}
