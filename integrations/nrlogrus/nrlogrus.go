// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package nrlogrus sends go-agent log messages to
// https://github.com/sirupsen/logrus.
//
// Use this package if you are using logrus in your application and would like
// the go-agent log messages to end up in the same place.  If you are using
// the logrus standard logger, use ConfigStandardLogger when creating your
// application:
//
//	app, err := pinpoint.NewApplication(
//		pinpoint.ConfigFromEnvironment(),
//		nrlogrus.ConfigStandardLogger(),
//	)
//
// If you are using a particular logrus Logger instance, then use ConfigLogger:
//
//	l := logrus.New()
//	l.SetLevel(logrus.DebugLevel)
//	app, err := pinpoint.NewApplication(
//		pinpoint.ConfigFromEnvironment(),
//		nrlogrus.ConfigLogger(l),
//	)
//
// This package requires logrus version v1.1.0 and above.
package nrlogrus

import (
	"github.com/dingyalin/pinpoint-go-agent/internal"
	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
	"github.com/sirupsen/logrus"
)

func init() { internal.TrackUsage("integration", "logging", "logrus") }

type shim struct {
	e *logrus.Entry
	l *logrus.Logger
}

func (s *shim) Error(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Error(msg)
}
func (s *shim) Warn(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Warn(msg)
}
func (s *shim) Info(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Info(msg)
}
func (s *shim) Debug(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Debug(msg)
}
func (s *shim) DebugEnabled() bool {
	lvl := s.l.GetLevel()
	return lvl >= logrus.DebugLevel
}

// StandardLogger returns a pinpoint.Logger which forwards agent log messages to
// the logrus package-level exported logger.
func StandardLogger() pinpoint.Logger {
	return Transform(logrus.StandardLogger())
}

// Transform turns a *logrus.Logger into a pinpoint.Logger.
func Transform(l *logrus.Logger) pinpoint.Logger {
	return &shim{
		l: l,
		e: l.WithFields(logrus.Fields{
			"component": "pinpoint",
		}),
	}
}

// ConfigLogger configures the pinpoint.Application to send log messsages to the
// provided logrus logger.
func ConfigLogger(l *logrus.Logger) pinpoint.ConfigOption {
	return pinpoint.ConfigLogger(Transform(l))
}

// ConfigStandardLogger configures the pinpoint.Application to send log
// messsages to the standard logrus logger.
func ConfigStandardLogger() pinpoint.ConfigOption {
	return pinpoint.ConfigLogger(StandardLogger())
}
