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

func (h *MenuHandler) GetActiveMenu(req *pb.GetMenuRequest, stream pb.MenuService_GetActiveMenuServer) error {
	items, err := h.menuService.GetActiveMenu(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get active menu: %w", err)
	}

	for _, item := range items {
		pbItem := &pb.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			IsAvailable: item.IsAvailable,
		}

		err := stream.Send(&pb.GetMenuResponse{
			Items: []*pb.Item{pbItem},
		})
		if err != nil {
			return fmt.Errorf("error streaming menu: %w", err)
		}
	}

	return nil
}

func (h *MenuHandler) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.Item, error) {
	item, err := h.menuService.GetItem(ctx, req.ItemId)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	if item == nil {
		return nil, fmt.Errorf("item with id %s not found", req.ItemId)
	}

	return &pb.Item{
		Id:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		IsAvailable: item.IsAvailable,
	}, nil
}
