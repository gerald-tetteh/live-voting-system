package election_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"geraldaddo.com/live-voting-system/domain/election"
	"geraldaddo.com/live-voting-system/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func SetupServer() *gin.Engine {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("requestId", uuid.New().String())
	})
	return server
}
func SetupTestAPI(ctrl *gomock.Controller) (*election.ElectionAPI, *mocks.MockElectionRepository) {
	repository := mocks.NewMockElectionRepository(ctrl)
	service := election.NewElectionService(repository, zap.NewNop())
	return election.NewElectionAPI(service, zap.NewNop()), repository
}

func TestCreateElectionAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := SetupServer()
	api, mockRepo := SetupTestAPI(ctrl)
	api.RegisterRoutes(server)

	validElection := election.Election{
		Title: "valid election",
		Description: "valid test election",
		StartTime: time.Now(),
		EndTime: time.Now().Add(time.Hour),
		Status: election.Draft,
	}
	invalidElection := election.Election{
		Title: "invalid election",
	}

	tests := []struct {
		name string
		input election.Election
		output string
		statusCode int
		shouldFail bool
	}{
		{"successfully create election", validElection, "created election", 200, false},
		{"fail to create election", invalidElection, "could not parse election", 400, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.shouldFail {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}
			recorder := httptest.NewRecorder()
			electionJson, _ := json.Marshal(test.input)
			request, _ := http.NewRequest("POST", "/elections", strings.NewReader(string(electionJson)))
			server.ServeHTTP(recorder, request)
			if recorder.Code != test.statusCode {
				t.Errorf("Expected error code: %d but got %d", test.statusCode, recorder.Code)
			}
			var response map[string]string
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			if err != nil {
				t.Fatal("Request did not return valid JSON")
			}
			message, exists := response["message"]
			if !exists {
				t.Fatal("JSON did not contain message key")
			}
			if message != test.output {
				t.Errorf("Expected message: %s but got %s", test.output, message)
			}
		})
	}
}

func TestGetElectionsAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := SetupServer()
	api, mockRepo := SetupTestAPI(ctrl)
	api.RegisterRoutes(server)

	expectedElections := []election.Election {
		{Title: "test-election-1"},
		{Title: "test-election-2"},
		{Title: "test-election-3"},
	}
	expectedBytes, _ := json.Marshal(expectedElections)
	expectedJson := string(expectedBytes)

	tests := []struct {
		name string
		status int
		query string
		expected string
		shouldFail bool
	}{
		{"Fail to parse status", 400, "status=invalid", "status is invalid", true},
		{"Fail to parse page number", 400, "page=t", "could not parse page number", true},
		{"Fail to parse page size", 400, "size=t", "could not parse page size", true},
		{"Successfully get elections", 200, "", expectedJson, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.shouldFail {
				mockRepo.
					EXPECT().
					GetAllWithFilters(gomock.Any(), gomock.Any()).
					Return(expectedElections, nil).
					Times(1)
			}
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/elections?" + test.query, nil)
			server.ServeHTTP(recorder, request)
			if recorder.Code != test.status {
				t.Errorf("Expect status code: %d but got %d", test.status, recorder.Code)
			}
			if test.shouldFail {
				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal("Request did not return valid JSON")
				}
				message, exists := response["message"]
				if !exists {
					t.Fatal("JSON does not contain message key")
				}
				if message != test.expected {
					t.Errorf("Expected message: %s but got %s", test.expected, message)
				}
			} else if recorder.Body.String() != test.expected {
				t.Errorf("JSON response: %s did not match expected: %s", recorder.Body.String(), test.expected)
			}
		})
	}
}

func TestGetElectionAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := SetupServer()
	api, mockRepo := SetupTestAPI(ctrl)
	api.RegisterRoutes(server)

	result := &election.Election{Title: "test-election"}
	resultBytes, _ := json.Marshal(result)
	mockRepo.
		EXPECT().
		GetById(gomock.Any(), gomock.Any()).
		Return(result, nil).
		Times(1)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/elections/test-request-id", nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("Expected status code: %d but got %d", 200, recorder.Code)
	}
	if !slices.Equal(recorder.Body.Bytes(), resultBytes) {
		t.Errorf("Request did not return expected JSON")
	}
}

func TestUpdateElectionAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := SetupServer()
	api, mockRepo := SetupTestAPI(ctrl)
	api.RegisterRoutes(server)

	now := time.Now()
	existingElection := &election.Election{
		Title: "test-election",
		Description: "test election description",
		StartTime: now.Add(time.Hour),
		EndTime: now.Add(2 * time.Hour),
		Status: election.Draft,
	}
	validElection := &election.Election{
		Title: "test-election",
		Description: "test election description",
		StartTime: now,
		EndTime: now.Add(time.Hour),
		Status: election.Draft,
	}
	invalidElection := &election.Election{Title: "invalid election"}

	validElectionBytes, _ := json.Marshal(validElection)
	invalidElectionBytes, _ := json.Marshal(invalidElection)

	tests := []struct {
		name string
		input []byte
		status int
		result string
		shouldFail bool
	}{
		{"Fail to parse election update", invalidElectionBytes, 400, "could not parse update information", true},
		{"Complete election update", validElectionBytes, 200, "updated election", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.shouldFail {
				mockRepo.
					EXPECT().
					GetById(gomock.Any(), gomock.Any()).
					Return(existingElection, nil).
					Times(1)
				mockRepo.
					EXPECT().
					UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			}
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("PATCH", "/elections/test-id", strings.NewReader(string(test.input)))
			server.ServeHTTP(recorder, request)

			if recorder.Code != test.status {
				t.Errorf("Expect status code: %d but got %d", test.status, recorder.Code)
			}
			var response map[string]string
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			if err != nil {
				t.Fatal("Request did not return valid JSON")
			}
			message, exists := response["message"]
			if !exists {
				t.Fatal("JSON does not contain message key")
			}
			if message != test.result {
				t.Errorf("Expected message: %s but got %s", test.result, message)
			}
		})
	}
}