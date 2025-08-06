package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	AuthKey          string
	Port             string
	OutputServiceURL string
	LogFile          string
	NumOfWorkers     int
	Timeout          time.Duration
}

// Load загружает конфигурацию из переменных окружения
// Для параметров, не указанных в окружении, используются значения по умолчанию
func Load() Config {
	log.Println("Loading config...")

	//Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := getEnv("PORT", "8080")
	outputServiceURL := getEnv("OUTPUT_SERVICE_URL", "http://localhost:8081/users")
	authKey := getEnv("AUTH_KEY", "")
	logFile := getEnv("LOG_FILE", "provider.log")
	numOfWorkers, err := strconv.Atoi(getEnv("NUM_OF_WORKERS", "10"))
	if err != nil {
		numOfWorkers = 10
		log.Println("NUM_OF_WORKERS env variable is incorrect")
	}

	timeoutStr := getEnv("TIMEOUT", "10s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeout = 10 * time.Second
		log.Println("TIMEOUT env variable is incorrect")
	}

	if authKey == "" {
		log.Fatal("AUTH_KEY is required in configuration")
	}

	log.Println("Loading config finished")
	return Config{
		AuthKey:          authKey,
		Port:             port,
		OutputServiceURL: outputServiceURL,
		LogFile:          logFile,
		NumOfWorkers:     numOfWorkers,
		Timeout:          timeout,
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
// если переменная не установлена
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
