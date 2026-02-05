package election

import (
	"context"
	"database/sql"

	"geraldaddo.com/live-voting-system/platform/models"
)

type ElectionQueryParams struct {
	Status ElectionStatus
	Limit int
	Offset int
}

type ElectionRepository interface {
	models.Repository[Election]
	GetAllWithFilters(ctx context.Context, params ElectionQueryParams) ([]Election, error)
}
type ElectionRepositoryImpl struct {
	db *sql.DB
}

func NewElectionRepository(db *sql.DB) *ElectionRepositoryImpl {
	return &ElectionRepositoryImpl{db: db}
}

func (repo *ElectionRepositoryImpl) Save(ctx context.Context, election *Election) error {
	insertStatement := `
	INSERT INTO elections(title, description, start_time, end_time, status)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := repo.db.Exec(
		insertStatement, election.Title, election.Description, election.StartTime, election.EndTime, election.Status)
	return err
}

func (repo *ElectionRepositoryImpl) GetById(ctx context.Context, id string) (*Election, error) {
	query := `
	SELECT id, title, description, start_time, end_time, status, created_at, updated_at
	FROM elections
	WHERE id = $1
	`
	row := repo.db.QueryRow(query, id)
	
	var e Election
	err := row.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (repo *ElectionRepositoryImpl) GetAllWithFilters(ctx context.Context, params ElectionQueryParams) ([]Election, error) {
	query := `
	SELECT id, title, description, start_time, end_time, status, created_at, updated_at
	FROM elections
	WHERE $1 = '' OR status = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := repo.db.Query(query, params.Status, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elections []Election
	for rows.Next() {
		var e Election
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		elections = append(elections, e)
	}
	return elections, nil
}

func (repo *ElectionRepositoryImpl) UpdateOne(ctx context.Context, id string, e *Election) error {
	updateStatement := `
	UPDATE elections
	SET title = $1, description = $2, start_time = $3, end_time = $4, status = $5
	WHERE id = $6
	`
	_, err := repo.db.Exec(updateStatement, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, id)
	return err
}