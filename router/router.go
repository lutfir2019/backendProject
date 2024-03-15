package router

import (
	"github.com/gofiber/fiber/v2"
	"go.mod/handlers"
	"go.mod/helper"
	"go.mod/middleware"
)

func Initalize(router *fiber.App) {

	router.Use(middleware.Security)

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Api Lawang Project v.1.0")
	})

	router.Use(middleware.Json)

	auth := router.Group("/api/auth")
	auth.Post("/login", handlers.Login)
	auth.Delete("/logout", handlers.Logout)

	// users := router.Group("/api/users")
	users := router.Group("/api/users", middleware.Authenticated)
	users.Post("/post-user", handlers.CreateUser)
	users.Post("/post-user/all", handlers.GetUsers)
	users.Post("/post-user/:unm", handlers.GetUserByUnm)
	users.Put("/put-user/:unm", handlers.UpdateUserByUnm)
	users.Delete("/delete-user/:unm", handlers.DeleteByUnm)
	users.Put("/put-user/changepassword", handlers.ChangePassword)

	products := router.Group("/api/products", middleware.Authenticated)
	products.Post("/post-product", handlers.CreateProduct)
	products.Post("/post-product/all", handlers.GetProducts)
	products.Post("/post-product/:pcd", handlers.GetProductByCode)
	products.Put("/put-product/:pcd", handlers.UpdateProductByCode)
	products.Delete("/delete-product/:pcd", handlers.DeleteProduct)

	// shops := router.Group("/api/shops")
	shops := router.Group("/api/shops", middleware.Authenticated)
	shops.Post("/post-shop", handlers.CreateShop)
	shops.Post("/post-shop/all", handlers.GetShops)
	shops.Post("/post-shop/:scd", handlers.GetShopByCode)
	shops.Put("/put-shop/:scd", handlers.UpdateShop)
	shops.Delete("/delete-shop/:scd", handlers.DeleteShop)

	router.Use(func(c *fiber.Ctx) error {
		return helper.ResponseBasic(c, 404, "404: Not Found")
	})

}
