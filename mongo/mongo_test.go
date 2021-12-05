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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestMongoSuite(t *testing.T) {
	suite.Run(t, &MongoSuite{})
}

type MongoSuite struct {
	suite.Suite
	err    error
	client *mongo.Client
	uri    string
	ctx    context.Context
	cancel context.CancelFunc
}

func (suite *MongoSuite) BeforeTest(suiteName, testName string) {
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
	suite.uri = "mongodb://root:root@localhost:27017/?maxPoolSize=20&w=majority"

	suite.client, suite.err = NewConnection("127.0.0.1:27017", "root", "root")
	assert.NoError(suite.T(), suite.err)
	assert.NotNil(suite.T(), suite.client)
}

func (suite *MongoSuite) Test_Ping() {
	err := suite.client.Ping(suite.ctx, readpref.Primary())
	assert.NoError(suite.T(), err)
}

func (suite *MongoSuite) Test_uri() {
	expected := `mongodb://root:root@localhost:27017/?connectTimeoutMS=30000000000&maxPoolSize=100&replicaSet=null&maxIdleTimeMS=0&minPoolSize=0&socketTimeoutMS=30000000000&serverSelectionTimeoutMS=10000000000&tls=false&w=null&directConnection=true`

	actual, err := uri(option("localhost:27017", "root", "root"))
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), actual)
	assert.Equal(suite.T(), expected, actual)
}
