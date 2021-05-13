package test

import (
	"reflect"
	"strconv"
	"strings"
)

func JointUrl(basic string, query interface{}) (url string) {
	v := reflect.ValueOf(query)
	v = getStructValue(v)
	if !v.IsValid() {
		return basic
	}

	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(basic)
	var start, prev = true, false
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsZero() {
			continue
		}
		if start {
			urlBuilder.WriteByte('?')
			start = false
		}
		if prev {
			urlBuilder.WriteByte('&')
		}
		prev = true
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		uri := tag.Get("uri")
		jointOneQueryUrl(uri, &urlBuilder, v.Field(i))
	}
	return urlBuilder.String()
}

func getStructValue(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Struct:
		return value
	case reflect.Ptr:
		return getStructValue(value.Elem())
	case reflect.Interface:
		return getStructValue(value.Elem())
	default:
		return reflect.Value{}
	}
}

func jointOneQueryUrl(uri string, builder *strings.Builder, fieldValue reflect.Value) {
	switch fieldValue.Kind() {
	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
            elem := fieldValue.Index(i)
            builder.WriteString(uri)
            builder.WriteByte('=')
            writeFieldValue(builder, elem)
            if i != fieldValue.Len() - 1 {
            	builder.WriteByte('&')
			}
		}
	default:
		builder.WriteString(uri)
		builder.WriteByte('=')
		writeFieldValue(builder, fieldValue)
	}
}

func writeFieldValue(builder *strings.Builder, elem reflect.Value) {
	switch elem.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		builder.WriteString(strconv.FormatInt(int64(elem.Int()), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		builder.WriteString(strconv.FormatUint(uint64(elem.Uint()), 10))
	case reflect.String:
		builder.WriteString(elem.String())
	}
}