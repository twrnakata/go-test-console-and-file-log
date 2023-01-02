package main

import (
	"gofiber_withlog/logs"
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// ดึงเวลามาใช้เป็นชื่อไฟล์ log
	now := time.Now().Format("2006-01-02")
	logFileName := "pathlog-" + now + ".log"

	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	// ถ้าต้องการให้แสดงผลทั้ง console และ log
	output := io.MultiWriter(os.Stdout, logFile)

	app := fiber.New()

	v1 := app.Group("/log", func(c *fiber.Ctx) error {
		// Set หรือ Get Header
		c.Set("Version", "v1")
		return c.Next()
	})

	// Add the logger middleware
	v1.Use("/format", logger.New(logger.Config{
		Format: "${time} ${method} ${path} ${ip} ${status} ${responseTime} ${error}\n",
		Output: output,
	}))
	// 16:46:50 GET /log/format 127.0.0.1 200  -

	v1.Get("/format", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1.Use("ipandport", logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Output: output,
	}))
	// [127.0.0.1]:50414 200 - GET /log/ipandport

	v1.Get("/ipandport", func(c *fiber.Ctx) error {
		return c.SendString("Hello, ip and port!")
	})

	ports := ":8000"
	logs.Info("Banking service started at port " + ports)
	app.Listen(ports)

}
