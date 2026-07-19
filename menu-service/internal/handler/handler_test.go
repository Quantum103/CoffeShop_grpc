package handler

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Quantum103/menu-service/internal/model"
	pb "github.com/Quantum103/menu-service/proto/v1"
	"google.golang.org/grpc"
)

type mockMenuService struct {
	getActiveMenuFn func(ctx context.Context) ([]model.Item, error)
	getItemFn       func(ctx context.Context, id string) (*model.Item, error)
}

func (m *mockMenuService) GetActiveMenu(ctx context.Context) ([]model.Item, error) {
	return m.getActiveMenuFn(ctx)
}

func (m *mockMenuService) GetItem(ctx context.Context, id string) (*model.Item, error) {
	return m.getItemFn(ctx, id)
}

type mockStream struct {
	grpc.ServerStream
	sendFn func(*pb.GetMenuResponse) error
	ctx    context.Context
}

func (m *mockStream) Send(resp *pb.GetMenuResponse) error {
	return m.sendFn(resp)
}

func (m *mockStream) Context() context.Context {
	return m.ctx
}

func TestMenuHandler_GetActiveMenu(t *testing.T) {
	tests := []struct {
		name          string
		serviceItems  []model.Item
		serviceErr    error
		streamSendErr error
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			serviceItems: []model.Item{
				{ID: "1", Name: "Burger", Description: "Tasty", Price: 10.0, IsAvailable: true},
				{ID: "2", Name: "Pizza", Description: "Cheesy", Price: 15.0, IsAvailable: false},
			},
			expectError: false,
		},
		{
			name:          "Service Error",
			serviceErr:    errors.New("database connection failed"),
			expectError:   true,
			errorContains: "failed to get active menu",
		},
		{
			name: "Stream Send Error",
			serviceItems: []model.Item{
				{ID: "1", Name: "Burger", Description: "Tasty", Price: 10.0, IsAvailable: true},
			},
			streamSendErr: errors.New("client disconnected"),
			expectError:   true,
			errorContains: "error streaming menu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockMenuService{
				getActiveMenuFn: func(ctx context.Context) ([]model.Item, error) {
					return tt.serviceItems, tt.serviceErr
				},
			}

			handler := NewMenuHandler(mockSvc)

			stream := &mockStream{
				sendFn: func(resp *pb.GetMenuResponse) error {
					return tt.streamSendErr
				},
				ctx: context.Background(),
			}

			err := handler.GetActiveMenu(&pb.GetMenuRequest{}, stream)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
				return
			}

			if tt.expectError && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
			}
		})
	}
}

func TestMenuHandler_GetItem(t *testing.T) {
	tests := []struct {
		name          string
		req           *pb.GetItemRequest
		serviceItem   *model.Item
		serviceErr    error
		expectError   bool
		errorContains string
	}{
		{
			name:        "Success",
			req:         &pb.GetItemRequest{ItemId: "1"},
			serviceItem: &model.Item{ID: "1", Name: "Burger", Description: "Tasty", Price: 10.0, IsAvailable: true},
			expectError: false,
		},
		{
			name:          "Service Error",
			req:           &pb.GetItemRequest{ItemId: "1"},
			serviceErr:    errors.New("database query failed"),
			expectError:   true,
			errorContains: "failed to get item",
		},
		{
			name:          "Item Not Found",
			req:           &pb.GetItemRequest{ItemId: "999"},
			serviceItem:   nil,
			expectError:   true,
			errorContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockMenuService{
				getItemFn: func(ctx context.Context, id string) (*model.Item, error) {
					return tt.serviceItem, tt.serviceErr
				},
			}

			handler := NewMenuHandler(mockSvc)

			item, err := handler.GetItem(context.Background(), tt.req)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
				return
			}

			if tt.expectError && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
			}

			if !tt.expectError {
				if item == nil {
					t.Errorf("expected item, got nil")
					return
				}
				if item.Id != tt.serviceItem.ID || item.Name != tt.serviceItem.Name {
					t.Errorf("expected item ID %s and Name %s, got ID %s and Name %s",
						tt.serviceItem.ID, tt.serviceItem.Name, item.Id, item.Name)
				}
			}
		})
	}
}
