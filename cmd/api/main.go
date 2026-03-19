package main

import (
	"context"
	_ "embed"
	"laundry-management-system/internal/config"
	"laundry-management-system/internal/db"
	"laundry-management-system/internal/model"
	"laundry-management-system/internal/repository"
	"os"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed static/index.html
var indexHTML []byte

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("Can't sync logger after shutdown", zap.Error(err))
		}
	}()

	sugar := logger.Sugar()

	cfg := config.ReadConfig(sugar)

	conn, err := pgxpool.New(ctx, cfg.GetPostgresUrl())
	if err != nil {
		sugar.Error("Cannot connect to PostgreSQL", zap.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	goose.SetLogger(zap.NewStdLog(logger))

	db.RunMigrations(conn)

	r := repository.Init(conn)

	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)

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

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.Send(indexHTML)
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/appliances", func(c *fiber.Ctx) error {
		var appliances []*model.Appliance
		var err error

		applianceType := c.Query("type")

		switch applianceType {
		case "washing_machine":
			appliances, err = r.GetWashingMachines(c.Context())
		case "tumble_dryer":
			appliances, err = r.GetTumbleDryers(c.Context())
		default:
			appliances, err = r.GetAppliances(c.Context())
		}
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(appliances)
	})

	app.Get("/appliances/:id/reservations", func(c *fiber.Ctx) error {
		applianceId := c.Params("id")

		applianceUuid, err := uuid.Parse(applianceId)
		if err != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		reservations, err := r.GetReservationsByApplianceId(c.Context(), applianceUuid)
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(reservations)
	})

	app.Post("/appliances/:id/reservations", func(c *fiber.Ctx) error {
		applianceId := c.Params("id")

		applianceUuid, err := uuid.Parse(applianceId)
		if err != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		t := &model.CreateReservation{}

		if err := c.BodyParser(t); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err = r.CreateReservationsByApplianceId(c.Context(), applianceUuid, t)
		if err != nil {
			return c.SendStatus(500)
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	app.Get("/reservations", func(c *fiber.Ctx) error {
		reservations, err := r.GetReservations(c.Context())
		if err != nil {
			return c.SendStatus(500)
		}
		return c.JSON(reservations)
	})

	sugar.Fatal(app.Listen(cfg.Addr))
}
