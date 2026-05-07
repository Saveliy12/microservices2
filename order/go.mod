module github.com.Saveliy12/microservices2.git/order

go 1.25.8

replace github.com.Saveliy12/microservices2.git/shared => ../shared

require (
	github.com.Saveliy12/microservices2.git/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.53.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
