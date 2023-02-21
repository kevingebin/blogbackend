package routes

import (
	"blogbackend/controller"
	"blogbackend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUp(app *fiber.App) {
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	app.Use(middleware.IsAuthenticate)
	app.Post("api/post", controller.CreatePost)
	app.Get("/api/allpost", controller.AllPost)
	app.Get("api/allpost/:id", controller.DetailPost)
	app.Put("api/update/:id", controller.UpdatePost)
	app.Get("api/uniquepost", controller.UniquePost)
	app.Delete("api/deletepost/:id", controller.DeletePost)
	app.Post("/api/uploadimage", controller.Upload)
	app.Static("/api/upload/", "./upload")
}
