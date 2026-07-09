module github.com/emreEngineering/kervan/services/auth-service

go 1.26.2

require (
	github.com/emreEngineering/kervan/gen/go v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/lib/pq v1.12.3
	golang.org/x/crypto v0.53.0
	google.golang.org/grpc v1.82.0
)

require (
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/emreEngineering/kervan/gen/go => ../../gen/go
