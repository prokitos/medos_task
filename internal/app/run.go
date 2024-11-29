package app

import (
	"fmt"
	"log"
	"mymod/internal/transport"
	"time"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	Server *fiber.App
}

func (a *App) NewServer(port string) {
	app := fiber.New()
	a.Server = app
	a.setHandler()
	a.launchServer(port)
}

func (a *App) Stop() {
	fmt.Println("Gracefully shutting down...")
	a.Server.ShutdownWithTimeout(50 * time.Second)
}

func (a *App) setHandler() {
	transport.SetHandlers(a.Server)
}

func (a *App) launchServer(port string) {
	log.Fatal(a.Server.Listen(port))
}
