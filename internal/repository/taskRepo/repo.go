package taskRepo

import (
	"LOTestTask/internal/model"
	"LOTestTask/tools/cache"
	"LOTestTask/tools/logger"
	"errors"
	"fmt"
)

const (
	layer             = "repository"
	taskAlreadyExists = "task already exist with this id"
)

type Repository struct {
	logger *logger.Logger
	cache  *cache.Cache[string, model.Task]
}

func New(cacheCfg cache.Config, logger *logger.Logger) *Repository {
	return &Repository{
		cache:  cache.NewCache[string, model.Task](cacheCfg),
		logger: logger,
	}
}

func (r *Repository) CreateTask(task model.Task) error {
	_, ok := r.cache.Get(task.ID)
	if ok {
		r.logger.Warning(fmt.Sprintf("%s - %s", taskAlreadyExists, task.ID), layer)
		return errors.New(taskAlreadyExists)
	}

	r.cache.Set(task.ID, task)
	r.logger.Info(fmt.Sprintf("task was set with id - %s", task.ID), layer)

	return nil
}

func (r *Repository) GetTaskByID(id string) (model.Task, bool) {
	task, ok := r.cache.Get(id)
	if ok {
		r.logger.Info(fmt.Sprintf("task got with id - %s", id), layer)
	} else {
		r.logger.Warning(fmt.Sprintf("task not found with id - %s", id), layer)
	}

	return task, ok
}

func (r *Repository) GetTasksByStatus(statuses []string) []model.Task {
	mp := r.cache.GetAll()

	resultTasks := make([]model.Task, 0, len(mp))

	for _, task := range mp {
		if len(statuses) == 0 {
			resultTasks = append(resultTasks, task)
			continue
		}

		for _, status := range statuses {
			if task.Status == status {
				resultTasks = append(resultTasks, task)
				break
			}
		}
	}

	if len(resultTasks) > 0 {
		r.logger.Info(fmt.Sprintf("tasks got with statuses - %v", statuses), layer)
	}

	return resultTasks
}
