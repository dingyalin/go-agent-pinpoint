package io

const (
	ServiceTypeHTTPClient        = 9052
	ServiceTypeUnkonwnDB         = 2050
	ServiceTypeMysql             = 2100
	ServiceTypeMysqlExecuteQuery = 2101
	ServiceTypeOracle            = 2300

	ServiceTypePython             = 1550
	ServiceTypePythonMethod       = 1551
	ServiceTypePythonRemoteMethod = 9800

	ServiceTypeGo             = 1000
	ServiceTypeGoMethod       = 0
	ServiceTypeGoAsyncMethod  = 0
	ServiceTypeGoRemoteMethod = 0

	ServiceTypeMemcached = 8050
	ServiceTypeRedis     = 8200
)
