package main

import (
	"TestTaskBeelineDev/internal/config"
	"TestTaskBeelineDev/internal/handlers"
	"TestTaskBeelineDev/internal/logger"
	"TestTaskBeelineDev/internal/middleware"
	"fmt"
	"log"
	"net/http"
)

// Точка входа приложения
func main() {
	// Загрузка конфигурации из переменных окружения
	cfg := config.Load()

	// Инициализация логгера с указанием файла для записи логов
	logger.InitLogger(cfg.LogFile)
	log.Println("Logger configured successfully")

	// Инициализация middleware с секретным ключом аутентификации
	middleware.InitMiddleware(cfg.AuthKey)
	log.Println("Config read successfully")

	// Создание экземпляра обработчика конвертации
	handler := &handlers.ConverterHandler{
		Cfg: &cfg,
	}

	// Регистрация обработчика для пути /convert_xml
	// с применением middleware аутентификации
	http.HandleFunc("/convert_xml", middleware.AuthMiddleware(handler.XmlConverterHandler))

	// Запуск HTTP-сервера на указанном порту
	log.Println("Listening on port " + cfg.Port)
	port := fmt.Sprintf(":%s", cfg.Port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Error while listening on port " + cfg.Port)
		return
	}
}
