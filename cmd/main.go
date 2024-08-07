package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rabboni171/url-shortener/configs"
	"github.com/rabboni171/url-shortener/internal/handler"
	"github.com/rabboni171/url-shortener/internal/repository"
	"github.com/rabboni171/url-shortener/internal/service"
	"github.com/rabboni171/url-shortener/pkg/db"
	"github.com/rabboni171/url-shortener/pkg/logger"
)

func main() {
	// Инициализируем конфигурацию
	cfg, err := configs.InitConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	zeroLog := logger.NewLogger()

	// Подключаемся к Redis
	redisClient := db.NewRedisClient(cfg.DBRedisParams)

	// Пробуем записать и прочитать значение(FOR TEST!!!)
	ctx := context.Background()
	err = redisClient.Set(ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	value, err := redisClient.Get(ctx, "test_key").Result()
	if err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}

	fmt.Printf("Value: %s\n", value)

	// Создаем экземпляры слоев
	repo := repository.NewRedisRepo(redisClient)
	svc := service.NewURLServiceImpl(repo)
	h := handler.NewURLHandler(cfg, zeroLog, svc)

	r := handler.InitRoutes(h)

	server := &http.Server{
		Addr:         cfg.AppParams.ServerURL + ":" + cfg.AppParams.PortRun,
		Handler:      r,
		WriteTimeout: time.Duration(cfg.AppParams.WriteTimeout),
		ReadTimeout:  time.Duration(cfg.AppParams.ReadTimeout),
	}

	// Запуск HTTP-сервера в горутине
	go func() {
		log.Printf("Server is starting on %s", server.Addr)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// Создаем канал для получения сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Блокируем главный поток до получения сигнала завершения

	log.Println("Server is shutting down...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем HTTP-сервер
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	// Сообщаем о завершении graceful shutdown
	log.Println("Server stopped gracefully")
}
