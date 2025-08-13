package createTask

import (
	"LOTestTask/internal/model"
	"LOTestTask/pkg/createTaskPkg"
	"LOTestTask/tools/logger"
	"encoding/json"
	"io"
	"net/http"
)

const layer = "create_task_usecase"

type UseCase struct {
	createTask createTask
	logger     *logger.Logger
}

func New(createTask createTask, logger *logger.Logger) *UseCase {
	return &UseCase{
		createTask: createTask,
		logger:     logger,
	}
}

func (u *UseCase) Execute(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		u.logger.Error(err, layer)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	data := createTaskPkg.CreateTaskV0Request{}

	if err := json.Unmarshal(body, &data); err != nil {
		u.logger.Error(err, layer)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := data.Validate(); err != nil {
		u.logger.Warning(err, layer)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := u.createTask.CreateTask(model.ConvertCreateTaskPkgToModel(data)); err != nil {
		u.logger.Warning(err, layer)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := createTaskPkg.CreateTaskV0Response{ID: data.ID}

	respJson, err := json.Marshal(resp)
	if err != nil {
		u.logger.Error(err, layer)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respJson); err != nil {
		u.logger.Error(err, layer)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
