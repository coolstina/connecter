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
	"reflect"
	"time"

	"github.com/go-redis/redis"
)

// NewConnection create a new gorm db instance with the given options.
func NewConnection(config *Config, ops ...Option) (*redis.Client, error) {
	configure := configuration(config)

	for _, o := range ops {
		o.apply(configure)
	}

	return redis.NewClient(options(configure)), nil
}

func configuration(config *Config) *Config {
	configure := &Config{
		Network:            "tcp",
		Host:               config.Host,
		Password:           config.Password,
		Database:           config.Database,
		MinRetryBackoff:    time.Millisecond * 8,
		MaxRetryBackoff:    time.Millisecond * 512,
		DialTimeout:        time.Second * 5,
		ReadTimeout:        time.Second * 3,
		WriteTimeout:       time.Second * 3,
		MinIdleConns:       10,
		MaxConnAge:         time.Second * 3600 * 8, // 8 hours.
		PoolTimeout:        time.Second * 4,
		IdleTimeout:        time.Minute * 5,
		IdleCheckFrequency: time.Minute,
	}
	return configure
}

func options(configure *Config) *redis.Options {
	options := &redis.Options{}
	ov := reflect.ValueOf(options).Elem()
	ct := reflect.TypeOf(configure).Elem()
	cv := reflect.ValueOf(configure).Elem()

	for i := 0; i < cv.NumField(); i++ {
		if !cv.Field(i).IsValid() {
			continue
		}

		switch ct.Field(i).Name {
		case "Host":
			value := ov.FieldByName("Addr")
			if value.IsValid() && value.CanSet() {
				if value.Kind() == reflect.String {
					value.SetString(cv.Field(i).String())
				}
			}
		case "Database":
			value := ov.FieldByName("DB")
			if value.IsValid() && value.CanSet() {
				if value.Kind() == reflect.Int {
					value.SetInt(cv.Field(i).Int())
				}
			}
		default:
			if cv.Field(i).IsZero() {
				continue
			}

			field := ov.FieldByName(ct.Field(i).Name)

			if field.IsValid() && field.CanSet() {
				if field.Kind() == reflect.String {
					field.SetString(cv.Field(i).String())
				}

				if field.Kind() == reflect.Int64 {
					field.SetInt(cv.Field(i).Int())
				}

				if field.Kind() == reflect.Int {
					field.SetInt(cv.Field(i).Int())
				}
			}
		}
	}
	return options
}
