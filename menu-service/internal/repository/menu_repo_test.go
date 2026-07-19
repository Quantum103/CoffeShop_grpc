package repository

import (
	"context"
	"errors"
	"strings"
	"testing"

	models "github.com/Quantum103/menu-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
)

func NewPgxMenuRepository(db Querier) *pgxMenuRepository {
	return &pgxMenuRepository{db: db}
}

func TestGetAllItems(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer mock.Close()
	repo := &pgxMenuRepository{db: mock}
	mockRows := pgxmock.NewRows([]string{"id", "name", "description", "price", "is_available"}).
		AddRow("1", "Ежь", "Описание", 200.00, true).
		AddRow("2", "жук", "описание", 200.00, false)

	mock.ExpectQuery("SELECT id, name, description, price, is_available FROM menu_items").WillReturnRows(mockRows)

	items, err := repo.GetAllItems(context.Background())
	if err != nil {
		t.Fatalf("ожидали что ошибки не будет, получили: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("ожидали 2 эл, получили: %d", len(items))
	}
	if items[0].Name != "Ежь" {
		t.Errorf("Ожидали 'Ежь', получили: %s", items[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("невыполненные ожидания мока: %v", err)
	}
}

func FallTest_GetAllItems(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("ошибка созд мока %v", err)
	}
	defer mock.Close()
	repo := &pgxMenuRepository{db: mock}

	mockRows := pgxmock.NewRows([]string{"id", "name", "description", "price", "is_available"}).
		AddRow(1, "а", "выа", "вав", "вав")

	mock.ExpectQuery("SELECT id, name, description, price, is_available FROM menu_items").WillReturnRows(mockRows)

	items, err := repo.GetAllItems(context.Background())

	if len(items) != 1 {
		t.Fatal("ошибка, ожидалось %d, получили %d", 1, len(items))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("невыполненные ожидания мока: %v", err)
	}
}

func TestPgxMenuRepository_GetItemByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer mock.Close()

	repo := &pgxMenuRepository{db: mock}

	queryRegex := "SELECT id, name, description, price, is_available FROM menu_items WHERE id = \\$1"

	tests := []struct {
		name          string
		id            string
		mockSetup     func()
		expectedItem  *models.Item
		expectError   bool
		errorMsgCheck string
	}{
		{
			name: "Success: item found",
			id:   "1",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"id", "name", "description", "price", "is_available"}).
					AddRow("1", "Бургер", "Сочный бургер", 250.50, true)

				mock.ExpectQuery(queryRegex).WithArgs("1").WillReturnRows(rows)
			},
			expectedItem: &models.Item{ID: "1", Name: "Бургер", Description: "Сочный бургер", Price: 250.50, IsAvailable: true},
			expectError:  false,
		},
		{
			name: "Error: item not found (pgx.ErrNoRows)",
			id:   "999",
			mockSetup: func() {
				mock.ExpectQuery(queryRegex).WithArgs("999").WillReturnError(pgx.ErrNoRows)
			},
			expectedItem:  nil,
			expectError:   true,
			errorMsgCheck: "failed to get item by id 999",
		},
		{
			name: "Error: database connection failed",
			id:   "1",
			mockSetup: func() {
				mock.ExpectQuery(queryRegex).WithArgs("1").WillReturnError(errors.New("connection refused"))
			},
			expectedItem:  nil,
			expectError:   true,
			errorMsgCheck: "failed to get item by id 1",
		},
		{
			name: "Error: scan failed (wrong data type)",
			id:   "1",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"id", "name", "description", "price", "is_available"}).
					AddRow("1", "Бургер", "Сочный", "not_a_number", true)

				mock.ExpectQuery(queryRegex).WithArgs("1").WillReturnRows(rows)
			},
			expectedItem:  nil,
			expectError:   true,
			errorMsgCheck: "failed to get item by id 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			item, err := repo.GetItemByID(context.Background(), tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("GetItemByID() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError {
				if !strings.Contains(err.Error(), tt.errorMsgCheck) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsgCheck, err.Error())
				}

				if tt.name == "Error: item not found (pgx.ErrNoRows)" {
					if !errors.Is(err, pgx.ErrNoRows) {
						t.Errorf("Expected error to wrap pgx.ErrNoRows, got %v", err)
					}
				}
			}

			if !tt.expectError {
				if item == nil {
					t.Fatalf("Expected item, got nil")
				}
				if item.ID != tt.expectedItem.ID || item.Name != tt.expectedItem.Name {
					t.Errorf("GetItemByID() got = %v, want %v", item, tt.expectedItem)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}
