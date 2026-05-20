package device

import (
	"context"
	"database/sql"
	"errors"

	"mertani/internal/shared/id"
)

var ErrNotFound = errors.New("device not found")

type Repository interface {
	Create(ctx context.Context, device *Device) error
	FindAll(ctx context.Context) ([]Device, error)
	FindByID(ctx context.Context, id id.ID) (Device, error)
	Update(ctx context.Context, device *Device) error
	Delete(ctx context.Context, id id.ID) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, device *Device) error {
	query := `
		INSERT INTO devices (id, name, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		device.ID,
		device.Name,
		device.Location,
		device.CreatedAt,
		device.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]Device, error) {
	query := `
		SELECT id, name, location, created_at, updated_at
		FROM devices
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := make([]Device, 0)
	for rows.Next() {
		var device Device
		if err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Location,
			&device.CreatedAt,
			&device.UpdatedAt,
		); err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id id.ID) (Device, error) {
	query := `
		SELECT id, name, location, created_at, updated_at
		FROM devices
		WHERE id = $1
	`

	var device Device
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&device.ID,
		&device.Name,
		&device.Location,
		&device.CreatedAt,
		&device.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Device{}, ErrNotFound
	}
	if err != nil {
		return Device{}, err
	}

	return device, nil
}

func (r *PostgresRepository) Update(ctx context.Context, device *Device) error {
	query := `
		UPDATE devices
		SET name = $2,
			location = $3,
			updated_at = $4
		WHERE id = $1
		RETURNING id, name, location, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		device.ID,
		device.Name,
		device.Location,
		device.UpdatedAt,
	).Scan(
		&device.ID,
		&device.Name,
		&device.Location,
		&device.CreatedAt,
		&device.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id id.ID) error {
	query := `DELETE FROM devices WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
