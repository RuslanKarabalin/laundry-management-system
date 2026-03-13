package main

import (
	"log"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("Can't sync logger after shutdown", zap.Error(err))
		}
	}()

	sugar := logger.Sugar()

	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		sugar.Fatal("Error reading config file, %s", err)
	}
	viper.SetDefault("APP_PORT", ":8080")
	viper.SetDefault("TG_TOKEN", "")

	tgToken := viper.GetString("TG_TOKEN")
	addr := viper.GetString("APP_PORT")

	sugar.Info(tgToken)
	sugar.Info(addr)

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

	tg.Start()

	// -------------

	app := fiber.New()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	sugar.Fatal(app.Listen(addr))
	// -------------
}
