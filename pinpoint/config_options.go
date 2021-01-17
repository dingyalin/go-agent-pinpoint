// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pinpoint

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	yaml "gopkg.in/yaml.v2"
)

// ConfigOption configures the Config when provided to NewApplication.
type ConfigOption func(*Config)

// ConfigEnabled sets the whether or not the agent is enabled.
func ConfigEnabled(enabled bool) ConfigOption {
	return func(cfg *Config) { cfg.Enabled = enabled }
}

// ConfigAppName sets the application name.
func ConfigAppName(appName string) ConfigOption {
	return func(cfg *Config) { cfg.AppName = appName }
}

// ConfigAgentID sets the agent id.
func ConfigAgentID(agentID string) ConfigOption {
	return func(cfg *Config) { cfg.AgentID = agentID }
}

// ConfigCollectorIP sets the collector ip.
func ConfigCollectorIP(collectorIP string) ConfigOption {
	return func(cfg *Config) { cfg.Collector.IP = collectorIP }
}

// ConfigCollectorTCPPort sets the collector tcp port.
func ConfigCollectorTCPPort(tcpPort int) ConfigOption {
	return func(cfg *Config) { cfg.Collector.TCPPort = tcpPort }
}

// ConfigCollectorStatPort sets the collector stat port.
func ConfigCollectorStatPort(statPort int) ConfigOption {
	return func(cfg *Config) { cfg.Collector.StatPort = statPort }
}

// ConfigCollectorSpanPort sets the collector span port.
func ConfigCollectorSpanPort(spanPort int) ConfigOption {
	return func(cfg *Config) { cfg.Collector.SpanPort = spanPort }
}

// ConfigCollectorUploaded set is upload.
func ConfigCollectorUploaded(uploaded bool) ConfigOption {
	return func(cfg *Config) { cfg.Collector.Uploaded = uploaded }
}

// ConfigCollectorUploadedAgentStat set is upload.
func ConfigCollectorUploadedAgentStat(uploadedAgentStat bool) ConfigOption {
	return func(cfg *Config) { cfg.Collector.UploadedAgentStat = uploadedAgentStat }
}

// ConfigServiceType sets the service type.
func ConfigServiceType(serviceType int16) ConfigOption {
	return func(cfg *Config) { cfg.ServiceType = serviceType }
}

// ConfigLicense sets the license.
func ConfigLicense(license string) ConfigOption {
	return func(cfg *Config) { cfg.License = license }
}

// ConfigDistributedTracerEnabled populates the Config's
// DistributedTracer.Enabled setting.
func ConfigDistributedTracerEnabled(enabled bool) ConfigOption {
	return func(cfg *Config) { cfg.DistributedTracer.Enabled = enabled }
}

// ConfigLogger populates the Config's Logger.
func ConfigLogger(l Logger) ConfigOption {
	return func(cfg *Config) { cfg.Logger = l }
}

// ConfigInfoLogger populates the config with basic Logger at info level.
func ConfigInfoLogger(w io.Writer) ConfigOption {
	return ConfigLogger(NewLogger(w))
}

// ConfigDebugLogger populates the config with a Logger at debug level.
func ConfigDebugLogger(w io.Writer) ConfigOption {
	return ConfigLogger(NewDebugLogger(w))
}

// ConfigFromEnvironment populates the config based on environment variables:
//
//  PINPOINT_APP_NAME                                sets AppName
//  PINPOINT_ATTRIBUTES_EXCLUDE                      sets Attributes.Exclude using a comma-separated list, eg. "request.headers.host,request.method"
//  PINPOINT_ATTRIBUTES_INCLUDE                      sets Attributes.Include using a comma-separated list
//  PINPOINT_DISTRIBUTED_TRACING_ENABLED             sets DistributedTracer.Enabled using strconv.ParseBool
//  PINPOINT_ENABLED                                 sets Enabled using strconv.ParseBool
//  PINPOINT_HIGH_SECURITY                           sets HighSecurity using strconv.ParseBool
//  PINPOINT_HOST                                    sets Host
//  PINPOINT_INFINITE_TRACING_SPAN_EVENTS_QUEUE_SIZE sets InfiniteTracing.SpanEvents.QueueSize using strconv.Atoi
//  PINPOINT_INFINITE_TRACING_TRACE_OBSERVER_PORT    sets InfiniteTracing.TraceObserver.Port using strconv.Atoi
//  PINPOINT_INFINITE_TRACING_TRACE_OBSERVER_HOST    sets InfiniteTracing.TraceObserver.Host
//  PINPOINT_LABELS                                  sets Labels using a semi-colon delimited string of colon-separated pairs, eg. "Server:One;DataCenter:Primary"
//  PINPOINT_LICENSE_KEY                             sets License
//  PINPOINT_LOG                                     sets Logger to log to either "stdout" or "stderr" (filenames are not supported)
//  PINPOINT_LOG_LEVEL                               controls the PINPOINT_LOG level, must be "debug" for debug, or empty for info
//  PINPOINT_PROCESS_HOST_DISPLAY_NAME               sets HostDisplayName
//  PINPOINT_SECURITY_POLICIES_TOKEN                 sets SecurityPoliciesToken
//  PINPOINT_UTILIZATION_BILLING_HOSTNAME            sets Utilization.BillingHostname
//  PINPOINT_UTILIZATION_LOGICAL_PROCESSORS          sets Utilization.LogicalProcessors using strconv.Atoi
//  PINPOINT_UTILIZATION_TOTAL_RAM_MIB               sets Utilization.TotalRAMMIB using strconv.Atoi
//
// This function is strict and will assign Config.Error if any of the
// environment variables cannot be parsed.
func ConfigFromEnvironment() ConfigOption {
	return configFromEnvironment(os.Getenv)
}

func configFromEnvironment(getenv func(string) string) ConfigOption {
	return func(cfg *Config) {
		// Because fields could have been assigned in a previous
		// ConfigOption, we only want to assign fields using environment
		// variables that have been populated.  This is especially
		// relevant for the string case where no processing occurs.
		assignBool := func(field *bool, name string) {
			if env := getenv(name); env != "" {
				if b, err := strconv.ParseBool(env); nil != err {
					cfg.Error = fmt.Errorf("invalid %s value: %s", name, env)
				} else {
					*field = b
				}
			}
		}
		assignInt := func(field *int, name string) {
			if env := getenv(name); env != "" {
				if i, err := strconv.Atoi(env); nil != err {
					cfg.Error = fmt.Errorf("invalid %s value: %s", name, env)
				} else {
					*field = i
				}
			}
		}
		assignString := func(field *string, name string) {
			if env := getenv(name); env != "" {
				*field = env
			}
		}

		assignBool(&cfg.Enabled, "PINPOINT_ENABLED")

		assignString(&cfg.AppName, "PINPOINT_APP_NAME")
		assignString(&cfg.AgentID, "PINPOINT_AGENT_ID")
		assignString(&cfg.Collector.IP, "PINPOINT_COLLECTOR_IP")

		assignInt(&cfg.Collector.TCPPort, "PINPOINT_COLLECTOR_TCP_PORT")
		assignInt(&cfg.Collector.StatPort, "PINPOINT_COLLECTOR_STAT_PORT")
		assignInt(&cfg.Collector.SpanPort, "PINPOINT_COLLECTOR_SPAN_PORT")

		assignBool(&cfg.HighSecurity, "PINPOINT_HIGH_SECURITY")
		assignString(&cfg.Host, "PINPOINT_HOST")
		assignString(&cfg.HostDisplayName, "PINPOINT_PROCESS_HOST_DISPLAY_NAME")
		assignInt(&cfg.InfiniteTracing.SpanEvents.QueueSize, "PINPOINT_INFINITE_TRACING_SPAN_EVENTS_QUEUE_SIZE")

		//assignString(&cfg.License, "PINPOINT_LICENSE_KEY")
		//assignBool(&cfg.DistributedTracer.Enabled, "PINPOINT_DISTRIBUTED_TRACING_ENABLED")
		//assignString(&cfg.SecurityPoliciesToken, "PINPOINT_SECURITY_POLICIES_TOKEN")

		//assignString(&cfg.Utilization.BillingHostname, "PINPOINT_UTILIZATION_BILLING_HOSTNAME")
		//assignString(&cfg.InfiniteTracing.TraceObserver.Host, "PINPOINT_INFINITE_TRACING_TRACE_OBSERVER_HOST")
		//assignInt(&cfg.InfiniteTracing.TraceObserver.Port, "PINPOINT_INFINITE_TRACING_TRACE_OBSERVER_PORT")
		//assignInt(&cfg.Utilization.LogicalProcessors, "PINPOINT_UTILIZATION_LOGICAL_PROCESSORS")
		//assignInt(&cfg.Utilization.TotalRAMMIB, "PINPOINT_UTILIZATION_TOTAL_RAM_MIB")

		if env := getenv("PINPOINT_LABELS"); env != "" {
			if labels := getLabels(getenv("PINPOINT_LABELS")); len(labels) > 0 {
				cfg.Labels = labels
			} else {
				cfg.Error = fmt.Errorf("invalid PINPOINT_LABELS value: %s", env)
			}
		}

		if env := getenv("PINPOINT_ATTRIBUTES_INCLUDE"); env != "" {
			cfg.Attributes.Include = strings.Split(env, ",")
		}
		if env := getenv("PINPOINT_ATTRIBUTES_EXCLUDE"); env != "" {
			cfg.Attributes.Exclude = strings.Split(env, ",")
		}

		if env := getenv("PINPOINT_LOG"); env != "" {
			if dest := getLogDest(env); dest != nil {
				if isDebugEnv(getenv("PINPOINT_LOG_LEVEL")) {
					cfg.Logger = NewDebugLogger(dest)
				} else {
					cfg.Logger = NewLogger(dest)
				}
			} else {
				cfg.Error = fmt.Errorf("invalid PINPOINT_LOG value %s", env)
			}
		}
	}
}

func getLogDest(env string) io.Writer {
	switch env {
	case "stdout", "Stdout", "STDOUT":
		return os.Stdout
	case "stderr", "Stderr", "STDERR":
		return os.Stderr
	default:
		return nil
	}
}

func isDebugEnv(env string) bool {
	switch env {
	case "debug", "Debug", "DEBUG", "d", "D":
		return true
	default:
		return false
	}
}

// getLabels reads Labels from the env string, expressed as a semi-colon
// delimited string of colon-separated pairs (for example, "Server:One;Data
// Center:Primary").  Label keys and values must be 255 characters or less in
// length.  No more than 64 Labels can be set.
func getLabels(env string) map[string]string {
	out := make(map[string]string)
	env = strings.Trim(env, ";\t\n\v\f\r ")
	for _, entry := range strings.Split(env, ";") {
		if entry == "" {
			return nil
		}
		split := strings.Split(entry, ":")
		if len(split) != 2 {
			return nil
		}
		left := strings.TrimSpace(split[0])
		right := strings.TrimSpace(split[1])
		if left == "" || right == "" {
			return nil
		}
		if utf8.RuneCountInString(left) > 255 {
			runes := []rune(left)
			left = string(runes[:255])
		}
		if utf8.RuneCountInString(right) > 255 {
			runes := []rune(right)
			right = string(runes[:255])
		}
		out[left] = right
		if len(out) >= 64 {
			return out
		}
	}
	return out
}

type yamlConfig struct {
	Enabled   *bool  `yaml:"enabled"`
	AppName   string `yaml:"app_name"`
	AgentID   string `yaml:"agent_id"`
	Collector struct {
		IP       string `yaml:"ip"`
		TCPPort  int    `yaml:"tcp_port"`
		StatPort int    `yaml:"stat_port"`
		SpanPort int    `yaml:"span_port"`
	}
	Log struct {
		STD   string `yaml:"std"`
		Level string `yaml:"level"`
	}
}

// ConfigFromYaml ...
func ConfigFromYaml(path string) ConfigOption {
	data, err := ioutil.ReadFile(path)
	return configFromYaml(data, err)
}

func configFromYaml(data []byte, err error) ConfigOption {
	return func(cfg *Config) {
		if err != nil {
			cfg.Error = err
			return
		}

		yc := yamlConfig{}
		err := yaml.Unmarshal([]byte(data), &yc)
		if err != nil {
			cfg.Error = err
			return
		}

		if yc.Enabled != nil {
			cfg.Enabled = *yc.Enabled
		}
		if yc.AppName != "" {
			cfg.AppName = yc.AppName
		}
		if yc.AgentID != "" {
			cfg.AgentID = yc.AgentID
		}
		if yc.Collector.IP != "" {
			cfg.Collector.IP = yc.Collector.IP
		}
		if yc.Collector.TCPPort != 0 {
			cfg.Collector.TCPPort = yc.Collector.TCPPort
		}
		if yc.Collector.StatPort != 0 {
			cfg.Collector.StatPort = yc.Collector.StatPort
		}
		if yc.Collector.SpanPort != 0 {
			cfg.Collector.SpanPort = yc.Collector.SpanPort
		}

		if std := yc.Log.STD; std != "" {
			if dest := getLogDest(std); dest != nil {
				if isDebugEnv(yc.Log.Level) {
					cfg.Logger = NewDebugLogger(dest)
				} else {
					cfg.Logger = NewLogger(dest)
				}
			} else {
				cfg.Error = fmt.Errorf("invalid log std value %s", std)
			}
		}
	}

}
