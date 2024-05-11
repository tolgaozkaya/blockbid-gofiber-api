package routes

import (
	Ihaleler "blockchain-smart-tender-platform/app/controllers/ihaleler"
	Kullanicilar "blockchain-smart-tender-platform/app/controllers/kullanicilar"
	Teklifler "blockchain-smart-tender-platform/app/controllers/teklifler"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// API v1

	app.Get("/network-status", func(c *fiber.Ctx) error {
		info, err := Ihaleler.FetchNetworkInfo()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch network information: %v", err)})
		}
		return c.JSON(info)
	})

	ihaleGroup := app.Group("/api/v1/ihaleler")
	{
		ihaleGroup.Post("/create", Ihaleler.CreateIhale)
		ihaleGroup.Put("/update", Ihaleler.UpdateIhale)
		ihaleGroup.Get("/query/:ihaleNumarasi", Ihaleler.QueryIhale)
		ihaleGroup.Get("/list", Ihaleler.ListIhaleler)
		ihaleGroup.Put("/close/:ihaleNumarasi", Ihaleler.CloseIhale)
		ihaleGroup.Put("/close", Ihaleler.CheckAndCloseIhales)
		ihaleGroup.Get("/user", Ihaleler.QueryKullaniciIhaleleri)
		ihaleGroup.Get("/others", Ihaleler.QueryDigerIhaleler)
		ihaleGroup.Post("/start/:ihaleNumarasi", Ihaleler.StartIhale)
		ihaleGroup.Put("/extend/:ihaleNumarasi", Ihaleler.ExtendIhale)
		ihaleGroup.Get("/history/:ihaleNumarasi", Ihaleler.GetHistoryForIhale)
	}

	teklifGroup := app.Group("/api/v1/teklifler")
	{
		teklifGroup.Post("/:ihaleId/make", Teklifler.MakeTeklif)
		teklifGroup.Get("/:ihaleId/query", Teklifler.QueryTeklif)
		teklifGroup.Get("/:ihaleId/queryAll", Teklifler.QueryAllTeklif)
	}
	kullaniciGroup := app.Group("/api/v1/kullanicilar")
	{
		kullaniciGroup.Get("/get", Kullanicilar.GetUserByID)
		kullaniciGroup.Get("/query/:userID", Kullanicilar.QueryUser)
		kullaniciGroup.Put("/update", Kullanicilar.UpdateUser)
		kullaniciGroup.Delete("/delete", Kullanicilar.DeleteUser)
		kullaniciGroup.Get("/userdata", Kullanicilar.GetUserDashboardData)
	}

}
