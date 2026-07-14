package handler

import (
	"context"
	"fmt"

	"github.com/Quantum103/menu-service/internal/service"
	pb "github.com/Quantum103/menu-service/proto/v1"
)

// MenuHandler реализует интерфейс MenuService из proto
type MenuHandler struct {
	pb.UnimplementedMenuServiceServer
	menuService service.MenuService
}

func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: svc}
}

// GetActiveMenu реализует Server Streaming из proto:
// rpc GetActiveMenu(GetMenuRequest) returns (stream GetMenuResponse)
func (h *MenuHandler) GetActiveMenu(req *pb.GetMenuRequest, stream pb.MenuService_GetActiveMenuServer) error {
	// 1. Получаем данные из сервиса (пока игнорируем req.Category, можно добавить фильтрацию позже)
	items, err := h.menuService.GetActiveMenu(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get active menu: %w", err)
	}

	// 2. Стримим данные: отправляем каждое сообщение по отдельности
	for _, item := range items {
		// Конвертируем внутреннюю модель в proto-модель
		pbItem := &pb.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			IsAvailable: item.IsAvailable,
		}

		// Отправляем одно сообщение в поток
		err := stream.Send(&pb.GetMenuResponse{
			Items: []*pb.Item{pbItem},
		})
		if err != nil {
			return fmt.Errorf("error streaming menu: %w", err)
		}
	}

	return nil
}

// GetItem реализует Unary RPC из proto:
// rpc GetItem(GetItemRequest) returns (Item)
func (h *MenuHandler) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.Item, error) {
	// 1. Получаем товар из сервиса по ID
	item, err := h.menuService.GetItem(ctx, req.ItemId)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	if item == nil {
		return nil, fmt.Errorf("item with id %s not found", req.ItemId)
	}

	// 2. Возвращаем proto-модель
	return &pb.Item{
		Id:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		IsAvailable: item.IsAvailable,
	}, nil
}
