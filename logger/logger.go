package logger

import "log"

func D(message string, params ...interface{}) {
	printf("DEBUG", message, params...)
}

func W(message string, params ...interface{}) {
	printf("WARNING", message, params...)
}

func E(message string, params ...interface{}) {
	printf("ERROR", message, params...)
}

func I(message string, params ...interface{}) {
	printf("INFO", message, params...)
}

func F(message string, params ...interface{}) {
	printf("FATAL", message, params...)
}

func printf(level string, message string, params ...interface{}) {
	if level == "FATAL" {
		log.Fatalf(level + ": " + message + "\n", params...)
	} else {
		log.Printf(level + ": " + message + "\n", params...)
	}
}
