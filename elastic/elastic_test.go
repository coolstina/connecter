package elastic

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestElasticSuite(t *testing.T) {
	suite.Run(t, &ElasticSuite{})
}

type ElasticSuite struct {
	suite.Suite
	client *elastic.Client
}

func (e *ElasticSuite) BeforeTest(suiteName, testName string) {
	var err error
	e.client, err = NewElasticConnection(
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
		WithSetURL("http://127.0.0.1:9200"),
	)
	assert.NoError(e.T(), err)
	assert.NotNil(e.T(), e.client)
}

func (e *ElasticSuite) Test_GetInstance() {
	actual, err := e.client.IndexExists("hello_world").Do(context.Background())
	assert.NoError(e.T(), err)
	assert.Equal(e.T(), false, actual)
}
