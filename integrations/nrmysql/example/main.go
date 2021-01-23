// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/dingyalin/pinpoint-go-agent/integrations/nrmysql"
	"github.com/dingyalin/pinpoint-go-agent/pinpoint"
)

func main() {
	// Set up a local mysql docker container with:
	// docker run -it -p 3306:3306 --net "bridge" -e MYSQL_ALLOW_EMPTY_PASSWORD=true mysql

	db, err := sql.Open("nrmysql", "root:123456@(127.0.0.1:3306)/information_schema")
	if nil != err {
		panic(err)
	}

	app, err := pinpoint.NewApplication(
		pinpoint.ConfigFromYaml("./pinpoint.yml"),
		pinpoint.ConfigFromEnvironment(),
	)
	if nil != err {
		panic(err)
	}
	app.WaitForConnection(5 * time.Second)
	txn := app.StartTransaction("mysqlQuery")

	ctx := pinpoint.NewContext(context.Background(), txn)
	row := db.QueryRowContext(ctx, "SELECT count(*) from tables")
	var count int
	err = row.Scan(&count)
	if nil != err {
		fmt.Printf("row scan err: %s \n", err)
	}

	txn.End()
	app.Shutdown(5 * time.Second)

	fmt.Println("number of tables in information_schema", count)
}
