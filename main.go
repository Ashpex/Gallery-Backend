package main

import (
	"example.com/gallery/controller"
	"example.com/gallery/database"
	"example.com/gallery/middleware"
	"example.com/gallery/repository"
	"example.com/gallery/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	// db is a global variable that represents the database connection
	db             *gorm.DB                  = database.SetupDB()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	// jwtService is a global variable that represents the jwt service (json web token)
	jwtService service.JWTService = service.NewJWTService()
	// Authentication service and controller
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	// User service and controller
	userService    service.UserService       = service.NewUserService(userRepository)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
)

func main() {
	r := gin.Default()
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	err := r.Run()
	if err != nil {
		return
	}
}
