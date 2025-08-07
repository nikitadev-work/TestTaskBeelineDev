package main

import (
	"TestTaskBeelineDev/internal/config"
	"TestTaskBeelineDev/internal/handlers"
	"TestTaskBeelineDev/internal/logger"
	"TestTaskBeelineDev/internal/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

// Точка входа приложения
func main() {
	//Создаем контекст, который будем слушать для последующей организации Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx); err != nil {
		log.Fatal(err)
	}
}

// Функция для запуска сервера
func runServer(ctx context.Context) error {

	// Загрузка конфигурации из переменных окружения
	cfg := config.Load()

	port := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr: port,
	}

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
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//Ожидаем завершение переданного контекста
	log.Println("Listening on port " + cfg.Port)
	<-ctx.Done()

	log.Println("Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %s", err)
		return err
	}

	log.Println("Server shutdown successfully")

	return nil
}
