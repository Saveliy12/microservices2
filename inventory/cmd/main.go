package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	inventoryV1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const grpcPort = 50051

type InventoryServiceServer struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func defaultParts() map[string]*inventoryV1.Part {
	engine := &inventoryV1.Part{
		Uuid:        "11111111-1111-1111-1111-111111111111",
		Name:        "Impulse Engine X1",
		Description: "Starter engine for test spaceship",
		Price:       1200,
		Category:    inventoryV1.Category_CATEGORY_ENGINE,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Nova Engines",
			Country: "RU",
			Website: "https://nova.example",
		},
		Tags: []string{"engine", "starter"},
	}

	fuel := &inventoryV1.Part{
		Uuid:        "22222222-2222-2222-2222-222222222222",
		Name:        "Fuel Tank F1",
		Description: "Compact fuel tank",
		Price:       800,
		Category:    inventoryV1.Category_CATEGORY_FUEL,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Orbital Fuel",
			Country: "RU",
			Website: "https://orbital.example",
		},
		Tags: []string{"fuel", "tank"},
	}

	return map[string]*inventoryV1.Part{
		engine.GetUuid(): engine,
		fuel.GetUuid():   fuel,
	}
}

func (s *InventoryServiceServer) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetResponse{
		Part: part,
	}, nil
}

func (s *InventoryServiceServer) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	current := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, p := range s.parts {
		current = append(current, p)
	}

	filter := req.GetFilter()
	if filter == nil {
		return &inventoryV1.ListPartsResponse{Parts: current}, nil
	}

	// UUID
	if len(filter.GetUuids()) > 0 {
		next := make([]*inventoryV1.Part, 0, len(current))
		for _, p := range current {
			if matchStringWrapperList(filter.GetUuids(), p.GetUuid()) {
				next = append(next, p)
			}
		}
		current = next
	}

	// Name
	if len(filter.GetNames()) > 0 {
		next := make([]*inventoryV1.Part, 0, len(current))
		for _, p := range current {
			if matchStringWrapperList(filter.GetNames(), p.GetName()) {
				next = append(next, p)
			}
		}
		current = next
	}

	// Category
	if len(filter.GetCategories()) > 0 {
		next := make([]*inventoryV1.Part, 0, len(current))
		for _, p := range current {
			if matchCategoryList(filter.GetCategories(), p.GetCategory()) {
				next = append(next, p)
			}
		}
		current = next
	}

	// Manufactured
	if len(filter.GetManufacturerCountries()) > 0 {
		next := make([]*inventoryV1.Part, 0, len(current))
		for _, p := range current {
			manufacturer := p.GetManufacturer()
			if manufacturer != nil && matchStringWrapperList(filter.GetManufacturerCountries(), manufacturer.GetCountry()) {
				next = append(next, p)
			}
		}
		current = next
	}

	// Tags
	if len(filter.GetTags()) > 0 {
		next := make([]*inventoryV1.Part, 0, len(current))
		for _, p := range current {
			if matchAnyTag(filter.GetTags(), p.GetTags()) {
				next = append(next, p)
			}
		}
		current = next
	}

	return &inventoryV1.ListPartsResponse{
		Parts: current,
	}, nil
}

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

	// Регистрируем наш сервис
	service := &InventoryServiceServer{
		// make(map[string]*inventoryV1.Part)
		parts: defaultParts(),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 Inventory gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Inventory gRPC server...")
	s.GracefulStop()
	log.Println("✅ Inventory Server stopped")
}

func matchStringWrapperList(values []*wrapperspb.StringValue, target string) bool {
	for _, v := range values {
		if v.GetValue() == target {
			return true
		}
	}
	return false
}

func matchCategoryList(values []inventoryV1.Category, target inventoryV1.Category) bool {
	return slices.Contains(values, target)
}

// Пересечение массивов (тегов)
func matchAnyTag(filterTags []*wrapperspb.StringValue, partTags []string) bool {
	for _, t := range partTags {
		if matchStringWrapperList(filterTags, t) {
			return true
		}
	}
	return false
}
