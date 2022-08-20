package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadENVData() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Fail to load the environment configuration")
	}
}
