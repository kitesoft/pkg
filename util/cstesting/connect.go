// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cstesting

import (
	"fmt"
	"io"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// EnvDSN is the name of the environment variable
const EnvDSN string = "CS_DSN"

func getDSN(env string) (string, error) {
	dsn := os.Getenv(env)
	if dsn == "" {
		return "", errors.NewNotFoundf("DSN in environment variable %q not found.", EnvDSN)
	}
	return dsn, nil
}

// MustGetDSN returns the data source name from an environment variable or
// panics on error.
func MustGetDSN() string {
	d, err := getDSN(EnvDSN)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return d
}

type skipper interface {
	Skipf(format string, args ...interface{})
}

// MustConnectDB is a helper function that creates a new database connection
// using a DSN from an environment variable found in the constant csdb.EnvDSN.
// If the DSN environment variable has not been set it skips the test.
// Argument t specified usually the *testing.T/B struct.
func MustConnectDB(t skipper, opts ...dbr.ConnectionOption) *dbr.Connection {
	if _, err := getDSN(EnvDSN); errors.IsNotFound(err) {
		t.Skipf("%s", err)
	}
	cos := append([]dbr.ConnectionOption{}, dbr.WithDSN(MustGetDSN()))
	dbc := dbr.MustConnectAndVerify(append(cos, opts...)...)
	return dbc
}

// Close for usage in conjunction with defer.
// 		defer cstesting.Close(t,db)
func Close(t errorFormatter, c io.Closer) {
	if err := c.Close(); err != nil {
		t.Errorf("%+v", err)
	}
}

// MockDB creates a mocked database connection. Fatals on error.
func MockDB(t fataler) (*dbr.Connection, sqlmock.Sqlmock) {
	db, sm, err := sqlmock.New()
	fatalIfError(t, err)

	dbc, err := dbr.NewConnection(dbr.WithDB(db))
	fatalIfError(t, err)
	return dbc, sm
}

// MockClose for usage in conjunction with defer.
// 		defer cstesting.MockClose(t,db,dbMock)
func MockClose(t fataler, c io.Closer, m sqlmock.Sqlmock) {
	m.ExpectClose()
	if err := c.Close(); err != nil {
		t.Fatalf("%+v", err)
	}
	if err := m.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations:\n%+v", err)
	}
}
