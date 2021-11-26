// Copyright 2021 helloshaohua <wu.shaohua@foxmail.com>;
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

package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/olivere/elastic"
)

const (
	DefaultURL = "http://127.0.0.1:9200"
)

// NewConnection initialize elastic client instance for connection.
func NewConnection(ops ...Option) (*elastic.Client, error) {
	opts := &options{}

	for _, o := range ops {
		o.apply(opts)
	}

	fs := make([]elastic.ClientOptionFunc, 0)
	rt := reflect.TypeOf(opts).Elem()
	rv := reflect.ValueOf(opts).Elem()

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			switch rt.Field(i).Name {
			case "httpClient":
			case "snifferEnabled":
				fs = append(fs, elastic.SetSniff(rv.Field(i).Elem().Bool()))
			case "healthCheckEnabled":
				fs = append(fs, elastic.SetHealthcheck(rv.Field(i).Elem().Bool()))
			case "gzipEnabled":
				fs = append(fs, elastic.SetGzip(rv.Field(i).Elem().Bool()))
			case "errorlog":
				fs = append(fs, elastic.SetErrorLog(private(rv, i).Interface().(elastic.Logger)))
			case "infolog":
				fs = append(fs, elastic.SetInfoLog(private(rv, i).Interface().(elastic.Logger)))
			case "tracelog":
				fs = append(fs, elastic.SetTraceLog(private(rv, i).Interface().(elastic.Logger)))
			case "healthCheckTimeoutStartup":
				fs = append(fs, elastic.SetHealthcheckTimeoutStartup(time.Duration(rv.Field(i).Elem().Int())))
			case "healthCheckTimeout":
				fs = append(fs, elastic.SetHealthcheckTimeout(time.Duration(rv.Field(i).Elem().Int())))
			case "healthCheckInterval":
				fs = append(fs, elastic.SetHealthcheckInterval(time.Duration(rv.Field(i).Elem().Int())))
			case "snifferTimeoutStartup":
				fs = append(fs, elastic.SetSnifferTimeoutStartup(time.Duration(rv.Field(i).Elem().Int())))
			case "snifferTimeout":
				fs = append(fs, elastic.SetSnifferTimeout(time.Duration(rv.Field(i).Elem().Int())))
			case "snifferInterval":
				fs = append(fs, elastic.SetSnifferInterval(time.Duration(rv.Field(i).Elem().Int())))
			case "urls":
				fs = append(fs, elastic.SetURL(private(rv, i).Interface().([]string)...))
			case "scheme":
				fs = append(fs, elastic.SetScheme(rv.Field(i).Elem().String()))
			}
		}
	}

	return elastic.NewClient(fs...)
}

// CreateIndexIfNotExists Create elastic mappings if not exists.
func CreateIndexIfNotExists(ctx context.Context, client *elastic.Client, index, body string) (*elastic.IndicesCreateResult, error) {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return nil, err
	}

	if exists {
		resp := &elastic.IndicesCreateResult{
			Acknowledged:       true,
			ShardsAcknowledged: true,
			Index:              index,
		}
		return resp, nil
	}

	resp, err := client.CreateIndex(index).Body(body).Do(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteIndexIfExists delete elastic mappings if exists.
func DeleteIndexIfExists(ctx context.Context, client *elastic.Client, index string) (*elastic.IndicesDeleteResponse, error) {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		resp := &elastic.IndicesDeleteResponse{
			Acknowledged: true,
		}
		return resp, nil
	}

	resp, err := client.DeleteIndex(index).Do(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// BulkInsert execute a new BulkIndexRequest.
// The operation type is "mappings" by default.
func BulkInsert(ctx context.Context, client *elastic.Client, list []string) (*elastic.BulkResponse, error) {
	if list == nil {
		return nil, fmt.Errorf("bulk install data can't empty")
	}

	bulk := client.Bulk()
	requests := make([]elastic.BulkableRequest, 0, len(list))

	for _, item := range list {
		var r elastic.GetResult

		err := json.Unmarshal([]byte(item), &r)
		if err != nil {
			return nil, err
		}

		source, err := r.Source.MarshalJSON()
		if err != nil {
			return nil, err
		}

		request := elastic.NewBulkIndexRequest().
			Index(r.Index).
			Type(r.Type).
			Id(r.Id).
			Doc(string(source))

		requests = append(requests, request)
	}

	// Do sends the bulk requests to Elasticsearch
	resp, err := bulk.Add(requests...).Do(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Access unexported struct fields.
func private(value reflect.Value, index int) reflect.Value {
	fi := value.Field(index)
	return reflect.NewAt(fi.Type(), unsafe.Pointer(fi.UnsafeAddr())).Elem()
}
