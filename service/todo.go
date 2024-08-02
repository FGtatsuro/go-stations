package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var rSubject, rDescription string
	var rCreatedAt, rUpdatedAt time.Time
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(
		&rSubject,
		&rDescription,
		&rCreatedAt,
		&rUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     rSubject,
		Description: rDescription,
		CreatedAt:   rCreatedAt,
		UpdatedAt:   rUpdatedAt,
	}
	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var rows *sql.Rows
	var err error
	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*model.TODO, 0)
	for rows.Next() {
		var rID int64
		var rSubject, rDescription string
		var rCreatedAt, rUpdatedAt time.Time
		if err := rows.Scan(
			&rID,
			&rSubject,
			&rDescription,
			&rCreatedAt,
			&rUpdatedAt,
		); err != nil {
			return nil, err
		}

		todo := &model.TODO{
			ID:          rID,
			Subject:     rSubject,
			Description: rDescription,
			CreatedAt:   rCreatedAt,
			UpdatedAt:   rUpdatedAt,
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return nil, &model.ErrNotFound{}
	}

	var rSubject, rDescription string
	var rCreatedAt, rUpdatedAt time.Time
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(
		&rSubject,
		&rDescription,
		&rCreatedAt,
		&rUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     rSubject,
		Description: rDescription,
		CreatedAt:   rCreatedAt,
		UpdatedAt:   rUpdatedAt,
	}
	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) == 0 {
		return nil
	}

	stmt := fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1))

	var args []interface{}
	for _, v := range ids {
		args = append(args, v)
	}
	res, err := s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
