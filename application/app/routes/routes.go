package routes

import (
	Ihaleler "blockchain-smart-tender-platform/app/controllers/ihaleler"
	Kullanicilar "blockchain-smart-tender-platform/app/controllers/kullanicilar"
	Teklifler "blockchain-smart-tender-platform/app/controllers/teklifler"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {

	ihaleGroup := app.Group("/api/v1/ihaleler")
	{
		ihaleGroup.Post("/create", Ihaleler.CreateIhale)
		ihaleGroup.Put("/update", Ihaleler.UpdateIhale)
		ihaleGroup.Get("/query/:ihaleNumarasi", Ihaleler.QueryIhale)
		ihaleGroup.Get("/list", Ihaleler.ListIhaleler)
		ihaleGroup.Put("/close/:ihaleNumarasi", Ihaleler.CloseIhale)
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
		kullaniciGroup.Put("/update", Kullanicilar.UpdateUser)
		kullaniciGroup.Delete("/delete", Kullanicilar.DeleteUser)
		kullaniciGroup.Get("/userdata", Kullanicilar.GetUserDashboardData)
	}

}
