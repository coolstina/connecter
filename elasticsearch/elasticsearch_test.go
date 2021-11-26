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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/coolstina/fsfire"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestElasticSuite(t *testing.T) {
	suite.Run(t, &ElasticSuite{})
}

type ElasticSuite struct {
	suite.Suite
	err         error
	client      *elastic.Client
	ctx         context.Context
	baseDir     string
	testDataDir string
	dropIndex   bool
}

func (e *ElasticSuite) BeforeTest(suiteName, testName string) {
	e.ctx = context.TODO()
	e.baseDir = "./../"
	e.dropIndex = true
	e.testDataDir = fsfire.MustGetFilePathWithFSPath(
		e.baseDir, fsfire.WithSpecificFSPath("test/data/elastic"),
	)

	ops := []Option{
		WithGzip(false),
		WithHealthCheck(false),
		WithSniff(false),
		WithErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		WithInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		WithTraceLog(log.New(os.Stdout, "", log.LstdFlags)),
		WithHealthCheckTimeoutStartup(5),
		WithHealthCheckTimeout(5),
		WithHealthCheckInterval(5),
		WithSnifferTimeoutStartup(5),
		WithSnifferTimeout(5),
		WithSnifferInterval(5),
		WithScheme("http"),
		WithSetURL("http://127.0.0.1:9600"),
	}

	e.client, e.err = NewConnection(ops...)
	assert.NoError(e.T(), e.err)
	assert.NotNil(e.T(), e.client)

	// Create mappings
	data, err := ioutil.ReadFile(filepath.Join(e.testDataDir, "mappings/hello_world.json"))
	assert.NoError(e.T(), err)

	resp, err := CreateIndexIfNotExists(e.ctx, e.client, "hello_world", string(data))
	assert.NoError(e.T(), err)
	assert.NotNil(e.T(), resp)
}

func (e *ElasticSuite) AfterTest(suiteName, testName string) {
	if e.dropIndex {
		resp, err := DeleteIndexIfExists(e.ctx, e.client, "hello_world")
		assert.NoError(e.T(), err)
		assert.NotNil(e.T(), resp)
	}
}

func (e *ElasticSuite) Test_BulkInstall() {
	data, err := fsfire.GetFileContentWithStringSlice(filepath.Join(e.testDataDir, "data/hello_world.json"))
	assert.NoError(e.T(), err)
	assert.NotNil(e.T(), data)

	insert, err := BulkInsert(e.ctx, e.client, data)
	assert.NoError(e.T(), err)
	assert.False(e.T(), insert.Errors)
}
