module github.com/dingyalin/pinpoint-go-agent/integrations/nrredis-v7

// As of Jan 2020, go 1.11 is in the go-redis go.mod file:
// https://github.com/go-redis/redis/blob/master/go.mod
go 1.11

require (
	github.com/dingyalin/pinpoint-go-agent v1.0.0
	github.com/go-redis/redis/v7 v7.2.0
)

replace github.com/dingyalin/pinpoint-go-agent v1.0.0 => ../../
