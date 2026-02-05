package election_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"geraldaddo.com/live-voting-system/domain/election"
	"geraldaddo.com/live-voting-system/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateElection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockElectionRepository := mocks.NewMockElectionRepository(ctrl)

	now := time.Now()
	input := &election.Election{
		Title: "test",
		Description: "test election",
		StartTime: now,
		EndTime: now.Add(time.Hour),
		Status: election.Draft,
	}

	mockElectionRepository.
		EXPECT().
		Save(gomock.Any(), input).
		Return(nil).
		Times(1)
	service := election.NewElectionService(mockElectionRepository, zap.NewNop())
	ctx := context.WithValue(context.Background(), "requestId", "test-request-id")
	err := service.CreateElection(ctx, input)

	if err != nil {
		t.Error("Create election returned an error", err.Error())
	}
}

func TestCreateElectionShouldFailIfStartTimeIsNotBeforeEndTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	tests := []struct {
		name string
		startTime time.Time
		endTime time.Time
	}{
		{"endTime before startTime", now, now.Add(-1 * time.Hour)},
		{"endTime equal to startTime", now, now},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockElectionRepository := mocks.NewMockElectionRepository(ctrl)
			input := &election.Election{StartTime: test.startTime, EndTime: test.endTime}
			service := election.NewElectionService(mockElectionRepository, zap.NewNop())
			ctx := context.WithValue(context.Background(), "requestId", "test-request-id")
			err := service.CreateElection(ctx, input)
			if err == nil {
				t.Fatal("Should create election where start date is after end date")
			}
			if err.Error() != "Election start time must be before end time" {
				t.Error("Did not return expected exception")
			}
		})
	}
}

func TestGetElections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockElectionRepository := mocks.NewMockElectionRepository(ctrl)

	elections := []election.Election {
		{Title: "test-1"},
		{Title: "test-2"},
		{Title: "test-3"},
	}
	queryParams := election.ElectionQueryParams{}

	mockElectionRepository.
		EXPECT().
		GetAllWithFilters(gomock.Any(), gomock.Any()).
		Return(elections, nil).
		Times(1)

	service := election.NewElectionService(mockElectionRepository, zap.NewNop())
	ctx := context.WithValue(context.Background(), "requestId", "test-request-id")
	result, err := service.GetElections(ctx, queryParams)

	if err != nil {
		t.Error("Could not get list of elections", err.Error())
	}
	if !slices.Equal(result, elections) {
		t.Error("Did not return expected elections")
	}
}

func TestGetElection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockElectionRepository := mocks.NewMockElectionRepository(ctrl)

	electionId := "test-election-id"
	expected := &election.Election{ID: electionId}

	mockElectionRepository.
		EXPECT().
		GetById(gomock.Any(), electionId).
		Return(expected, nil).
		Times(1)

	service := election.NewElectionService(mockElectionRepository, zap.NewNop())
	ctx := context.WithValue(context.Background(), "requestId", "test-request-id")
	result, err := service.GetElection(ctx, electionId)

	if err != nil {
		t.Error("Could not get election", err.Error())
	}
	if result != expected {
		t.Error("Did not return expected election")
	}
}

func TestUpdateElection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockElectionRepository := mocks.NewMockElectionRepository(ctrl)

	electionId := "test-election-id"
	now := time.Now()
	existingElection := &election.Election{
		ID: electionId, 
		Title: "Existing", 
		StartTime: now.Add(2 * time.Hour), 
		EndTime: now.Add(3 * time.Hour),
		Status: election.Draft,
	}
	updatedElection := &election.Election{ID: electionId, Title: "Updated"}

	mockElectionRepository.
		EXPECT().
		GetById(gomock.Any(), electionId).
		Return(existingElection, nil).
		Times(1)
	mockElectionRepository.
		EXPECT().
		UpdateOne(gomock.Any(), electionId, updatedElection).
		Return(nil).
		Times(1)
	
	service := election.NewElectionService(mockElectionRepository, zap.NewNop())
	ctx := context.WithValue(context.Background(), "requestId", "test-request-id")
	err := service.UpdateElection(ctx, electionId, updatedElection)

	if err != nil {
		t.Error("Could not update election", err.Error())
	}
}

func TestUpdateElectionShouldNotUpdateActiveOrPastElection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		status election.ElectionStatus
		startTime time.Time
		endTime time.Time
	}{
		{"Active election", election.Active, time.Now(), time.Now().Add(time.Hour)},
		{"Closed election", election.Closed, time.Now(), time.Now().Add(time.Hour)},
		{"Archived election", election.Archived, time.Now(), time.Now().Add(time.Hour)},
		{"Draft election that has already started", election.Draft, time.Now().Add(-1 * time.Hour), time.Now()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			existingElection := &election.Election{
				Status: test.status, 
				StartTime: test.startTime, 
				EndTime: test.endTime,
			}
			mockElectionRepository := mocks.NewMockElectionRepository(ctrl)
			service := election.NewElectionService(mockElectionRepository, zap.NewNop())
			ctx := context.WithValue(context.Background(), "requestId", "test-request-id")
			mockElectionRepository.
				EXPECT().
				GetById(gomock.Any(), gomock.Any()).
				Return(existingElection, nil).
				Times(1)
			err := service.UpdateElection(ctx, "test-id", &election.Election{})
			if err == nil {
				t.Error("Attempted to update an active or closed election")
			}
		})
	}
}