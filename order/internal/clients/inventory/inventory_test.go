package inventory_test

import (
	"context"
	"net"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/clients/inventory"
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type InventoryClientSuite struct {
	suite.Suite
	ctx    context.Context
	conn   *grpc.ClientConn
	client *inventory.Client
	server *inventoryServiceServer
}

func (s *InventoryClientSuite) SetupTest() {
	s.ctx = context.Background()

	s.server = newInventoryServiceServer()
	lis := bufconn.Listen(1 << 20)

	grpcServer := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(grpcServer, s.server)
	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	s.conn = conn
	s.client = inventory.New(conn)
}

func (s *InventoryClientSuite) TearDownTest() {
	_ = s.conn.Close()
}

func (s *InventoryClientSuite) TestListParts_Success() {
	s.server.listPartsResponse = &inventoryv1.ListPartsResponse{
		Parts: []*inventoryv1.Part{
			{Uuid: "part-1", Price: 100.5},
			{Uuid: "part-2", Price: 200.75},
		},
	}

	parts, err := s.client.ListParts(s.ctx, []string{"part-1", "part-2"})

	s.Require().NoError(err)
	s.Equal([]model.PartQuote{
		{UUID: "part-1", Price: 100.5},
		{UUID: "part-2", Price: 200.75},
	}, parts)
	s.Equal([]string{"part-1", "part-2"}, s.server.lastFilterUUIDs)
}

func (s *InventoryClientSuite) TestListParts_ServerError() {
	s.server.listPartsErr = status.Error(codes.Internal, "inventory unavailable")

	parts, err := s.client.ListParts(s.ctx, []string{"part-1"})

	s.Require().Error(err)
	s.Nil(parts)
}

func (s *InventoryClientSuite) TestListParts_EmptyResponse() {
	s.server.listPartsResponse = &inventoryv1.ListPartsResponse{}

	parts, err := s.client.ListParts(s.ctx, []string{"part-1"})

	s.Require().NoError(err)
	s.Empty(parts)
}

func TestInventoryClientSuite(t *testing.T) {
	suite.Run(t, new(InventoryClientSuite))
}

// inventoryServiceServer — тестовая реализация inventoryv1.InventoryServiceServer
type inventoryServiceServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	listPartsResponse *inventoryv1.ListPartsResponse
	listPartsErr      error
	lastFilterUUIDs   []string
}

func newInventoryServiceServer() *inventoryServiceServer {
	return &inventoryServiceServer{}
}

func (s *inventoryServiceServer) ListParts(_ context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	if req.GetFilter() != nil {
		for _, uuid := range req.GetFilter().GetUuids() {
			s.lastFilterUUIDs = append(s.lastFilterUUIDs, uuid.GetValue())
		}
	}

	if s.listPartsErr != nil {
		return nil, s.listPartsErr
	}

	return s.listPartsResponse, nil
}
