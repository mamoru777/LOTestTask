package api

import (
	"LOTestTask/internal/usecase/createTask"
	"LOTestTask/internal/usecase/getTaskByID"
	"LOTestTask/internal/usecase/getTasksByStatus"
	"fmt"
	"net/http"
	"strings"
)

const MethodNotAllowed = "Method Not Allowed"

type Api struct {
	mux                     *http.ServeMux
	createTaskUseCase       *createTask.UseCase
	getTaskByIDUseCase      *getTaskByID.UseCase
	getTasksByStatusUseCase *getTasksByStatus.UseCase
}

func New(
	createTaskUseCase *createTask.UseCase,
	getTaskByIDUseCase *getTaskByID.UseCase,
	getTasksByStatusUseCase *getTasksByStatus.UseCase,
) *Api {
	return &Api{
		createTaskUseCase:       createTaskUseCase,
		getTaskByIDUseCase:      getTaskByIDUseCase,
		getTasksByStatusUseCase: getTasksByStatusUseCase,
		mux:                     http.NewServeMux(),
	}
}

func (a *Api) Route() {
	a.mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	a.mux.HandleFunc("/", a.handleTasks)
}

func (a *Api) handleTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request:", r.Method, r.URL.Path)

	if r.URL.Path == "/tasks" {
		switch r.Method {
		case http.MethodGet:
			a.getTasksByStatusUseCase.Execute(w, r)
			return
		case http.MethodPost:
			a.createTaskUseCase.Execute(w, r)
			return
		default:
			http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	}

	if strings.HasPrefix(r.URL.Path, "/tasks/") && r.Method == http.MethodGet {
		a.getTaskByIDUseCase.Execute(w, r)
		return
	}

	http.NotFound(w, r)
}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
