package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/klausfun/ToDo_project"
	"github.com/klausfun/ToDo_project/pkg/handler"
	"github.com/klausfun/ToDo_project/pkg/repository"
	"github.com/klausfun/ToDo_project/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// задаем формат JSON для логеров(для удобства)
	logrus.SetFormatter(new(logrus.JSONFormatter))
	// инициализируем конфиги
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	// загрузка переменного окружения из файла .env с помощью функции Load()
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// инициализируем БД, передавая все необходимые значения
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"), // из переменного окружения получаем пароль по имени
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// зависимости
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// инициализируем сервер
	srv := new(todo.Server)
	// получаем порт по ключу
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Fatalf("error occured on db connection close: %s", err.Error())
	}
}

// инициализация конфигурационных файлов
func initConfig() error {
	viper.AddConfigPath("configs") // имя директории
	viper.SetConfigName("config")  // имя файла

	// считывает значения конфигов и записывает их во внутрений объект viper
	return viper.ReadInConfig()
}
