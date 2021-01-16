// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// +build go1.10

// Package nrmysql instruments https://github.com/go-sql-driver/mysql.
//
// Use this package to instrument your MySQL calls without having to manually
// create DatastoreSegments.  This is done in a two step process:
//
// 1. Use this package's driver in place of the mysql driver.
//
// If your code is using sql.Open like this:
//
//	import (
//		_ "github.com/go-sql-driver/mysql"
//	)
//
//	func main() {
//		db, err := sql.Open("mysql", "user@unix(/path/to/socket)/dbname")
//	}
//
// Then change the side-effect import to this package, and open "nrmysql" instead:
//
//	import (
//		_ "github.com/dingyalin/pinpoint-go-agent/integrations/nrmysql"
//	)
//
//	func main() {
//		db, err := sql.Open("nrmysql", "user@unix(/path/to/socket)/dbname")
//	}
//
// 2. Provide a context containing a pinpoint.Transaction to all exec and query
// methods on sql.DB, sql.Conn, sql.Tx, and sql.Stmt.  This requires using the
// context methods ExecContext, QueryContext, and QueryRowContext in place of
// Exec, Query, and QueryRow respectively.  For example, instead of the
// following:
//
//	row := db.QueryRow("SELECT count(*) from tables")
//
// Do this:
//
//	ctx := pinpoint.NewContext(context.Background(), txn)
//	row := db.QueryRowContext(ctx, "SELECT count(*) from tables")
//
// A working example is shown here:
// https://github.com/pinpoint/go-agent/tree/master/v3/integrations/nrmysql/example/main.go
package nrmysql

import (
	"database/sql"
	"net"

	"github.com/dingyalin/pinpoint-go-agent/internal"
	"github.com/dingyalin/pinpoint-go-agent/pinpoint"
	"github.com/dingyalin/pinpoint-go-agent/pinpoint/sqlparse"
	"github.com/go-sql-driver/mysql"
)

var (
	baseBuilder = pinpoint.SQLDriverSegmentBuilder{
		BaseSegment: pinpoint.DatastoreSegment{
			Product: pinpoint.DatastoreMySQL,
		},
		ParseQuery: sqlparse.ParseQuery,
		ParseDSN:   parseDSN,
	}
)

func init() {
	sql.Register("nrmysql", pinpoint.InstrumentSQLDriver(mysql.MySQLDriver{}, baseBuilder))
	internal.TrackUsage("integration", "driver", "mysql")
}

func parseDSN(s *pinpoint.DatastoreSegment, dsn string) {
	cfg, err := mysql.ParseDSN(dsn)
	if nil != err {
		return
	}
	parseConfig(s, cfg)
}

func parseConfig(s *pinpoint.DatastoreSegment, cfg *mysql.Config) {
	s.DatabaseName = cfg.DBName

	var host, ppoid string
	switch cfg.Net {
	case "unix", "unixgram", "unixpacket":
		host = "localhost"
		ppoid = cfg.Addr
	case "cloudsql":
		host = cfg.Addr
	default:
		var err error
		host, ppoid, err = net.SplitHostPort(cfg.Addr)
		if nil != err {
			host = cfg.Addr
		} else if host == "" {
			host = "localhost"
		}
	}

	s.Host = host
	s.PortPathOrID = ppoid
}
