module github.com.Saveliy12/microservices2.git/payment

go 1.25.8

replace github.com.Saveliy12/microservices2.git/shared => ../shared

replace github.com.Saveliy12/microservices2.git/platform => ../platform

require (
	github.com.Saveliy12/microservices2.git/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.83.1
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
)

require (
	github.com.Saveliy12/microservices2.git/platform v0.0.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/stretchr/testify v1.11.1
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260720211330-0afa2a65878a // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
