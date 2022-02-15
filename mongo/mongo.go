// Copyright 2021 helloshaohua <wu.shaohua@foxmail.com>;;
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/coolstina/httpclient/rawquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewConnection initialize mongodb client for connection instance.
func NewConnection(host, username, password string, ops ...Option) (*mongo.Client, error) {
	opts := option(host, username, password)
	for _, o := range ops {
		o.apply(opts)
	}

	uri, err := uri(opts)
	if err != nil {
		return nil, err
	}

	fmt.Printf("uri: %+v\n", uri)
	// uri = "mongodb://root:root@localhost:27017/?maxPoolSize=20&w=majority"
	// uri = "mongodb://root:root@localhost:27017/?maxPoolSize=100&replicaSet=null&maxIdleTimeMS=0&minPoolSize=0&tls=false&w=null&directConnection=false"
	// replicaSet https://www.mongodb.com/community/forums/t/server-selection-error-server-selection-timeout-current-topology/2930
	uri = "mongodb://root:root@localhost:27017/?maxPoolSize=100&replicaSet=hello_world&maxIdleTimeMS=0&minPoolSize=0&tls=false&w=null&directConnection=false"

	connect, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return connect, nil
}

func uri(opts *opts) (string, error) {
	uri := fmt.Sprintf(`mongodb://%s:%s@%s/`,
		opts.username,
		opts.password,
		strings.Join(opts.hosts, ","),
	)

	parse, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	ot := reflect.TypeOf(opts).Elem()
	ov := reflect.ValueOf(opts).Elem()
	queries := make([]rawquery.Query, 0, ot.NumField())

	for i := 0; i < ot.NumField(); i++ {
		name := ot.Field(i).Name
		switch name {
		case "hosts":
			fallthrough
		case "username":
			fallthrough
		case "password":
			continue
		default:
			query := rawquery.Query{Field: name}

			switch ov.Field(i).Kind() {
			case reflect.String:
				query.Value = ov.Field(i).String()
			case reflect.Bool:
				query.Value = ov.Field(i).Bool()
			case reflect.Int64:
				query.Value = ov.Field(i).Int()
			case reflect.Int:
				query.Value = ov.Field(i).Int()
			}

			queries = append(queries, query)
		}
	}

	parse.RawQuery = rawquery.NewRawQueryWithQueries(queries)
	return parse.String(), err
}

func option(host string, username string, password string) *opts {
	var opts = &opts{
		hosts:                    []string{host},
		username:                 username,
		password:                 password,
		connectTimeoutMS:         time.Second * 30,
		maxPoolSize:              100,
		replicaSet:               "null",
		maxIdleTimeMS:            0,
		minPoolSize:              0,
		socketTimeoutMS:          time.Second * 30,
		serverSelectionTimeoutMS: time.Second * 10,
		tls:                      false,
		w:                        "null",
		directConnection:         false,
	}
	return opts
}
