package config

// import ("os" "strconv" )

import (
	"os"
	"strconv"

	logrus "github.com/sirupsen/logrus"
)

const (
	LogLevelEnvVar       = "LOG_LEVEL"
	portEnvVar           = "PORT"
	BrokerAddressEnvVar  = "BROKER_ADDRESS"
	defaultLogLevel      = logrus.DebugLevel
	defaultPort          = 8080
	defaultBrokerAddress = "localhost"
)

func LogLevel() logrus.Level {
	var (
		level logrus.Level
		err   error
	)

	if level, err = logrus.ParseLevel(os.Getenv(LogLevelEnvVar)); err != nil {
		return defaultLogLevel
	}

	return level
}

// Port returns the poer the service should listen on, or 3000 if not defined or\
// is not a valid port
func Port() int {
	var (
		rawPort string
		found   bool
		port    int
		err     error
	)

	if rawPort, found = os.LookupEnv(portEnvVar); !found {
		return defaultPort
	}

	if port, err = strconv.Atoi(rawPort); err != nil {
		return defaultPort
	}

	return port
}

// BrokerAddress returns the address the kafka broker is listeneing on, or localohost if not defined
func BrokerAddress() string {
	var brokerAddress string
	var found bool

	if brokerAddress, found = os.LookupEnv(BrokerAddressEnvVar); !found {
		return defaultBrokerAddress
	}

	if len(brokerAddress) == 0 {
		return defaultBrokerAddress
	}

	return brokerAddress
}
