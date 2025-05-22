package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "github.com/bllooop/nameservice/internal/delivery/api"
	applog "github.com/bllooop/nameservice/internal/log"
	"github.com/bllooop/nameservice/internal/repository"
	"github.com/bllooop/nameservice/internal/usecase"
	"github.com/spf13/viper"
)

func Run(cfg repository.Config) error {
	dbpool, err := repository.NewPostgresDB(cfg)
	if err != nil {
		applog.Logger.Error().Err(err).Msg("Не удалось установить соединение с базой данных")
		//        return fmt.Errorf("подключение к БД: %w", err)
		applog.Logger.Fatal().Msg("Произошла ошибка с базой данных")

	}
	applog.Logger.Debug().Msg("База данных успешно подключена")

	migratePath := "./migrations"
	applog.Logger.Debug().Msgf("Running database migrations from path: %s", migratePath)
	if err = repository.RunMigrate(cfg, migratePath); err != nil {
		applog.Logger.Error().Err(err).Msg("")
		//applog.Logger.Fatal().Msg("Возникла ошибка при переносе")
		return fmt.Errorf("ошибка при миграции: %w", err)

	}

	applog.Logger.Debug().Msg("Инициализация слоя репозитория")
	repos := repository.NewRepository(dbpool)
	applog.Logger.Debug().Msg("Инициализация usecase слоя")
	usecases := usecase.NewUsecase(repos)
	applog.Logger.Debug().Msg("Инициализация обработчиков API")
	handler := handlers.NewHandler(usecases)
	srv := new(Server)
	//http serv
	go func() {
		applog.Logger.Info().Msg("Запуск сервера...")
		if err := srv.StartHTTP(cfg.ServerPort, handler.InitRoutes()); err != nil && err != http.ErrServerClosed {
			applog.Logger.Error().Err(err).Msg("При запуске HTTP сервера произошла ошибка")
			//applog.Logger.Fatal().Msg("При запуске HTTP сервера произошла ошибка")
			os.Exit(1)
		} else {
			applog.Logger.Info().Msg("HTTP сервер был закрыт аккуратно")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	applog.Logger.Debug().Msg("Прослушивание сигналов завершения работы ОС")
	<-quit
	applog.Logger.Info().Msg("Сервер отключается")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer dbpool.Close()
	applog.Logger.Debug().Msg("Закрытие соединения с базой данных ")
	if err := srv.Shutdown(ctx); err != nil {
		applog.Logger.Error().Err(err).Msg("")
		//applog.Logger.Fatal().Msg("При выключении сервера произошла ошибка")
		return fmt.Errorf("ошибка остановки сервера: %w", err)

	}
	return nil
}

func initConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
