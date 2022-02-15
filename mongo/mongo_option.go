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
	"sync"
	"time"
)

type Option interface {
	apply(*opts)
}

type optionFunc func(ops *opts)

func (o optionFunc) apply(ops *opts) {
	o(ops)
}

type opts struct {
	hosts                    []string
	username                 string
	password                 string
	connectTimeoutMS         time.Duration
	maxPoolSize              int
	replicaSet               string
	maxIdleTimeMS            time.Duration
	minPoolSize              int
	socketTimeoutMS          time.Duration
	serverSelectionTimeoutMS time.Duration
	tls                      bool
	w                        string
	directConnection         bool
}

var access sync.Mutex

// WithHosts Specifies mongod server hosts, you can specify one or more then.
func WithHosts(hosts ...string) Option {
	access.Lock()
	defer access.Unlock()
	return optionFunc(func(ops *opts) {
		if ops.hosts == nil {
			ops.hosts = make([]string, 0)
		}
		for _, host := range hosts {
			if !in(ops.hosts, host) {
				ops.hosts = append(ops.hosts, host)
			}
		}
	})
}

// WithUsername Specifies user's username.
func WithUsername(username string) Option {
	return optionFunc(func(ops *opts) {
		ops.username = username
	})
}

// WithPassword Specifies user's password.
func WithPassword(password string) Option {
	return optionFunc(func(ops *opts) {
		ops.password = password
	})
}

// WithConnectTimeoutMS Specifies the number of
// milliseconds to wait before timeout on a TCP connection.
// Default value 30000.
func WithConnectTimeoutMS(timeout time.Duration) Option {
	return optionFunc(func(ops *opts) {
		ops.connectTimeoutMS = timeout
	})
}

// WithMaxPoolSize Specifies the maximum number of
// connections that a connection pool may have at a given time.
// Default value 100.
func WithMaxPoolSize(maximum int) Option {
	return optionFunc(func(ops *opts) {
		ops.maxPoolSize = maximum
	})
}

// WithReplicaSet Specifies the replica set name for the cluster.
// All nodes in the replica set must have the same replica set name,
// or the Client will not consider them as part of the set.
// Default value "null".
func WithReplicaSet(replica string) Option {
	return optionFunc(func(ops *opts) {
		ops.replicaSet = replica
	})
}

// WithMaxIdleTimeMS Specifies the maximum amount of time a connection
// can remain idle in the connection pool before being removed and closed.
// The default is 0, meaning a connection can remain unused indefinitely.
func WithMaxIdleTimeMS(idletime time.Duration) Option {
	return optionFunc(func(ops *opts) {
		ops.maxIdleTimeMS = idletime
	})
}

// WithMinPoolSize Specifies the minimum number of connections
// that the driver maintains in a single connection pool.
// Default value 0.
func WithMinPoolSize(minimum int) Option {
	return optionFunc(func(ops *opts) {
		ops.minPoolSize = minimum
	})
}

// WithSocketTimeoutMS Specifies the number of milliseconds to wait
// for a socket read or write to return before returning a network error.
// The 0 default value indicates that there is no timeout.
func WithSocketTimeoutMS(timeout time.Duration) Option {
	return optionFunc(func(ops *opts) {
		ops.socketTimeoutMS = timeout
	})
}

// WithTLS Specifies whether to establish a Transport Layer Security (TLS)
// connection with the instance.
// This is automatically set to true when using a DNS seedlist (SRV) in the connection string.
// You can override this behavior by setting the value to false.
func WithTLS(enabled bool) Option {
	return optionFunc(func(ops *opts) {
		ops.tls = enabled
	})
}

// WithWriteConcern Specifies the write concern. For more information on values,
// see the server documentation on Write Concern opts(https://docs.mongodb.com/manual/reference/write-concern/).
func WithWriteConcern(concern string) Option {
	return optionFunc(func(ops *opts) {
		ops.w = concern
	})
}

// WithDirectConnection Specifies whether to force dispatch all
// operations to the host specified in the connection URI.
// default is false.
func WithDirectConnection(direct bool) Option {
	return optionFunc(func(ops *opts) {
		ops.directConnection = direct
	})
}

func in(source []string, find string) bool {
	for _, item := range source {
		if item == find {
			return true
		}
	}
	return false
}
