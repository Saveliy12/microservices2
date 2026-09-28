//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Order Repository Integration Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	if err := logger.Init(loggerLevelValue, true); err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = ginkgo.AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")

	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}

	suiteCancel()
})
