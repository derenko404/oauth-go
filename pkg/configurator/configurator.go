package configurator

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	ENV_TAG         = "env"
	ENV_DEFAULT_TAG = "env_default"
)

type Options struct {
	File string
}

func loadEnvFile(file string) error {
	if file != "" {
		err := godotenv.Load(file)

		if err != nil {
			return err
		}
	} else {
		err := godotenv.Load()

		if err != nil {
			return err
		}
	}

	return nil
}

func getEnvValue(key, defaultValue string) (string, error) {
	envValue := os.Getenv(key)

	if envValue == "" {
		if defaultValue != "" {
			return defaultValue, nil
		} else {
			return "", fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	return envValue, nil
}

func createTransformationError(typ string, tag string, err error) error {
	return fmt.Errorf("invalid %s value for %s: %v", typ, tag, err)
}

// Load parses environment variables into the provided struct based on tags.
// It supports default values using the `env_default` tag.
// Supported types are: string, int, float64, and bool.
//
// Example:
//
//	type Config struct {
//		AppPort string `env:"APP_PORT" env_default:"8080"`
//		AppHost string `env:"APP_HOST" env_default:"localhost"`
//	}
//
//	var config Config
//	err := configurator.Load(&config, &configurator.Options{})
func Load(config any, options *Options) error {
	if config == nil {
		return fmt.Errorf("config is nil %v", config)
	}

	err := loadEnvFile(options.File)

	if err != nil {
		return fmt.Errorf("error loading env file: %w", err)
	}

	value := reflect.ValueOf(config)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for i := range value.NumField() {
		fieldType := value.Type().Field(i)
		fieldTag := fieldType.Tag.Get(ENV_TAG)
		defaultValue := fieldType.Tag.Get(ENV_DEFAULT_TAG)

		envValue, err := getEnvValue(fieldTag, defaultValue)

		if err != nil {
			return fmt.Errorf("error getting env value for field %s: %v", fieldTag, err)
		}

		fieldVal := value.Field(i)

		switch fieldType.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(envValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				return createTransformationError("int", fieldTag, err)
			}

			fieldVal.SetInt(int64(intValue))
		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return createTransformationError("float", fieldTag, err)
			}

			fieldVal.SetFloat(floatValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return createTransformationError("bool", fieldTag, err)
			}

			fieldVal.SetBool(boolValue)
		default:
			return fmt.Errorf("unsupported field type for %s", fieldTag)
		}
	}

	return err
}
