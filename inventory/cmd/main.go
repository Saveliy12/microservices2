package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiv1 "github.com.Saveliy12/microservices2.git/inventory/internal/api/inventory/v1"
	partrepo "github.com.Saveliy12/microservices2.git/inventory/internal/repository/part"
	partservice "github.com.Saveliy12/microservices2.git/inventory/internal/service/part"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	partRepository := partrepo.NewRepository()
	partService := partservice.NewService(partRepository)
	inventoryAPI := apiv1.NewAPI(partService)

	inventoryv1.RegisterInventoryServiceServer(s, inventoryAPI)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 Inventory gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Inventory gRPC server: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Inventory gRPC server...")
	s.GracefulStop()
	log.Println("✅ Inventory gRPC server stopped")
}
