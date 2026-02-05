package election

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ElectionService struct {
	repo ElectionRepository
	log *zap.Logger
}

func NewElectionService(repo ElectionRepository, logger *zap.Logger) *ElectionService {
	return &ElectionService{repo: repo, log: logger}
}

func (service *ElectionService) CreateElection(ctx context.Context, election *Election) error {
	requestId, _ := ctx.Value("requestId").(string)
	election.Status = Draft
	if election.StartTime.After(election.EndTime) || election.StartTime.Equal(election.EndTime) {
		service.log.Warn("Election start time must be before end time")
		return errors.New("Election start time must be before end time")
	}
	err := service.repo.Save(ctx, election)
	if err != nil {
		service.log.Error(err.Error())
		service.log.Error("Could not create election", zap.String("request_id", requestId))
		return errors.New("Could not create election")
	}
	service.log.Info("Created election", zap.String("request_id", requestId))
	return nil
}

func (service *ElectionService) GetElections(ctx context.Context, params ElectionQueryParams) ([]Election, error) {
	requestId, _ := ctx.Value("requestId").(string)
	elections, err := service.repo.GetAllWithFilters(ctx, params)
	if err != nil {
		service.log.Error(err.Error())
		service.log.Error("Failed to get list of elections", zap.String("request_id", requestId))
		return nil, errors.New("Failed to get elections")
	}
	service.log.Info(fmt.Sprintf("Got elections of length: %d", len(elections)), zap.String("request_id", requestId))
	return elections, nil
}

func (service *ElectionService) GetElection(ctx context.Context, id string) (*Election, error) {
	requestId, _ := ctx.Value("requestId").(string)
	election, err := service.repo.GetById(ctx, id)
	if err != nil {
		service.log.Error(err.Error())
		service.log.Error("Could not find election with id: " + id, zap.String("request_id", requestId))
		return nil, errors.New("Failed to get election with ID: " + id)
	}
	service.log.Info("Found election with ID: " + id, zap.String("request_id", requestId))
	return election, nil
}

func (service *ElectionService) UpdateElection(ctx context.Context, id string, updatedElection *Election) error {
	requestId, _ := ctx.Value("requestId").(string)
	election, err := service.repo.GetById(ctx, id)
	if err != nil {
		service.log.Error(err.Error())
		service.log.Error("Election with id: " + id + " does not exist", zap.String("request_id", requestId))
		return errors.New("Election with ID: " + id + " does not exist")
	}
	now := time.Now()
	if election.Status != Draft || election.StartTime.Before(now) || election.StartTime.Equal(now) {
		service.log.Warn("Cannot update active or closed election: " + id, zap.String("request_id", requestId))
		return errors.New("Cannot update active or closed elections")
	}
	err = service.repo.UpdateOne(ctx, id, updatedElection)
	if err != nil {
		service.log.Error(err.Error())
		service.log.Error("Could not update election: " + id, zap.String("request_id", requestId))
		return errors.New("Could not update election: " + id)
	}
	service.log.Info("Updated election: " + id, zap.String("request_id", requestId))
	return nil
}