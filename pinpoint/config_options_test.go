// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pinpoint

import (
	"reflect"
	"testing"
)

func TestConfigFromYaml(t *testing.T) {
	var data = `
enabled: false
app_name: my_app
agent_id: my_agent
collector:
  ip: 10.10.10.10
  tcp_port: 9984
  stat_port: 9985
  span_port: 9986
`

	cfgOpt := configFromYaml([]byte(data), nil)
	cfg := defaultConfig()
	cfgOpt(&cfg)

	expect := defaultConfig()
	expect.Enabled = false
	expect.AppName = "my_app"
	expect.AgentID = "my_agent"
	expect.Collector.IP = "10.10.10.10"
	expect.Collector.TCPPort = 9984
	expect.Collector.StatPort = 9985
	expect.Collector.SpanPort = 9986

	if !reflect.DeepEqual(expect, cfg) {
		t.Errorf("cfg   : %#v", cfg)
		t.Errorf("expect: %#v", expect)
	}
}

func TestConfigFromEnvironment(t *testing.T) {
	cfgOpt := configFromEnvironment(func(s string) string {
		switch s {
		case "PINPOINT_APP_NAME":
			return "my app"
		case "PINPOINT_AGENT_ID":
			return "my agent id"
		case "PINPOINT_COLLECTOR_IP":
			return "10.10.10.10"
		case "PINPOINT_COLLECTOR_TCP_PORT":
			return "9984"
		case "PINPOINT_COLLECTOR_STAT_PORT":
			return "9985"
		case "PINPOINT_COLLECTOR_SPAN_PORT":
			return "9986"
		case "PINPOINT_DISTRIBUTED_TRACING_ENABLED":
			return "true"
		case "PINPOINT_ENABLED":
			return "false"
		case "PINPOINT_HIGH_SECURITY":
			return "1"
		case "PINPOINT_HOST":
			return "my host"
		case "PINPOINT_PROCESS_HOST_DISPLAY_NAME":
			return "my display host"
		case "PINPOINT_LABELS":
			return "star:car;far:bar"
		case "PINPOINT_ATTRIBUTES_INCLUDE":
			return "zip,zap"
		case "PINPOINT_ATTRIBUTES_EXCLUDE":
			return "zop,zup,zep"
		case "PINPOINT_INFINITE_TRACING_SPAN_EVENTS_QUEUE_SIZE":
			return "98765"
		}
		return ""
	})
	expect := defaultConfig()
	expect.AppName = "my app"
	expect.AgentID = "my agent id"
	expect.Collector.IP = "10.10.10.10"
	expect.Collector.TCPPort = 9984
	expect.Collector.StatPort = 9985
	expect.Collector.SpanPort = 9986
	expect.DistributedTracer.Enabled = true
	expect.Enabled = false
	expect.HighSecurity = true
	expect.Host = "my host"
	expect.HostDisplayName = "my display host"
	expect.Labels = map[string]string{"star": "car", "far": "bar"}
	expect.Attributes.Include = []string{"zip", "zap"}
	expect.Attributes.Exclude = []string{"zop", "zup", "zep"}
	expect.InfiniteTracing.SpanEvents.QueueSize = 98765

	cfg := defaultConfig()
	cfgOpt(&cfg)

	if !reflect.DeepEqual(expect, cfg) {
		t.Errorf("cfg   : %#v", cfg)
		t.Errorf("expect: %#v", expect)
	}
}

func TestConfigFromEnvironmentIgnoresUnset(t *testing.T) {
	// test that configFromEnvironment ignores unset env vars
	cfgOpt := configFromEnvironment(func(string) string { return "" })
	cfg := defaultConfig()
	cfg.AppName = "something"
	cfg.Labels = map[string]string{"hello": "world"}
	cfg.Attributes.Include = []string{"zip", "zap"}
	cfg.Attributes.Exclude = []string{"zop", "zup", "zep"}
	cfg.License = "something"
	cfg.DistributedTracer.Enabled = true
	cfg.HighSecurity = true
	cfg.Host = "something"
	cfg.HostDisplayName = "something"
	cfg.SecurityPoliciesToken = "something"
	cfg.Utilization.BillingHostname = "something"
	cfg.Utilization.LogicalProcessors = 42
	cfg.Utilization.TotalRAMMIB = 42

	cfgOpt(&cfg)

	if cfg.AppName != "something" {
		t.Error("config value changed:", cfg.AppName)
	}
	if len(cfg.Labels) != 1 {
		t.Error("config value changed:", cfg.Labels)
	}
	if cfg.License != "something" {
		t.Error("config value changed:", cfg.License)
	}
	if !cfg.DistributedTracer.Enabled {
		t.Error("config value changed:", cfg.DistributedTracer.Enabled)
	}
	if !cfg.HighSecurity {
		t.Error("config value changed:", cfg.HighSecurity)
	}
	if cfg.Host != "something" {
		t.Error("config value changed:", cfg.Host)
	}
	if cfg.HostDisplayName != "something" {
		t.Error("config value changed:", cfg.HostDisplayName)
	}
	if cfg.SecurityPoliciesToken != "something" {
		t.Error("config value changed:", cfg.SecurityPoliciesToken)
	}
	if cfg.Utilization.BillingHostname != "something" {
		t.Error("config value changed:", cfg.Utilization.BillingHostname)
	}
	if cfg.Utilization.LogicalProcessors != 42 {
		t.Error("config value changed:", cfg.Utilization.LogicalProcessors)
	}
	if cfg.Utilization.TotalRAMMIB != 42 {
		t.Error("config value changed:", cfg.Utilization.TotalRAMMIB)
	}
	if len(cfg.Attributes.Include) != 2 {
		t.Error("config value changed:", cfg.Attributes.Include)
	}
	if len(cfg.Attributes.Exclude) != 3 {
		t.Error("config value changed:", cfg.Attributes.Exclude)
	}
}

func TestConfigFromEnvironmentAttributes(t *testing.T) {
	cfgOpt := configFromEnvironment(func(s string) string {
		switch s {
		case "PINPOINT_ATTRIBUTES_INCLUDE":
			return "zip,zap"
		case "PINPOINT_ATTRIBUTES_EXCLUDE":
			return "zop,zup,zep"
		default:
			return ""
		}
	})
	cfg := defaultConfig()
	cfgOpt(&cfg)
	if !reflect.DeepEqual(cfg.Attributes.Include, []string{"zip", "zap"}) {
		t.Error("incorrect config value:", cfg.Attributes.Include)
	}
	if !reflect.DeepEqual(cfg.Attributes.Exclude, []string{"zop", "zup", "zep"}) {
		t.Error("incorrect config value:", cfg.Attributes.Exclude)
	}
}

func TestConfigFromEnvironmentInvalidBool(t *testing.T) {
	cfgOpt := configFromEnvironment(func(s string) string {
		switch s {
		case "PINPOINT_ENABLED":
			return "BOGUS"
		default:
			return ""
		}
	})
	cfg := defaultConfig()
	cfgOpt(&cfg)
	if cfg.Error == nil {
		t.Error("error expected")
	}
}

func TestConfigFromEnvironmentInvalidInt(t *testing.T) {
	cfgOpt := configFromEnvironment(func(s string) string {
		switch s {
		case "PINPOINT_UTILIZATION_LOGICAL_PROCESSORS":
			return "BOGUS"
		default:
			return ""
		}
	})
	cfg := defaultConfig()
	cfgOpt(&cfg)
	if cfg.Error == nil {
		t.Error("error expected")
	}
}

func TestConfigFromEnvironmentInvalidLogger(t *testing.T) {
	cfgOpt := configFromEnvironment(func(s string) string {
		switch s {
		case "PINPOINT_LOG":
			return "BOGUS"
		default:
			return ""
		}
	})
	cfg := defaultConfig()
	cfgOpt(&cfg)
	if cfg.Error == nil {
		t.Error("error expected")
	}
}

func TestConfigFromEnvironmentInvalidLabels(t *testing.T) {
	cfgOpt := configFromEnvironment(func(s string) string {
		switch s {
		case "PINPOINT_LABELS":
			return ";;;"
		default:
			return ""
		}
	})
	cfg := defaultConfig()
	cfgOpt(&cfg)
	if cfg.Error == nil {
		t.Error("error expected")
	}
}

func TestConfigFromEnvironmentLabelsSuccess(t *testing.T) {
	cfgOpt := configFromEnvironment(func(s string) string {
		switch s {
		case "PINPOINT_LABELS":
			return "zip:zap; zop:zup"
		default:
			return ""
		}
	})
	cfg := defaultConfig()
	cfgOpt(&cfg)
	if !reflect.DeepEqual(cfg.Labels, map[string]string{"zip": "zap", "zop": "zup"}) {
		t.Error(cfg.Labels)
	}
}
