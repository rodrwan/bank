package accounts

import "log"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

// AccountServiceConfig ...
type AccountServiceConfig struct {
	RabbitMQURL string
	PostgresDNS string
}
