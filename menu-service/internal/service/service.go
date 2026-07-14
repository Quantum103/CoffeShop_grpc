package service

import (
	"context"
	"errors"
	"fmt"

	models "github.com/Quantum103/menu-service/internal/model"
	"github.com/Quantum103/menu-service/internal/repository"
)

type MenuService interface {
	GetActiveMenu(ctx context.Context) ([]models.Item, error)
	GetItem(ctx context.Context, id string) (*models.Item, error) // <-- ДОБАВЛЕНО

}
type menuServiceImpl struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuServiceImpl{repo: repo}
}

// получаем меню и фильтруем только те, которые активные
func (s *menuServiceImpl) GetActiveMenu(ctx context.Context) ([]models.Item, error) {
	allItems, err := s.repo.GetAllItems(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch items from database")
	}
	var activeItems []models.Item
	for _, item := range allItems {
		if item.IsAvailable {
			activeItems = append(activeItems, item)
		}
	}
	return activeItems, nil

}
func (s *menuServiceImpl) GetItem(ctx context.Context, id string) (*models.Item, error) {
	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	// Дополнительная бизнес-логика: например, проверять, доступен ли товар
	if item != nil && !item.IsAvailable {
		return nil, fmt.Errorf("item is not available")
	}

	return item, nil
}
