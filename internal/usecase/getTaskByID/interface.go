package getTaskByID

import "LOTestTask/internal/model"

type getTaskByID interface {
	GetTaskByID(id string) (model.Task, bool)
}
