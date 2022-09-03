package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/routes"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	HOST          = "127.0.0.1"
	PORT          = "4000"
)

var (
	PROD       = flag.Bool("prod", false, "Enable prefork in Production")
	app        *fiber.App
)

func loadRoutes() {
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
}

func initLogger() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.ToSlash("./log/logfile.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	}
	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(log.TextFormatter)
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
}

func main() {
	config := fiber.Config{
		Prefork: *PROD, // go run app.go -prod
	}

	initLogger()

	app = fiber.New(config)

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	database.InitDatabase()
	log.Info("Connected to DB!")

	loadRoutes()

	log.Fatal(app.Listen(HOST + ":" + PORT))
}
