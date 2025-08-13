package di

import (
	"LOTestTask/cfg"
	"LOTestTask/internal/api"
	"LOTestTask/internal/repository/taskRepo"
	"LOTestTask/internal/usecase/createTask"
	"LOTestTask/internal/usecase/getTaskByID"
	"LOTestTask/internal/usecase/getTasksByStatus"
	"LOTestTask/tools/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const layer = "di"

type DI struct {
	config *cfg.Config

	infr struct {
		logger *logger.Logger
	}

	repos struct {
		taskRepo *taskRepo.Repository
	}

	useCases struct {
		createTask       *createTask.UseCase
		getTaskByID      *getTaskByID.UseCase
		getTasksByStatus *getTasksByStatus.UseCase
	}

	api *api.Api

	server *http.Server
}

func (di *DI) Init() error {
	if err := di.loadCfg(); err != nil {
		return err
	}

	di.initLogger()
	di.initRepos()
	di.initUseCases()
	di.initApi()
	di.initServer()

	return nil
}

func (di *DI) loadCfg() error {
	di.config = &cfg.Config{}

	return di.config.LoadConfig()
}

func (di *DI) initLogger() {
	di.infr.logger = logger.InitLogger(di.config.LoggerCfg, os.Stdout)
}

func (di *DI) initRepos() {
	di.repos.taskRepo = taskRepo.New(di.config.CacheCfg, di.infr.logger)
}

func (di *DI) initUseCases() {
	di.useCases.createTask = createTask.New(di.repos.taskRepo, di.infr.logger)

	di.useCases.getTaskByID = getTaskByID.New(di.repos.taskRepo, di.infr.logger)

	di.useCases.getTasksByStatus = getTasksByStatus.New(di.repos.taskRepo, di.infr.logger)
}

func (di *DI) initApi() {
	di.api = api.New(
		di.useCases.createTask,
		di.useCases.getTaskByID,
		di.useCases.getTasksByStatus,
	)

	di.api.Route()
}

func (di *DI) initServer() {
	di.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", di.config.ServerCfg.Addr, di.config.ServerCfg.Port),
		Handler: di.api,
	}
}

func (di *DI) Start() error {
	di.infr.logger.Info(fmt.Sprintf("starting server on port - %s", di.config.ServerCfg.Port), layer)

	err := di.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (di *DI) Stop(ctx context.Context) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopCh)

	select {
	case <-stopCh:
	case <-ctx.Done():
	}

	di.infr.logger.Info("stoping server", layer)

	if err := di.server.Shutdown(ctx); err != nil {
		di.infr.logger.Error(fmt.Sprintf("server shutdown error: %v", err), layer)
	}

	di.infr.logger.Close()
}
