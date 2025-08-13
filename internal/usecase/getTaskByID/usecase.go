package getTaskByID

import (
	"LOTestTask/internal/model"
	"LOTestTask/pkg/getTaskByIDPkg"
	"LOTestTask/tools/logger"
	"encoding/json"
	"errors"
	"net/http"
	"path"
)

const layer = "get_task_by_id_usecase"

var ErrTaskNotFound = errors.New("task with this id not found")

type UseCase struct {
	getTaskByID getTaskByID
	logger      *logger.Logger
}

func New(getTaskByID getTaskByID, logger *logger.Logger) *UseCase {
	return &UseCase{
		getTaskByID: getTaskByID,
		logger:      logger,
	}
}

func (u *UseCase) Execute(w http.ResponseWriter, r *http.Request) {
	data := getTaskByIDPkg.GetTaskByIDV0Request{}

	data.ID = path.Base(r.URL.Path)

	if err := data.Validate(); err != nil {
		u.logger.Warning(err, layer)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskModel, ok := u.getTaskByID.GetTaskByID(data.ID)
	if !ok {
		u.logger.Warning(ErrTaskNotFound, layer)
		http.Error(w, ErrTaskNotFound.Error(), http.StatusNotFound)
		return
	}

	resp := model.ConvertModelToGetTaskByIDPkg(taskModel)

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
