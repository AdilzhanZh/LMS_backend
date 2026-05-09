package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/AdilzhanZh/LMS_backend/internal/config"
	"github.com/AdilzhanZh/LMS_backend/internal/handler"
	"github.com/AdilzhanZh/LMS_backend/internal/pkg/logger"
	"github.com/AdilzhanZh/LMS_backend/internal/repository"
	"github.com/AdilzhanZh/LMS_backend/internal/server"
	"github.com/AdilzhanZh/LMS_backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	slogger := logger.New(cfg.LogLevel)
	slog.SetDefault(slogger)

	router, err := buildApp(cfg)

	if err != nil {
		slog.Error("failed to init router", "router", err)
		return
	}
	srv := server.New(router, cfg.Port)
	err = srv.Run()
	if err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}

func buildApp(cfg *config.Config) (*gin.Engine, error) {
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("failed to connect db", "db", err)
		return nil, err
	}

	courseRepo := repository.NewPsgCourseRepo(db)
	lessonRepo := repository.NewPsgLessonRepo(db)
	userRepo := repository.NewPsgUserRepo(db)

	services := &service.Services{
		Course: service.NewCourseService(courseRepo, lessonRepo),
		Lesson: service.NewLessonService(lessonRepo, courseRepo),
		Auth:   service.NewAuthService(userRepo),
	}

	h := handler.NewHandler(services)
	router, err := h.InitRoutes()

	return router, err
}
