package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiv1 "github.com.Saveliy12/microservices2.git/inventory/internal/api/inventory/v1"
	dbmongo "github.com.Saveliy12/microservices2.git/inventory/internal/client/db/mongo"
	partrepo "github.com.Saveliy12/microservices2.git/inventory/internal/repository/part"
	partservice "github.com.Saveliy12/microservices2.git/inventory/internal/service/part"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort        = 50051
	partsCollection = "parts"
)

// initMongo устанавливает соединение с MongoDB, проверяет доступность,
// создаёт уникальный индекс по uuid и наполняет коллекцию стартовыми данными
// Возвращает *mongo.Collection и функцию закрытия соединения
func initMongo(ctx context.Context) (*mongo.Collection, func(), error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, nil, errors.New("MONGO_URI is not set")
	}
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		return nil, nil, errors.New("MONGO_DB is not set")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	collection := client.Database(dbName).Collection(partsCollection)

	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "uuid", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("failed to create unique index on uuid: %w", err)
	}

	if err := partrepo.SeedParts(ctx, collection); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("failed to seed parts: %w", err)
	}

	return collection, func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from mongodb: %v", cerr)
		}
	}, nil
}

// buildServiceStack собирает слои чистой архитектуры: репозиторий -> сервис -> API
func buildServiceStack(collection *mongo.Collection) *apiv1.API {
	dbCollection := dbmongo.NewMongoCollection(collection)
	repo := partrepo.NewRepository(dbCollection)
	service := partservice.NewService(repo)
	return apiv1.NewAPI(service)
}

func main() {
	// Загружаем переменные окружения из .env перед чтением настроек
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	// Настройка graceful shutdown с обработкой сигналов
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		stop()
		log.Fatalf("failed to listen: %v", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	// Инициализация соединения с MongoDB: индекс, seed
	collection, closeMongo, err := initMongo(ctx)
	if err != nil {
		stop()
		log.Fatalf("failed to initialize mongodb: %v", err)
	}
	defer closeMongo()

	api := buildServiceStack(collection)

	// Создаем gRPC сервер
	s := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 Inventory gRPC server listening on %d\n", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve Inventory gRPC server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Println("🛑 Shutting down Inventory gRPC server...")
	s.GracefulStop()
	log.Println("✅ Inventory gRPC server stopped")
}
