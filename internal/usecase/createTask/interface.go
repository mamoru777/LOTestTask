package createTask

import "LOTestTask/internal/model"

type createTask interface {
	CreateTask(task model.Task) error
}
