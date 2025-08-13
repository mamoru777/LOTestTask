package getTaskByIDTest

import (
	"LOTestTask/internal/di"
	"LOTestTask/pkg"
	"LOTestTask/pkg/createTaskPkg"
	"LOTestTask/pkg/getTaskByIDPkg"
	"LOTestTask/tests/integration/usecases"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestGetTaskByID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	initInfr(ctx, cancel)

	testCases := initTestCases()

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			if tc.ExpectedOutput.ID != "" {
				prepare(tc.ExpectedOutput.ID, tc.ExpectedOutput.Discription, tc.ExpectedOutput.AuthorID, tc.ExpectedOutput.Status)
			}

			resp, err := http.Get(fmt.Sprintf("http://localhost:9191/tasks/%s", tc.TestData.ID))
			if err != tc.ExpectedError {
				t.Errorf("wrong error, expected - %v, actual - %v", tc.ExpectedError, err)
			}

			if resp.StatusCode != tc.ExpectedStatusCode {
				t.Errorf("wrong status code, expected - %d, actual - %d", tc.ExpectedStatusCode, resp.StatusCode)
			}

			if resp.StatusCode == 200 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("can not read body - %v", err)
				}

				data := getTaskByIDPkg.GetTaskByIDV0Response{}

				if err := json.Unmarshal(body, &data); err != nil {
					t.Errorf("can not unmarshal json - %v", err)
				}

				if data.ID != tc.ExpectedOutput.ID {
					t.Errorf("wrong output id, expected - %s, actual - %s", tc.ExpectedOutput.ID, data.ID)
				}

				if data.Discription != tc.ExpectedOutput.Discription {
					t.Errorf("wrong output discription, expected - %s, actual - %s", tc.ExpectedOutput.Discription, data.Discription)
				}

				if data.AuthorID != tc.ExpectedOutput.AuthorID {
					t.Errorf("wrong output author id, expected - %s, actual - %s", tc.ExpectedOutput.AuthorID, data.AuthorID)
				}

				if string(data.Status) != string(tc.ExpectedOutput.Status) {
					t.Errorf("wrong output status, expected - %s, actual - %s", string(tc.ExpectedOutput.Status), string(data.Status))
				}

				if !(data.CreatedAt.Sub(tc.ExpectedOutput.CreatedAt) < 3*time.Second) || !(tc.ExpectedOutput.CreatedAt.Sub(data.CreatedAt) < 3*time.Second) {
					t.Errorf("wrong output created at, expected - %s, actual - %s", tc.ExpectedOutput.CreatedAt, data.CreatedAt)
				}

				if !(data.UpdatedAt.Sub(tc.ExpectedOutput.UpdatedAt) < 3*time.Second) || !(tc.ExpectedOutput.UpdatedAt.Sub(data.UpdatedAt) < 3*time.Second) {
					t.Errorf("wrong output updated at, expected - %s, actual - %s", tc.ExpectedOutput.UpdatedAt, data.UpdatedAt)
				}
			}

			resp.Body.Close()
		})
	}

	cancel()
}

func initInfr(ctx context.Context, cancel context.CancelFunc) {
	os.Setenv("CACHE_SPACE", "30")
	os.Setenv("SERVER_ADDR", "localhost")
	os.Setenv("SERVER_PORT", "9191")
	os.Setenv("LOGGER_LOG_LEVEL", "info")

	di := di.DI{}
	if err := di.Init(); err != nil {
		log.Fatalf("can not init service - %v", err)
	}

	go func() {
		if err := di.Start(); err != nil {
			log.Printf("error ocured while starting server - %v", err)
			cancel()
		}
	}()

	go func() {
		di.Stop(ctx)
	}()
}

func prepare(id, discription, authorID string, status pkg.Status) error {
	data := createTaskPkg.CreateTaskV0Request{
		ID:          id,
		Discription: discription,
		AuthorID:    authorID,
		Status:      status,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:9191/tasks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("status code not 200")
	}

	return nil
}

func initTestCases() []usecases.GetTaskByIDTestCase {
	cases := []usecases.GetTaskByIDTestCase{
		{
			TestName: "[TestGetTaskByID] successfull get",
			TestData: getTaskByIDPkg.GetTaskByIDV0Request{
				ID: "1",
			},
			ExpectedError: nil,
			ExpectedOutput: getTaskByIDPkg.GetTaskByIDV0Response{
				ID:          "1",
				Discription: "test",
				AuthorID:    "11",
				Status:      "open",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			ExpectedStatusCode: 200,
		},
		{
			TestName: "[TestGetTaskByID] not found",
			TestData: getTaskByIDPkg.GetTaskByIDV0Request{
				ID: "2",
			},
			ExpectedError:      nil,
			ExpectedStatusCode: 404,
		},
	}

	return cases
}
