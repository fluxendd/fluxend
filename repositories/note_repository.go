package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"myapp/errs"
	"myapp/models"
	"myapp/utils"
	"time"
)

type NoteRepository struct {
	db *sqlx.DB
}

func NewNoteRepository(injector *do.Injector) (*NoteRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)
	return &NoteRepository{db: db}, nil
}

func (r *NoteRepository) List(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Note, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := `
		SELECT 
			id, user_id, title, content, created_at, updated_at 
		FROM notes 
			WHERE user_id = :user_id
			ORDER BY :sort DESC
		LIMIT :limit 
		OFFSET :offset
	`

	params := map[string]interface{}{
		"user_id": authenticatedUserId,
		"sort":    paginationParams.Sort,
		"limit":   paginationParams.Limit,
		"offset":  offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.StructScan(&note); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return notes, nil
}

func (r *NoteRepository) GetByID(id, authenticatedUserId uint) (models.Note, error) {
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2"
	var note models.Note
	err := r.db.Get(&note, query, id, authenticatedUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Note{}, errs.NewNotFoundError("note.error.notFound")
		}

		return models.Note{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return note, nil
}

func (r *NoteRepository) ExistsByID(id, authenticatedUserId uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM notes WHERE id = $1 AND user_id = $2)"
	var exists bool
	err := r.db.Get(&exists, query, id, authenticatedUserId)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *NoteRepository) Create(note *models.Note) (*models.Note, error) {
	query := "INSERT INTO notes (title, content, user_id) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRowx(query, note.Title, note.Content, note.UserId).Scan(&note.ID)
	if err != nil {
		return &models.Note{}, fmt.Errorf("could not create row: %v", err)
	}

	return note, nil
}

func (r *NoteRepository) Update(id uint, note *models.Note) (*models.Note, error) {
	note.UpdatedAt = time.Now()
	note.ID = id

	query := `
		UPDATE notes 
		SET title = :title, content = :content, updated_at = :updated_at 
		WHERE id = :id`

	res, err := r.db.NamedExec(query, note)
	if err != nil {
		return &models.Note{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Note{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return note, nil
}

func (r *NoteRepository) Delete(noteId uint) (bool, error) {
	query := "DELETE FROM notes WHERE id = $1"
	res, err := r.db.Exec(query, noteId)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
