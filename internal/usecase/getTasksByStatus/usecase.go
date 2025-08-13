package getTasksByStatus

import (
	"LOTestTask/internal/model"
	"LOTestTask/pkg"
	"LOTestTask/pkg/getTasksByStatusPkg"
	"LOTestTask/tools/logger"
	"encoding/json"
	"net/http"
)

const layer = "get_tasks_by_status_usecase"

type UseCase struct {
	getTasksByStatus getTasksByStatus
	logger           *logger.Logger
}

func New(getTasksByStatus getTasksByStatus, logger *logger.Logger) *UseCase {
	return &UseCase{
		getTasksByStatus: getTasksByStatus,
		logger:           logger,
	}
}

func (u *UseCase) Execute(w http.ResponseWriter, r *http.Request) {
	data := getTasksByStatusPkg.GetTasksByStatusV0Request{}

	statuses := r.URL.Query()["statuses"]

	data.Statuses = pkg.ConvertStringArrToStatuses(statuses)

	if err := data.Validate(); err != nil {
		u.logger.Warning(err, layer)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskModels := u.getTasksByStatus.GetTasksByStatus(pkg.ConvertStatusesToStringArr(data.Statuses))

	resp := model.ConvertModelsToGetTasksByStatusPkg(taskModels)

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
