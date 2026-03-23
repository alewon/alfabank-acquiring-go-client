package alfabank

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func encodeForm(v any) (url.Values, error) {
	values := url.Values{}
	if v == nil {
		return values, nil
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return values, nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("form payload must be a struct")
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("form")
		if tag == "" || tag == "-" {
			continue
		}

		name, opts := parseFormTag(tag)
		if name == "" {
			continue
		}

		fv := rv.Field(i)
		if opts["omitempty"] && isZeroValue(fv) {
			continue
		}

		if opts["json"] {
			payload, err := json.Marshal(fv.Interface())
			if err != nil {
				return nil, err
			}
			values.Add(name, string(payload))
			continue
		}

		if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
			for j := 0; j < fv.Len(); j++ {
				values.Add(name, scalarToString(fv.Index(j)))
			}
			continue
		}

		values.Add(name, scalarToString(fv))
	}

	return values, nil
}

func parseFormTag(tag string) (string, map[string]bool) {
	parts := strings.Split(tag, ",")
	opts := make(map[string]bool, len(parts))
	for _, part := range parts[1:] {
		opts[part] = true
	}
	return parts[0], opts
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

func scalarToString(v reflect.Value) string {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		return fmt.Sprint(v.Interface())
	}
}