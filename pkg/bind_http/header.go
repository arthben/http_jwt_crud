package bindhttp

// Header Validation taken from :
// https://github.com/dennisstritzke/httpheader/blob/master/bind.go

import (
	"fmt"
	"net/http"
	"reflect"
)

const (
	tagIdentifier = "header"
)

// Bind processes the HTTP header fields and stores the result in the value pointed to by v.
func BindHeader(header http.Header, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		msg := reflect.TypeOf(v).String()
		return fmt.Errorf(msg)
	}

	return bind(v, header)
}

func bind(v interface{}, header http.Header) error {
	rv := reflect.Indirect(reflect.ValueOf(v))
	t := reflect.TypeOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if headerName, ok := field.Tag.Lookup(tagIdentifier); ok {
			headerValue := header[headerName]

			if len(headerValue) > 0 {
				value := rv.Field(i)
				setValue(field, value, headerValue)
			}
		}
	}

	return nil
}

func setValue(field reflect.StructField, value reflect.Value, headerValue []string) {
	switch field.Type.Kind() {
	case reflect.String:
		setStringValue(value, headerValue)
	case reflect.Slice:
		setSliceValue(value, headerValue)
	}
}

func setSliceValue(value reflect.Value, headerValue []string) {
	headerValueCount := len(headerValue)
	slice := reflect.MakeSlice(reflect.TypeOf([]string{}), headerValueCount, headerValueCount)
	for i := 0; i < headerValueCount; i++ {
		sliceItem := slice.Index(i)
		sliceItem.SetString(headerValue[i])
	}
	value.Set(slice)
}

func setStringValue(value reflect.Value, headerValue []string) {
	value.SetString(headerValue[0])
}
