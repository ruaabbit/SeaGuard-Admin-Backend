package main

import (
	"log"

	"seaguard-admin-backend/config"
	"seaguard-admin-backend/handlers"
	"seaguard-admin-backend/repository"
	"seaguard-admin-backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "seaguard-admin-backend/docs"
)

// @title           SeaGuard Admin API
// @version         1.0
// @description     这是海洋卫士志愿者管理系统的后端API文档
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api
func main() {
	// @Summary 主函数，初始化和启动服务器
	// 初始化数据库
	config.InitDatabase()

	// 初始化repository层
	activityRepo := repository.NewActivityRepository()
	volunteerRepo := repository.NewVolunteerRepository()
	registrationRepo := repository.NewRegistrationRepository()

	// 初始化service层
	activityService := service.NewActivityService(activityRepo)
	volunteerService := service.NewVolunteerService(volunteerRepo)
	registrationService := service.NewRegistrationService(registrationRepo, activityRepo)

	// 初始化handlers
	activityHandler := handlers.NewActivityHandler(activityService)
	volunteerHandler := handlers.NewVolunteerHandler(volunteerService)
	registrationHandler := handlers.NewRegistrationHandler(registrationService)

	// 创建gin引擎
	r := gin.Default()

	// CORS配置
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 活动相关路由
	r.GET("/api/activities", activityHandler.ListActivities)
	r.POST("/api/activities", activityHandler.CreateActivity)
	r.PUT("/api/activities/:id", activityHandler.UpdateActivity)
	r.DELETE("/api/activities/:id", activityHandler.DeleteActivity)

	// 志愿者相关路由
	r.GET("/api/volunteers", volunteerHandler.ListVolunteers)
	r.POST("/api/volunteers", volunteerHandler.CreateVolunteer)
	r.PUT("/api/volunteers/:id", volunteerHandler.UpdateVolunteer)
	r.DELETE("/api/volunteers/:id", volunteerHandler.DeleteVolunteer)

	// 报名相关路由
	r.GET("/api/activities/:id/registrations", registrationHandler.ListActivityRegistrations)
	r.PUT("/api/registrations/:id/status", registrationHandler.UpdateRegistrationStatus)

// Swagger API文档路由
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 启动服务器
if err := r.Run(":8080"); err != nil {
	log.Fatal("Failed to start server:", err)
}
}
