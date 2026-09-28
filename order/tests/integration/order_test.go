//go:build integration

package integration

import (
	"context"

	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("OrderRepository", func() {
	var ctx context.Context

	ginkgo.BeforeEach(func() {
		ctx = suiteCtx
		gomega.Expect(env.ClearOrders(ctx)).To(gomega.Succeed())
	})

	ginkgo.It("создаёт и читает заказ (roundtrip)", func() {
		order := BuildTestOrder()

		gomega.Expect(env.Repo.CreateOrder(ctx, order)).To(gomega.Succeed())

		got, err := env.Repo.GetOrder(ctx, order.OrderUUID)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(got.OrderUUID).To(gomega.Equal(order.OrderUUID))
		gomega.Expect(got.UserUUID).To(gomega.Equal(order.UserUUID))
		gomega.Expect(got.PartUUIDs).To(gomega.Equal(order.PartUUIDs))
		gomega.Expect(got.TotalPrice).To(gomega.Equal(order.TotalPrice))
		gomega.Expect(got.Status).To(gomega.Equal(order.Status))
		gomega.Expect(got.TransactionUUID).To(gomega.BeNil())
		gomega.Expect(got.PaymentMethod).To(gomega.BeNil())
	})

	ginkgo.It("возвращает ошибку для несуществующего заказа", func() {
		_, err := env.Repo.GetOrder(ctx, uuid.NewString())

		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("не позволяет создать заказ с дублирующимся UUID", func() {
		order := BuildTestOrder()

		gomega.Expect(env.Repo.CreateOrder(ctx, order)).To(gomega.Succeed())
		gomega.Expect(env.Repo.CreateOrder(ctx, order)).ToNot(gomega.Succeed())
	})
})
