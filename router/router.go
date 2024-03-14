package router

import (
	"github.com/gofiber/fiber/v2"
	"go.mod/handlers"
	"go.mod/middleware"
)

func Initalize(router *fiber.App) {

	router.Use(middleware.Security)

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Api Lawang Project v.1.0")
	})

	router.Use(middleware.Json)

	users := router.Group("/users")
	users.Post("/post-user", handlers.CreateUser)
	users.Post("/post-user/login", handlers.Login)
	users.Delete("/delete-user/logout", handlers.Logout)
	users.Delete("/delete-user/all", middleware.Authenticated, handlers.DeleteByUnm)
	users.Post("/post-user/:unm", middleware.Authenticated, handlers.GetUserByUnm)
	users.Delete("/delete-user/:unm", handlers.DeleteByUnm)
	users.Put("/put-user/changepassword", middleware.Authenticated, handlers.ChangePassword)

	products := router.Group("/products", middleware.Authenticated)
	products.Post("/post-product", handlers.CreateProduct)
	products.Post("/post-product/all", handlers.GetProducts)
	products.Post("/post-product/:pcd", handlers.GetProductByCode)
	products.Put("/put-product/:pid", handlers.UpdateProduct)
	products.Delete("/delete-product/:pcd", handlers.DeleteProduct)

	shops := router.Group("/shops")
	// shops := router.Group("/shops", middleware.Authenticated)
	shops.Post("/post-shop", handlers.CreateShop)
	shops.Post("/post-shop/all", handlers.GetShops)
	shops.Post("/post-shop/:sid", handlers.GetShopByCode)
	shops.Put("/put-shop/:sid", handlers.UpdateShop)
	shops.Delete("/delete-shop/:sid", handlers.DeleteShop)

	router.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"message": "404: Not Found",
		})
	})

}
