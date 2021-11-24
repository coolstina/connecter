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
