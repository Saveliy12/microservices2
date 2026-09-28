//go:build integration

package integration

import (
	"context"

	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Тесты общаются с сервисом только через gRPC
var _ = ginkgo.Describe("InventoryService e2e", func() {
	var (
		ctx    context.Context
		client inventoryv1.InventoryServiceClient
	)

	ginkgo.BeforeEach(func() {
		ctx = suiteCtx

		conn, err := grpc.NewClient(env.App.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		client = inventoryv1.NewInventoryServiceClient(conn)
	})

	ginkgo.It("возвращает список деталей с сидовыми данными", func() {
		resp, err := client.ListParts(ctx, &inventoryv1.ListPartsRequest{Filter: &inventoryv1.PartsFilter{}})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(resp.GetParts()).ToNot(gomega.BeEmpty())

		names := make([]string, 0, len(resp.GetParts()))
		for _, part := range resp.GetParts() {
			names = append(names, part.GetName())
		}
		gomega.Expect(names).To(gomega.ContainElement("Engine V8"))
	})

	ginkgo.It("возвращает деталь по UUID из списка", func() {
		listResp, err := client.ListParts(ctx, &inventoryv1.ListPartsRequest{Filter: &inventoryv1.PartsFilter{}})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(listResp.GetParts()).ToNot(gomega.BeEmpty())

		expected := listResp.GetParts()[0]

		getResp, err := client.GetPart(ctx, &inventoryv1.GetPartRequest{Uuid: expected.GetUuid()})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(getResp.GetPart().GetUuid()).To(gomega.Equal(expected.GetUuid()))
		gomega.Expect(getResp.GetPart().GetName()).To(gomega.Equal(expected.GetName()))
		gomega.Expect(getResp.GetPart().GetPrice()).To(gomega.Equal(expected.GetPrice()))
	})

	ginkgo.It("возвращает NotFound для неизвестного UUID", func() {
		_, err := client.GetPart(ctx, &inventoryv1.GetPartRequest{Uuid: uuid.NewString()})

		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(status.Code(err)).To(gomega.Equal(codes.NotFound))
	})

	ginkgo.It("читает деталь, вставленную напрямую в MongoDB", func() {
		partUUID, err := env.InsertTestPart(ctx)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		resp, err := client.GetPart(ctx, &inventoryv1.GetPartRequest{Uuid: partUUID})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(resp.GetPart().GetName()).To(gomega.Equal("Test Part"))
		gomega.Expect(resp.GetPart().GetPrice()).To(gomega.Equal(123.45))
	})
})
