module github.com/dingyalin/pinpoint-go-agent/integrations/nrmicro

// As of Dec 2019, the go-micro go.mod file uses 1.13:
// https://github.com/micro/go-micro/blob/master/go.mod
go 1.13

require (
	github.com/dingyalin/pinpoint-go-agent v1.0.0
	github.com/golang/protobuf v1.4.0
	github.com/micro/go-micro v1.8.0
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
)

replace github.com/dingyalin/pinpoint-go-agent v1.0.0 => ../../
