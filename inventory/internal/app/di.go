package app

import (
	"context"

	apiv1 "github.com.Saveliy12/microservices2.git/inventory/internal/api/inventory/v1"
	db "github.com.Saveliy12/microservices2.git/inventory/internal/client/db"
	dbmongo "github.com.Saveliy12/microservices2.git/inventory/internal/client/db/mongo"
	"github.com.Saveliy12/microservices2.git/inventory/internal/config"
	"github.com.Saveliy12/microservices2.git/inventory/internal/repository"
	partrepo "github.com.Saveliy12/microservices2.git/inventory/internal/repository/part"
	inventoryservice "github.com.Saveliy12/microservices2.git/inventory/internal/service"
	inventoryimpl "github.com.Saveliy12/microservices2.git/inventory/internal/service/part"
	"github.com.Saveliy12/microservices2.git/platform/pkg/closer"
	"github.com.Saveliy12/microservices2.git/platform/pkg/grpc/health"
	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serviceProvider struct {
	config         config.InventoryConfig
	mongoClient    *mongo.Client
	collection     *mongo.Collection
	dbCollection   db.Collection
	service        inventoryservice.InventoryService
	partRepository repository.PartRepository
	api            *apiv1.API
	grpcServer     *grpc.Server
	closer         *closer.Closer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) Config() (config.InventoryConfig, error) {
	if sp.config == nil {
		cfg, err := config.New()
		if err != nil {
			return nil, err
		}

		sp.config = cfg
	}

	return sp.config, nil
}

func (sp *serviceProvider) MongoClient() (*mongo.Client, error) {
	if sp.mongoClient == nil {
		cfg, err := sp.Config()
		if err != nil {
			return nil, err
		}

		client, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo().URI()))
		if err != nil {
			return nil, err
		}

		if err := client.Ping(context.Background(), nil); err != nil {
			_ = client.Disconnect(context.Background())
			return nil, err
		}

		sp.mongoClient = client
		sp.Closer().AddNamed("mongo", func(ctx context.Context) error { return client.Disconnect(ctx) })
	}

	return sp.mongoClient, nil
}

func (sp *serviceProvider) Collection() (*mongo.Collection, error) {
	if sp.collection == nil {
		client, err := sp.MongoClient()
		if err != nil {
			return nil, err
		}

		cfg, err := sp.Config()
		if err != nil {
			return nil, err
		}

		coll := client.Database(cfg.Mongo().Database()).Collection("parts")

		_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		})

		_ = partrepo.SeedParts(context.Background(), coll)

		sp.collection = coll
	}

	return sp.collection, nil
}

func (sp *serviceProvider) DBCollection() db.Collection {
	if sp.dbCollection == nil {
		coll, _ := sp.Collection()
		sp.dbCollection = dbmongo.NewMongoCollection(coll)
	}

	return sp.dbCollection
}

func (sp *serviceProvider) PartRepository() repository.PartRepository {
	if sp.partRepository == nil {
		sp.partRepository = partrepo.NewRepository(sp.DBCollection())
	}

	return sp.partRepository
}

func (sp *serviceProvider) InventoryService() inventoryservice.InventoryService {
	if sp.service == nil {
		sp.service = inventoryimpl.NewService(sp.PartRepository())
	}

	return sp.service
}

func (sp *serviceProvider) InventoryAPI() *apiv1.API {
	if sp.api == nil {
		sp.api = apiv1.NewAPI(sp.InventoryService())
	}

	return sp.api
}

func (sp *serviceProvider) GRPCServer() *grpc.Server {
	if sp.grpcServer == nil {
		server := grpc.NewServer()
		inventoryv1.RegisterInventoryServiceServer(server, sp.InventoryAPI())
		health.RegisterService(server)
		reflection.Register(server)

		sp.grpcServer = server
	}

	return sp.grpcServer
}

func (sp *serviceProvider) Closer() *closer.Closer {
	if sp.closer == nil {
		sp.closer = closer.NewWithLogger(logger.Logger())
	}

	return sp.closer
}
