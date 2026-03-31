package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware (simplified for MVP)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiGroup := r.Group("/api")
	{
		users := apiGroup.Group("/users")
		{
			users.POST("/register", Register)
			users.POST("/login", Login)
			users.GET("/me", AuthMiddleware(), GetCurrentUser)
		}

		rooms := apiGroup.Group("/rooms")
		{
			rooms.GET("", AuthMiddleware(), GetRooms)
			rooms.GET("/:id/bookings", AuthMiddleware(), GetRoomBookings) // 新增：查看房间预约情况
			// Admin only
			rooms.POST("", AuthMiddleware(), AdminRequired(), CreateRoom)
			rooms.PUT("/:id", AuthMiddleware(), AdminRequired(), UpdateRoom)
			rooms.DELETE("/:id", AuthMiddleware(), AdminRequired(), DeleteRoom)
		}

		bookings := apiGroup.Group("/bookings")
		bookings.Use(AuthMiddleware())
		{
			bookings.POST("", CreateBooking)
			bookings.GET("/my", GetMyBookings)
			bookings.PUT("/:id/cancel", CancelBooking)
			bookings.PUT("/:id/checkin", CheckInBooking) // 新增签到接口
			
			// Admin only
			bookings.GET("", AdminRequired(), GetAllBookings)
			bookings.PUT("/:id/approve", AdminRequired(), ApproveBooking)
		}
	}

	return r
}
