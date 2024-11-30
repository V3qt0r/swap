package main

import (
	"log"
	"os"
	"net/http"
	"swap/middleware"
	"swap/repository"
	"swap/services"

	shandlers "swap/handler"
	sdb "swap/datasources"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	ginEngine := gin.Default()
	ginEngine.MaxMultipartMemory = 8 << 20

	// Auto migrate database and init global userDb variable

	// Load dev env from .env file
	_, isEnvSet := os.LookupEnv("ENV")
	if !isEnvSet {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("Error loaing .env file")
		}
	}

	// initialize data sources
	swapDB, err := sdb.InitDS()
	if err != nil {
		log.Printf("Error on application startup: %v\n", err)
	}

	userRepository := repository.NewUserRepository(swapDB.DB)
	itemRepository := repository.NewItemRepository(swapDB.DB)

	userService := services.NewUserService(userRepository)
	itemService := services.NewItemService(itemRepository)

	userHandler := shandlers.NewUserHandler(userService)
	itemHandler := shandlers.NewItemHandler(itemService)

	jwtMiddleware, err := middleware.Middleware(userService)

	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}

	errInit := jwtMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error: " + errInit.Error())
	}



	ginEngine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to swap",
		})
	})

	ginEngine.GET("/cors", func(c *gin.Context) {
		c.File("cors.html")
	})

	ginEngine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Page not found",
		})
	})

	ginEngine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})


	// Create an account group
	userGroup := ginEngine.Group("/api/user")

	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", jwtMiddleware.LoginHandler)
	userGroup.GET("/refresh", jwtMiddleware.RefreshHandler)
	userGroup.POST("/one-time-login", jwtMiddleware.LoginHandler)
	userGroup.POST("/send-onetime-password", userHandler.SendOneTimePassword)
	userGroup.POST("/logout", jwtMiddleware.LogoutHandler)
	userGroup.POST("/upload", userHandler.UploadFile)

	userAuthRoutes := ginEngine.Group("/api/users").Use(jwtMiddleware.MiddlewareFunc())
	userAuthRoutes.GET("/:id", userHandler.GetUserById)
	userAuthRoutes.GET("/self", userHandler.GetLoggedInUser)
	userAuthRoutes.PUT("/self", userHandler.UpdateUser)
	userAuthRoutes.GET("/refresh_token", jwtMiddleware.RefreshHandler)
	userAuthRoutes.PUT("/update-password", userHandler.UpdatePassword)
	userAuthRoutes.PUT("/confirm-password", userHandler.ConfirmUserPassword)
	userAuthRoutes.POST("/search", userHandler.Search)
	userAuthRoutes.POST("/totp/enroll", userHandler.EnrollTOTP)
	userAuthRoutes.POST("/totp/verify", userHandler.VerifyTOTP)
	userAuthRoutes.POST("/totp/disable", userHandler.DisableTOTP)
	userAuthRoutes.POST("/totp/enable", userHandler.EnableTOTP)
	userAuthRoutes.GET("/emailOruserName", userHandler.FindUserByEmailOrUsername)
	userAuthRoutes.GET("/phoneNumber", userHandler.FindUserByPhoneNumber)


	itemGroup := ginEngine.Group("/api/items").Use(jwtMiddleware.MiddlewareFunc())
	itemGroup.POST("/register", itemHandler.RegisterItem)
	itemGroup.PUT("/buy", itemHandler.BuyItem)
	itemGroup.PUT("/swap", itemHandler.SwapItem)

	itemGroup.GET("/:id", itemHandler.GetItemById)
	itemGroup.PUT("/:id", itemHandler.UpdateItem)
	itemGroup.DELETE("/:id", itemHandler.DeleteItem)
	itemGroup.GET("/self", itemHandler.GetItemsByOwnerId)
	itemGroup.GET("/search", itemHandler.GetItemsByCategory)
	itemGroup.GET("/search-unsold", itemHandler.GetUnsoldItemsByCategory)
	itemGroup.POST("/upload/:id", itemHandler.UploadFile)

	
	ginEngine.Run(":" + os.Getenv("PORT"))
}
