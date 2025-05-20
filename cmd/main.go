package main

import (
	"fmt"
	"os"

	applog "github.com/bllooop/nameservice/internal/log"
	running "github.com/bllooop/nameservice/internal/server"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/spf13/viper"
)

func main() {
	applog.InitLogger(os.Stdout, zerolog.DebugLevel)
	err := running.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
