package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Quantum103/CoffeShop_grpc/proto/coffeeshop_proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeShopServer
}

func (s *server) GetMenu(MenuRequest *pb.MenuRequest, srv grpc.ServerStreamingServer[pb.Menu]) error {
	items := []*pb.Item{
		&pb.Item{
			Id:   "1",
			Name: "Americano",
		},
		&pb.Item{
			Id:   "2",
			Name: "Latte",
		},
		&pb.Item{
			Id:   "3",
			Name: "Capuchino",
		},
	}
	for i, _ := range items {
		srv.Send(&pb.Menu{
			Items: items[0 : i+1],
		})
	}
	return nil
}
func (s *server) PlaceOrder(context context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC123",
	}, nil
}
func (s *server) GetOrderStatus(context context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN PROGRESS",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("error: failed to listened port 8080")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterCoffeShopServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %s", err)
	}
}
