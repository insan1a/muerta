package storagetype

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type StorageTypeRepositorer interface {
	repositories.Repository
	FindByID(ctx context.Context, id int) (models.StorageType, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.StorageType, error)
	Create(ctx context.Context, storageType models.StorageType) error
	Update(ctx context.Context, storageType models.StorageType) error
	Delete(ctx context.Context, id int) error
	FindTips(ctx context.Context, id int) ([]models.Tip, error)
	CreateTip(ctx context.Context, id int, tipID int) (models.Tip, error)
	DeleteTip(ctx context.Context, id int, tipID int) error
	FindStorages(ctx context.Context, id int) ([]models.Storage, error)
}

type storageTypeRepository struct {
	client repositories.PostgresClient
}

func (r *storageTypeRepository) Count(ctx context.Context) (int, error) {
	var (
		query = `
			SELECT COUNT(*) FROM storages_types
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count types: %w", err)
	}
	return count, nil
}

// CreateTip implements StorageTypeRepositorer
func (r *storageTypeRepository) CreateTip(ctx context.Context, id int, tipID int) (models.Tip, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO  storages_types_tips (id_storage_type, id_tip)
				VALUES ($1, $2)
				RETURNING id_tip, id_storage_type
			)
			SELECT t.id, t.description
			JOIN inserted i ON i.id_tip = t.id
			WHERE t.id = i.id_tip
			LIMIT 1
		`
		result models.Tip
	)
	if err := r.client.QueryRow(ctx, query, id, tipID).Scan(&result.ID, &result.Description); err != nil {
		return result, fmt.Errorf("failed to create tip: %w", err)
	}
	return result, nil
}

// DeleteTip implements StorageTypeRepositorer
func (r *storageTypeRepository) DeleteTip(ctx context.Context, id int, tipID int) error {
	var query = `
		DELETE FROM storages_types_tips
		WHERE id_storage_type = $1 AND id_tip = $2
	`
	if _, err := r.client.Exec(ctx, query, id, tipID); err != nil {
		return fmt.Errorf("failed to delete tip: %w", err)
	}
	return nil
}

// FindStorages implements StorageTypeRepositorer
func (r *storageTypeRepository) FindStorages(ctx context.Context, id int) ([]models.Storage, error) {
	var (
		query = `
			SELECT s.id, s.name, s.temperature, s.humidity
			FROM storages s
			WHERE s.id_type = $1 AND s.deleted_at IS NULL
		`
		result []models.Storage
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find storages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var storage models.Storage
		if err := rows.Scan(&storage.ID, &storage.Name, &storage.Temperature, &storage.Humidity); err != nil {
			return nil, fmt.Errorf("failed to scan storage: %w", err)
		}
		result = append(result, storage)
	}
	return result, nil
}

// FindTips implements StorageTypeRepositorer
func (r *storageTypeRepository) FindTips(ctx context.Context, id int) ([]models.Tip, error) {
	var (
		query = `
			SELECT t.id, t.description
			FROM tips t
			JOIN storages_types_tips stt ON stt.id_tip = t.id
			WHERE stt.id_storage_type = $1 AND t.deleted_at IS NULL
		`
		result []models.Tip
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tips: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var tip models.Tip
		if err := rows.Scan(&tip.ID, &tip.Description); err != nil {
			return nil, fmt.Errorf("failed to scan tip: %w", err)
		}
		result = append(result, tip)
	}
	return result, nil
}

// Create implements StorageTypeRepositorer
func (r *storageTypeRepository) Create(ctx context.Context, storageType models.StorageType) error {
	var (
		query = `
			INSERT INTO storages_types (name)
			VALUES ($1)
		`
	)
	if _, err := r.client.Exec(ctx, query, storageType.Name); err != nil {
		return fmt.Errorf("failed to create storageType: %w", err)
	}
	return nil
}

// Delete implements StorageTypeRepositorer
func (r *storageTypeRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			DELETE FROM storages_types
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete storageType: %w", err)
	}
	return nil
}

// FindByID implements StorageTypeRepositorer
func (r *storageTypeRepository) FindByID(ctx context.Context, id int) (models.StorageType, error) {
	var (
		query = `
			SELECT id, name
			FROM storages_types
			WHERE id = $1
			LIMIT 1	
		`
		storageType models.StorageType
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&storageType.ID, &storageType.Name); err != nil {
		return models.StorageType{}, fmt.Errorf("failed to find storageType: %w", err)
	}
	return storageType, nil
}

// FindMany implements StorageTypeRepositorer
func (r *storageTypeRepository) FindMany(ctx context.Context, limit int, offset int, name string) ([]models.StorageType, error) {
	var (
		query = `
			SELECT id, name
			FROM storages_types
			WHERE name ILIKE $1
			LIMIT $2
			OFFSET $3
		`
		storageTypes = make([]models.StorageType, 0, limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find storageTypes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var storageType models.StorageType
		if err := rows.Scan(&storageType.ID, &storageType.Name); err != nil {
			return nil, fmt.Errorf("failed to scan storageType: %w", err)
		}
		storageTypes = append(storageTypes, storageType)
	}
	return storageTypes, nil
}

// Update implements StorageTypeRepositorer
func (r *storageTypeRepository) Update(ctx context.Context, storageType models.StorageType) error {
	var (
		query = `
			UPDATE storages_types
			SET name = $1
			WHERE id = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, storageType.Name, storageType.ID); err != nil {
		return fmt.Errorf("failed to update storageType: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) StorageTypeRepositorer {
	return &storageTypeRepository{
		client: client,
	}
}
