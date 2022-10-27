package router

import (
	"MyGram/controllers"
	"MyGram/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Uri() *fiber.App {
    app := fiber.New()

	// User route
    userRouter := app.Group("/users")
    userRouter.Get("/", controllers.GetUsers)
    userRouter.Get("/:Id", controllers.GetUser)
    userRouter.Post("/register", controllers.CreateUser)
    userRouter.Post("/login", controllers.LoginUser)
    userRouter.Use(middlewares.Authentication())
    userRouter.Put("/:Id", controllers.UpdateUser)
    userRouter.Delete("/:Id", controllers.DeleteUser)

	// Photo route
	photoRouter := app.Group("/photos")
	photoRouter.Use(middlewares.Authentication())
    photoRouter.Get("/:Id", controllers.GetPhoto)
	photoRouter.Get("/", controllers.GetPhotos)
    photoRouter.Post("/", controllers.CreatePhoto)
    photoRouter.Put("/:Id", controllers.UpdatePhoto)
    photoRouter.Delete("/:Id", controllers.DeletePhoto)
	
    // Comment route
	commentRouter := app.Group("/comments")
	commentRouter.Use(middlewares.Authentication())
    commentRouter.Get("/:Id", controllers.GetComment)
	commentRouter.Get("/", controllers.GetComments)
    commentRouter.Post("/", controllers.CreateComment)
    commentRouter.Put("/:Id", controllers.UpdateComment)
    commentRouter.Delete("/:Id", controllers.DeleteComment)
    
    // Social Media route
	socialMediaRouter := app.Group("/socialmedias")
	socialMediaRouter.Use(middlewares.Authentication())
    socialMediaRouter.Get("/:Id", controllers.GetSocialMedia)
	socialMediaRouter.Get("/", controllers.GetSocialMedias)
    socialMediaRouter.Post("/", controllers.CreateSocialMedia)
    socialMediaRouter.Put("/:Id", controllers.UpdateSocialMedia)
    socialMediaRouter.Delete("/:Id", controllers.DeleteSocialMedia)

	return app
}