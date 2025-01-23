package main

import (
	"github.com/joho/godotenv"
	"github.com/sajad-dev/eda-architecture/internal/exception"
)

func main() {
	err := godotenv.Load(".env")
	exception.Log(err)
	
}
