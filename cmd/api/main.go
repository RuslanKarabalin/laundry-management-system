package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("Can't sync logger after shutdown", zap.Error(err))
		}
	}()

	sugar := logger.Sugar()

	viper.AutomaticEnv()
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		sugar.Warn("Error reading config file, %s", err)
	}
	viper.SetDefault("APP_PORT", ":8080")
	viper.SetDefault("TG_TOKEN", "")

	tgToken := viper.GetString("TG_TOKEN")
	addr := viper.GetString("APP_PORT")
	pgUsername := viper.GetString("POSTGRES_USER")
	pgPassword := viper.GetString("POSTGRES_PASSWORD")
	pgHost := viper.GetString("POSTGRES_HOST")
	pgPort := viper.GetString("POSTGRES_PORT")
	pgBasename := viper.GetString("POSTGRES_DB")

	pgUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		pgUsername,
		pgPassword,
		pgHost,
		pgPort,
		pgBasename,
	)

	sugar.Info(tgToken)
	sugar.Info(addr)

	conn, err := pgxpool.New(ctx, pgUrl)
	if err != nil {
		sugar.Error("Cannot connect to PostgreSQL", zap.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	// -------------
	pref := telebot.Settings{
		Token:  tgToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	tg, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	tg.Handle("/hello", func(c telebot.Context) error {
		return c.Send("Hello!")
	})

	go tg.Start()

	// -------------

	app := fiber.New()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := conn.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"db":     "unreachable",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     "reachable",
		})
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	sugar.Fatal(app.Listen(addr))
	// -------------
}
