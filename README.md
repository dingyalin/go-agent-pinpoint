
# Pinpoint Go Agent

# Runnable Example

```
cd examples/server
go run main.go
```

# 配置文件
```
enabled: true
app_name: my_app
agent_id: my_agent
collector:
  ip: 127.0.0.1
  tcp_port: 9994
  stat_port: 9995
  span_port: 9996
log:	
  # sets Logger to log to either "stdout" or "stderr" (filenames are not supported)
  std: stdout
  # controls the pinpoint log level, must be "debug" for debug, or empty for info
  level: debug
```



# 环境变量

- PINPOINT_APP_NAME
- PINPOINT_AGENT_ID
- PINPOINT_COLLECTOR_IP
- PINPOINT_COLLECTOR_TCP_PORT
- PINPOINT_COLLECTOR_STAT_PORT
- PINPOINT_COLLECTOR_SPAN_PORT
- PINPOINT_LOG
sets Logger to log to either "stdout" or "stderr" (filenames are not supported)
- PINPOINT_LOG_LEVEL
controls the pinpoint log level, must be "debug" for debug, or empty for info
- PINPOINT_ENABLED


# 关闭监控

1. 设置环境变量

PINPOINT_ENABLED=false

2. 重启服务


# License

The Pinpoint Go agent is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
The Pinpoint Go agent Modified based on New Relic [newrelic/go-agent](https://github.com/newrelic/go-agent).
