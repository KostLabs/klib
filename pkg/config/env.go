package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// LoadEnv reads environment variables into one or more config structs.
// Each exported field must carry an `env:"VAR_NAME"` tag. Fields without
// the tag are skipped. If the environment variable is empty or unset the
// field is left unchanged, so YAML defaults set beforehand are preserved.
//
// Supported field types: string, int, bool, time.Duration.
//
// Example:
//
//	type AppConfig struct {
//	    APIKey        string        `env:"API_KEY"`
//	    RedisPassword string        `env:"REDIS_PASSWORD"`
//	    Timeout       time.Duration `env:"TIMEOUT"`
//	}
//
//	var cfg AppConfig
//	if err := config.LoadEnv(&cfg); err != nil { ... }
func LoadEnv(targets ...any) error {
	for _, target := range targets {
		if err := loadEnvInto(target); err != nil {
			return err
		}
	}
	return nil
}

func loadEnvInto(target any) error {
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config.LoadEnv: target must be a pointer to a struct, got %T", target)
	}

	structValue := value.Elem()
	structType := structValue.Type()

	for fieldIndex := range structType.NumField() {
		field := structType.Field(fieldIndex)
		fieldValue := structValue.Field(fieldIndex)

		// Recurse into nested structs (with or without an env tag).
		if field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Duration(0)) {
			if err := loadEnvInto(fieldValue.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		envKey, ok := field.Tag.Lookup("env")
		if !ok || envKey == "" {
			continue
		}

		rawValue, set := os.LookupEnv(envKey)
		if !set || rawValue == "" {
			continue
		}

		if err := setField(fieldValue, field.Name, rawValue); err != nil {
			return fmt.Errorf("config.LoadEnv: field %s (env %q): %w", field.Name, envKey, err)
		}
	}

	return nil
}

func setField(fieldValue reflect.Value, fieldName, rawValue string) error {
	switch fieldValue.Interface().(type) {
	case string:
		fieldValue.SetString(rawValue)

	case int, int8, int16, int32, int64:
		parsedInt, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as int: %w", rawValue, err)
		}
		fieldValue.SetInt(parsedInt)

	case bool:
		parsedBool, err := strconv.ParseBool(rawValue)
		if err != nil {
			return fmt.Errorf("cannot parse %q as bool: %w", rawValue, err)
		}
		fieldValue.SetBool(parsedBool)

	case time.Duration:
		parsedDuration, err := time.ParseDuration(rawValue)
		if err != nil {
			return fmt.Errorf("cannot parse %q as duration: %w", rawValue, err)
		}
		fieldValue.SetInt(int64(parsedDuration))

	default:
		return fmt.Errorf("unsupported field type %s for %s", fieldValue.Type(), fieldName)
	}

	return nil
}
