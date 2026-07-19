package repository

import (
	"context"
	"fmt"

	models "github.com/Quantum103/menu-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuRepository interface {
	GetAllItems(ctx context.Context) ([]models.Item, error)
	GetItemByID(ctx context.Context, id string) (*models.Item, error)
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// РЕАЛИЗАЦИЯ интерфейса для PostgreSQL
type pgxMenuRepository struct {
	db Querier
}

// Создаёт новый репозиторий с подключением к БД
func NewMenuRepository(db *pgxpool.Pool) MenuRepository {
	return &pgxMenuRepository{db: db}
}

func (r *pgxMenuRepository) GetAllItems(ctx context.Context) ([]models.Item, error) {
	query := `SELECT id, name, description, price, is_available FROM menu_items`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query menu items: %w", err)
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.IsAvailable,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating menu items: %w", err)
	}

	return items, nil
}

func (r *pgxMenuRepository) GetItemByID(ctx context.Context, id string) (*models.Item, error) {
	query := `SELECT id, name, description, price, is_available FROM menu_items WHERE id = $1`

	var item models.Item

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.IsAvailable,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get item by id %s: %w", id, err)
	}
	return &item, nil
}
