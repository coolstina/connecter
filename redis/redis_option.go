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

import "time"

type Option interface {
	apply(config *Config)
}

type optionFunc func(config *Config)

func (o optionFunc) apply(config *Config) {
	o(config)
}

// WithNetwork the network type, either tcp or unix.
// Default is tcp.
func WithNetwork(network string) Option {
	return optionFunc(func(config *Config) {
		config.Network = network
	})
}

// WithMaxRetries Specify maximum number of retries before giving up.
// Default is to not retry failed commands.
func WithMaxRetries(times int) Option {
	return optionFunc(func(config *Config) {
		config.MaxRetries = times
	})
}

// WithMinRetryBackoff Specify minimum backoff between each retry.
// Default is 8 milliseconds; -1 disables backoff.
func WithMinRetryBackoff(wait time.Duration) Option {
	return optionFunc(func(config *Config) {
		config.MinRetryBackoff = wait
	})
}

// WithMaxRetryBackoff Specify maximum backoff between each retry.
// Default is 512 milliseconds; -1 disables backoff.
func WithMaxRetryBackoff(wait time.Duration) Option {
	return optionFunc(func(config *Config) {
		config.MaxRetryBackoff = wait
	})
}

// WithDialTimeout Specify dial timeout for establishing new connections.
// Default is 5 seconds.
func WithDialTimeout(timeout time.Duration) Option {
	return optionFunc(func(config *Config) {
		config.DialTimeout = timeout
	})
}

// WithReadTimeout Specify Timeout for socket reads. If reached, commands will fail
// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
// Default is 3 seconds.
func WithReadTimeout(timeout time.Duration) Option {
	return optionFunc(func(config *Config) {
		config.ReadTimeout = timeout
	})
}
