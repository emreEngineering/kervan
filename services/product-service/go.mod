module github.com/emreEngineering/kervan/services/product-service

go 1.26.2

require (
	github.com/emreEngineering/kervan/gen/go v0.0.0-20260709194251-0333e7794538
	github.com/lib/pq v1.12.3
	google.golang.org/grpc v1.82.1
)

require (
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/emreEngineering/kervan/gen/go => ../../gen/go
