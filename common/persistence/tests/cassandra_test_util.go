// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tests

import (
	"fmt"
	"path/filepath"

	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	p "go.temporal.io/server/common/persistence"
	"go.temporal.io/server/common/persistence/cassandra"
	"go.temporal.io/server/common/persistence/nosql/nosqlplugin/cassandra/gocql"
	"go.temporal.io/server/common/resolver"
)

const (
	// TODO hard code this dir for now
	//  need to merge persistence test config / initialization in one place
	testCassandraExecutionSchema  = "../../../schema/cassandra/temporal/schema.cql"
	testCassandraVisibilitySchema = "../../../schema/cassandra/visibility/schema.cql"
)

func SetUpCassandraDatabase(cfg *config.Cassandra) {
	adminCfg := *cfg
	// NOTE need to connect with empty name to create new database
	adminCfg.Keyspace = "system"

	session, err := gocql.NewSession(adminCfg, resolver.NewNoopResolver(), log.NewNoopLogger())
	if err != nil {
		panic(fmt.Sprintf("unable to create Cassandra session: %v", err))
	}
	defer session.Close()

	if err := cassandra.CreateCassandraKeyspace(
		session,
		cfg.Keyspace,
		1,
		true,
		log.NewNoopLogger(),
	); err != nil {
		panic(fmt.Sprintf("unable to create Cassandra keyspace: %v", err))
	}
}

func SetUpCassandraSchema(cfg *config.Cassandra) {
	session, err := gocql.NewSession(*cfg, resolver.NewNoopResolver(), log.NewNoopLogger())
	if err != nil {
		panic(fmt.Sprintf("unable to create Cassandra session: %v", err))
	}
	defer session.Close()

	schemaPath, err := filepath.Abs(testCassandraExecutionSchema)
	if err != nil {
		panic(err)
	}

	statements, err := p.LoadAndSplitQuery([]string{schemaPath})
	if err != nil {
		panic(err)
	}

	for _, stmt := range statements {
		if err = session.Query(stmt).Exec(); err != nil {
			panic(err)
		}
	}

	schemaPath, err = filepath.Abs(testCassandraVisibilitySchema)
	if err != nil {
		panic(err)
	}

	statements, err = p.LoadAndSplitQuery([]string{schemaPath})
	if err != nil {
		panic(err)
	}

	for _, stmt := range statements {
		if err = session.Query(stmt).Exec(); err != nil {
			panic(err)
		}
	}
}

func TearDownCassandraKeyspace(cfg *config.Cassandra) {
	adminCfg := *cfg
	// NOTE need to connect with empty name to create new database
	adminCfg.Keyspace = "system"

	session, err := gocql.NewSession(adminCfg, resolver.NewNoopResolver(), log.NewNoopLogger())
	if err != nil {
		panic(fmt.Sprintf("unable to create Cassandra session: %v", err))
	}
	defer session.Close()

	if err := cassandra.DropCassandraKeyspace(
		session,
		cfg.Keyspace,
		log.NewNoopLogger(),
	); err != nil {
		panic(fmt.Sprintf("unable to drop Cassandra keyspace: %v", err))
	}
}