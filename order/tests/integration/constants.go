package integration

import "time"

const (
	postgresImage    = "postgres:17.0-alpine3.20"
	postgresUser     = "demo"
	postgresPassword = "demo" //nolint:gosec
	postgresDB       = "postgres"

	testsTimeout   = 5 * time.Minute
	startupTimeout = 2 * time.Minute

	loggerLevelValue = "debug"
)
