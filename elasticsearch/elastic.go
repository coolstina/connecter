package elasticsearch

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/olivere/elastic"
)

const (
	DefaultURL = "http://127.0.0.1:9200"
)

func NewConnection(ops ...Option) (*elastic.Client, error) {
	opts := &options{}
	for _, o := range ops {
		o.apply(opts)
	}

	fs := make([]elastic.ClientOptionFunc, 0)
	tf := reflect.TypeOf(opts).Elem()
	of := reflect.ValueOf(opts).Elem()

	for i := 0; i < of.NumField(); i++ {
		if !of.Field(i).IsNil() {
			switch tf.Field(i).Name {
			case "httpClient":
			case "snifferEnabled":
				fs = append(fs, elastic.SetSniff(of.Field(i).Elem().Bool()))
			case "healthCheckEnabled":
				fs = append(fs, elastic.SetHealthcheck(of.Field(i).Elem().Bool()))
			case "gzipEnabled":
				fs = append(fs, elastic.SetGzip(of.Field(i).Elem().Bool()))
			case "errorlog":
				fs = append(fs, elastic.SetErrorLog(private(of, i).Interface().(elastic.Logger)))
			case "infolog":
				fs = append(fs, elastic.SetInfoLog(private(of, i).Interface().(elastic.Logger)))
			case "tracelog":
				fs = append(fs, elastic.SetTraceLog(private(of, i).Interface().(elastic.Logger)))
			case "healthCheckTimeoutStartup":
				fs = append(fs, elastic.SetHealthcheckTimeoutStartup(time.Duration(of.Field(i).Elem().Int())))
			case "healthCheckTimeout":
				fs = append(fs, elastic.SetHealthcheckTimeout(time.Duration(of.Field(i).Elem().Int())))
			case "healthCheckInterval":
				fs = append(fs, elastic.SetHealthcheckInterval(time.Duration(of.Field(i).Elem().Int())))
			case "snifferTimeoutStartup":
				fs = append(fs, elastic.SetSnifferTimeoutStartup(time.Duration(of.Field(i).Elem().Int())))
			case "snifferTimeout":
				fs = append(fs, elastic.SetSnifferTimeout(time.Duration(of.Field(i).Elem().Int())))
			case "snifferInterval":
				fs = append(fs, elastic.SetSnifferInterval(time.Duration(of.Field(i).Elem().Int())))
			case "urls":
				fs = append(fs, elastic.SetURL(private(of, i).Interface().([]string)...))
			case "scheme":
				fs = append(fs, elastic.SetScheme(of.Field(i).Elem().String()))
			}
		}
	}

	return elastic.NewClient(fs...)
}

// Access unexported struct fields.
func private(value reflect.Value, index int) reflect.Value {
	fi := value.Field(index)
	return reflect.NewAt(fi.Type(), unsafe.Pointer(fi.UnsafeAddr())).Elem()
}
