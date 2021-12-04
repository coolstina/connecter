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

package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	connection, err := NewConnection(
		NewDefaultSimpleConfig("localhost:6379", "", 8),
		WithMaxRetries(3),
	)
	assert.NoError(t, err)
	assert.NotNil(t, connection)

	// Set key value.
	set := connection.Set("username", "helloshaohua", time.Second*5)
	assert.NotNil(t, set)
	actual, err := set.Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", actual)

	// Get key value.
	actual, err = connection.Get("username").Result()
	assert.NoError(t, err)
	assert.Equal(t, "helloshaohua", actual)
}
