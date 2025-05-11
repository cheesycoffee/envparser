package envparser

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Parse(target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}
	v := val.Elem()
	t := v.Type()

	var errs []error

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		tag := fieldType.Tag
		envKey := tag.Get("env")

		// Handle embedded/anonymous structs
		if fieldType.Anonymous || (fieldType.Type.Kind() == reflect.Struct && (envKey == "" || envKey == "-")) {
			if err := Parse(field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if envKey == "" || envKey == "-" {
			continue
		}

		val, ok := os.LookupEnv(envKey)
		if !ok {
			return fmt.Errorf("missing %s environment", envKey)
		}

		if err := setValueFromEnv(field, fieldType, val); err != nil {
			errs = append(errs, fmt.Errorf("env '%s': %v", envKey, err))
		}
	}

	if len(errs) > 0 {
		var builder strings.Builder
		builder.WriteString("error parsing environment to struct:\n")
		for _, err := range errs {
			builder.WriteString(err.Error() + "\n")
		}
		return errors.New(builder.String())
	}

	return nil
}

func setValueFromEnv(field reflect.Value, fieldType reflect.StructField, val string) error {
	switch field.Interface().(type) {
	case time.Duration:
		d, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(d))

	case time.Time:
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(t))

	case int, int32, int64:
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))

	case uint, uint32, uint64:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)

	case float32, float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)

	case bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case string:
		field.SetString(val)

	case []string:
		field.Set(reflect.ValueOf(strings.Split(val, ",")))

	case []int:
		numStrings := strings.Split(val, ",")
		ints := make([]int, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				return err
			}
			ints[i] = n
		}
		field.Set(reflect.ValueOf(ints))

	case []int32:
		numStrings := strings.Split(val, ",")
		ints := make([]int32, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				return err
			}
			ints[i] = int32(n)
		}
		field.Set(reflect.ValueOf(ints))

	case []int64:
		numStrings := strings.Split(val, ",")
		ints := make([]int64, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err != nil {
				return err
			}
			ints[i] = n
		}
		field.Set(reflect.ValueOf(ints))

	case []float32:
		numStrings := strings.Split(val, ",")
		float := make([]float32, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.ParseFloat(strings.TrimSpace(v), 32)
			if err != nil {
				return err
			}
			float[i] = float32(n)
		}
		field.Set(reflect.ValueOf(float))

	case []float64:
		numStrings := strings.Split(val, ",")
		float := make([]float64, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil {
				return err
			}
			float[i] = n
		}
		field.Set(reflect.ValueOf(float))

	case []uint:
		numStrings := strings.Split(val, ",")
		unsigned := make([]uint, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
			if err != nil {
				return err
			}
			unsigned[i] = uint(n)
		}
		field.Set(reflect.ValueOf(unsigned))

	case []uint32:
		numStrings := strings.Split(val, ",")
		unsigned := make([]uint32, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 32)
			if err != nil {
				return err
			}
			unsigned[i] = uint32(n)
		}
		field.Set(reflect.ValueOf(unsigned))

	case []uint64:
		numStrings := strings.Split(val, ",")
		unsigned := make([]uint64, len(numStrings))
		for i, v := range numStrings {
			n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
			if err != nil {
				return err
			}
			unsigned[i] = uint64(n)
		}
		field.Set(reflect.ValueOf(unsigned))

	default:
		switch fieldType.Tag.Get("encoding") {
		case "json":
			return json.Unmarshal([]byte(val), field.Addr().Interface())
		case "xml":
			return xml.Unmarshal([]byte(val), field.Addr().Interface())
		case "form":
			parsed, err := url.ParseQuery(val)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(parsed)) // if field is url.Values
		case "base64":
			decoded, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(decoded)) // if field is []byte
		}
	}
	return nil
}
