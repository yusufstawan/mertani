package sensor

import (
	"context"
	"database/sql"
	"errors"

	"mertani/internal/shared/id"
)

var ErrNotFound = errors.New("sensor not found")

type Repository interface {
	Create(ctx context.Context, sensor *Sensor) error
	FindAll(ctx context.Context) ([]Sensor, error)
	FindByID(ctx context.Context, id id.ID) (Sensor, error)
	Update(ctx context.Context, sensor *Sensor) error
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

func (r *PostgresRepository) Create(ctx context.Context, sensor *Sensor) error {
	query := `
		INSERT INTO sensors (id, device_id, sensor_type, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		sensor.ID,
		sensor.DeviceID,
		sensor.SensorType,
		sensor.Value,
		sensor.CreatedAt,
		sensor.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]Sensor, error) {
	query := `
		SELECT id, device_id, sensor_type, value, created_at, updated_at
		FROM sensors
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sensors := make([]Sensor, 0)
	for rows.Next() {
		var sensor Sensor
		if err := rows.Scan(
			&sensor.ID,
			&sensor.DeviceID,
			&sensor.SensorType,
			&sensor.Value,
			&sensor.CreatedAt,
			&sensor.UpdatedAt,
		); err != nil {
			return nil, err
		}

		sensors = append(sensors, sensor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id id.ID) (Sensor, error) {
	query := `
		SELECT id, device_id, sensor_type, value, created_at, updated_at
		FROM sensors
		WHERE id = $1
	`

	var sensor Sensor
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sensor.ID,
		&sensor.DeviceID,
		&sensor.SensorType,
		&sensor.Value,
		&sensor.CreatedAt,
		&sensor.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Sensor{}, ErrNotFound
	}
	if err != nil {
		return Sensor{}, err
	}

	return sensor, nil
}

func (r *PostgresRepository) Update(ctx context.Context, sensor *Sensor) error {
	query := `
		UPDATE sensors
		SET device_id = $2,
			sensor_type = $3,
			value = $4,
			updated_at = $5
		WHERE id = $1
		RETURNING id, device_id, sensor_type, value, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		sensor.ID,
		sensor.DeviceID,
		sensor.SensorType,
		sensor.Value,
		sensor.UpdatedAt,
	).Scan(
		&sensor.ID,
		&sensor.DeviceID,
		&sensor.SensorType,
		&sensor.Value,
		&sensor.CreatedAt,
		&sensor.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id id.ID) error {
	query := `DELETE FROM sensors WHERE id = $1`

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
