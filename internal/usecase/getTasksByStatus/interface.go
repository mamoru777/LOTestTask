package getTasksByStatus

import "LOTestTask/internal/model"

type getTasksByStatus interface {
	GetTasksByStatus(statuses []string) []model.Task
}
