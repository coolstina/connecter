// Copyright 2021 helloshaohua &lt;wu.shaohua@foxmail.com&gt;
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

package mysql

type Option interface {
	apply(*options)
}

type optionFunc func(ops *options)

func (o optionFunc) apply(ops *options) {
	o(ops)
}

type options struct {
	charset   string
	parseTime bool
	location  string
}

func WithCharset(charset string) Option {
	return optionFunc(func(ops *options) {
		ops.charset = charset
	})
}

func WithParseTime(parseTime bool) Option {
	return optionFunc(func(ops *options) {
		ops.parseTime = parseTime
	})
}

func WithLocation(loc string) Option {
	return optionFunc(func(ops *options) {
		ops.location = loc
	})
}
