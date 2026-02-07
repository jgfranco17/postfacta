package environment

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	APPLICATION_ENV_LOCAL = "local"
	APPLICATION_ENV_DEV   = "dev"
	APPLICATION_ENV_STAGE = "stage"
	APPLICATION_ENV_PROD  = "prod"
)

const (
	ENV_KEY_ENVIRONMENT = "ENVIRONMENT"
	ENV_KEY_VERSION     = "APP_VERSION"
	ENV_KEY_JWT_SECRET  = "POSTFACTA_JWT_SECRET"
	ENV_KEY_DB_URL      = "POSTFACTA_DB_URL"
	ENV_KEY_DB_KEY      = "POSTFACTA_DB_KEY"
	ENV_KEY_LOG_LEVEL   = "LOG_LEVEL"
	ENV_LOG_FORMAT      = "LOG_FORMAT"
)

func IsRunningLocally() bool {
	return GetApplicationEnv() == APPLICATION_ENV_LOCAL
}

func GetEnvWithDefault(key string, defaultValue string) string {
	value, present := os.LookupEnv(key)
	if present {
		return value
	}
	return defaultValue
}

func GetApplicationEnv() string {
	return GetEnvWithDefault(ENV_KEY_ENVIRONMENT, APPLICATION_ENV_LOCAL)
}

func GetLogLevel() logrus.Level {
	appEnv := GetEnvWithDefault(ENV_KEY_LOG_LEVEL, "INFO")
	stringToLogLevel := map[string]logrus.Level{
		"DEBUG": logrus.DebugLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"PANIC": logrus.PanicLevel,
		"FATAL": logrus.FatalLevel,
		"TRACE": logrus.TraceLevel,
	}

	level, exists := stringToLogLevel[strings.ToUpper(appEnv)]
	if !exists {
		return logrus.InfoLevel
	}
	return level
}

func GetLogFormatter() logrus.Formatter {
	format := GetEnvWithDefault(ENV_LOG_FORMAT, "DEFAULT")
	switch strings.ToUpper(format) {
	case "JSON":
		return &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
		}
	case "TEXT":
		return &logrus.TextFormatter{
			DisableColors:          false,
			PadLevelText:           true,
			QuoteEmptyFields:       true,
			FullTimestamp:          true,
			DisableSorting:         true,
			DisableLevelTruncation: true,
			TimestampFormat:        time.DateTime,
		}
	}
	return &logrus.TextFormatter{
		DisableColors:    false,
		PadLevelText:     true,
		QuoteEmptyFields: true,
		DisableSorting:   true,
		FullTimestamp:    true,
		TimestampFormat:  time.RFC822,
	}
}
