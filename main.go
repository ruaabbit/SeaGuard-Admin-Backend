package main

import (
	"log"

	"seaguard-admin-backend/config"
	"seaguard-admin-backend/handlers"
	"seaguard-admin-backend/middleware"
	"seaguard-admin-backend/repository"
	"seaguard-admin-backend/service"

	_ "seaguard-admin-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	userRepo := repository.NewUserRepository(config.DB)
	activityRepo := repository.NewActivityRepository()
	volunteerRepo := repository.NewVolunteerRepository()
	registrationRepo := repository.NewRegistrationRepository()

	// 初始化service层
	userService := service.NewUserService(userRepo)
	activityService := service.NewActivityService(activityRepo)
	volunteerService := service.NewVolunteerService(volunteerRepo)
	registrationService := service.NewRegistrationService(registrationRepo, activityRepo)

	// 初始化handlers
	userHandler := handlers.NewUserHandler(userService)
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

	// 认证相关路由（无需认证）
	r.POST("/api/auth/register", userHandler.Register)
	r.POST("/api/auth/login", userHandler.Login)

	// 用户相关路由（需要认证）
	auth := r.Group("/api", middleware.AuthMiddleware(userService))
	{
		// 用户管理（仅管理员）
		admin := auth.Group("", middleware.AdminRequired())
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
		}

		// 活动管理
		auth.GET("/activities", activityHandler.ListAvailableActivities)       // 所有认证用户可查看活动
		admin.GET("/admin/activities", activityHandler.ListActivitiesForAdmin) // 管理员专用查看
		admin.POST("/activities", activityHandler.CreateActivity)
		admin.PUT("/activities/:id", activityHandler.UpdateActivity)
		admin.DELETE("/activities/:id", activityHandler.DeleteActivity)

		// 志愿者相关路由
		// 管理员权限
		admin.GET("/volunteers", volunteerHandler.ListVolunteers)
		admin.POST("/volunteers", volunteerHandler.CreateVolunteer)
		admin.PUT("/volunteers/:id", volunteerHandler.UpdateVolunteer)
		admin.DELETE("/volunteers/:id", volunteerHandler.DeleteVolunteer)

		// 志愿者个人路由（需要志愿者权限）
		volunteer := auth.Group("", middleware.VolunteerRequired())
		{
			// 志愿者个人信息
			volunteer.GET("/volunteer/my-info", volunteerHandler.GetMyInfo)
			volunteer.PUT("/volunteer/my-info", volunteerHandler.UpdateMyInfo)

			// 活动报名相关
			volunteer.GET("/activities/:id/registrations", registrationHandler.ListActivityRegistrations)
			volunteer.PUT("/registrations/:id/status", registrationHandler.UpdateRegistrationStatus)
			volunteer.POST("/activities/:id/register", registrationHandler.Register)             // 活动报名
			volunteer.GET("/activities/:id/registration", registrationHandler.GetMyRegistration) // 查询个人报名状态
		}

		// 通用功能
		auth.PUT("/auth/password", userHandler.ChangePassword)
	}

	// Swagger API文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
