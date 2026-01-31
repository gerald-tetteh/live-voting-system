package services

import (
	"database/sql"

	"geraldaddo.com/live-voting-system/models"
	"go.uber.org/zap"
)

type ElectionService struct {
	Logger *zap.Logger
	DB *sql.DB
}

type ElectionQueryParams struct {
	Status models.ElectionStatus
	Limit int
	Offset int
}

func (service *ElectionService) Save(election *models.Election) error {
	insertStatement := `
	INSERT INTO elections(title, description, start_time, end_time, status)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := service.DB.Exec(
		insertStatement, election.Title, election.Description, election.StartTime, election.EndTime, election.Status)
	if err != nil {
		return err
	}
	return nil
}

func (service *ElectionService) GetAll(params ElectionQueryParams) ([]models.Election, error) {
	query := `
	SELECT id, title, description, start_time, end_time, status, created_at, updated_at
	FROM elections
	WHERE $1 = '' OR status = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := service.DB.Query(query, params.Status, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elections []models.Election
	for rows.Next() {
		var e models.Election
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		elections = append(elections, e)
	}
	return elections, nil
}

func (service *ElectionService) GetOne(id string) (*models.Election, error) {
	query := `
	SELECT id, title, description, start_time, end_time, status, created_at, updated_at
	FROM elections
	WHERE id = $1
	`
	row := service.DB.QueryRow(query, id)
	
	var e models.Election
	err := row.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}