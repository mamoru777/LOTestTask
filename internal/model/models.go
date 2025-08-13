package model

import (
	"LOTestTask/pkg"
	"LOTestTask/pkg/createTaskPkg"
	"LOTestTask/pkg/getTaskByIDPkg"
	"LOTestTask/pkg/getTasksByStatusPkg"
	"time"
)

type Task struct {
	ID          string
	Discription string
	AuthorID    string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ConvertCreateTaskPkgToModel(req createTaskPkg.CreateTaskV0Request) Task {
	return Task{
		ID:          req.ID,
		Discription: req.Discription,
		AuthorID:    req.AuthorID,
		Status:      string(req.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func ConvertModelToGetTaskByIDPkg(task Task) getTaskByIDPkg.GetTaskByIDV0Response {
	return getTaskByIDPkg.GetTaskByIDV0Response{
		ID:          task.ID,
		Discription: task.Discription,
		AuthorID:    task.AuthorID,
		Status:      pkg.Status(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func ConvertModelsToGetTasksByStatusPkg(tasks []Task) getTasksByStatusPkg.GetTasksByStatusV0Response {
	resp := getTasksByStatusPkg.GetTasksByStatusV0Response{}

	if len(tasks) == 0 {
		return resp
	}

	resp.Tasks = make([]getTasksByStatusPkg.Task, 0, len(tasks))

	for _, task := range tasks {
		resp.Tasks = append(resp.Tasks, getTasksByStatusPkg.Task{
			ID:          task.ID,
			Discription: task.Discription,
			AuthorID:    task.AuthorID,
			Status:      pkg.Status(task.Status),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	return resp
}
