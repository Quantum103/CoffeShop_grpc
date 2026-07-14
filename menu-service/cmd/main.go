package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Quantum103/menu-service/internal/database"
	"github.com/Quantum103/menu-service/internal/handler"
	"github.com/Quantum103/menu-service/internal/repository"
	"github.com/Quantum103/menu-service/internal/service"
	pb "github.com/Quantum103/menu-service/proto/v1"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Menu service is starting")

	dbConfig := database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "coffee_user",
		Password: "coffee_password",
		DBName:   "menu_db",
	}

	db, err := database.NewPostgresPool(dbConfig)
	if err != nil {
		log.Fatalf(" Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	repo := repository.NewMenuRepository(db)
	svc := service.NewMenuService(repo)
	h := handler.NewMenuHandler(svc)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterMenuServiceServer(grpcServer, h)

	fmt.Println("Menu Service запущен")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
