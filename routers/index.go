package routers

import (
	"gin-boilerplate/controllers"
	"gin-boilerplate/routers/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	//Add All route
	api_version := "/api/v1"
	route.GET(api_version+"/register", controllers.UserRegister)
	route.GET(api_version+"/login", controllers.UserLogin)
	route.GET("/updateUserProfile", controllers.UserUpdateProfile)

	adminGroup := route.Group(api_version+"/admin", middleware.UserRoleAuthMiddleware([]string{"系统管理员"}))
	{
		// user ops
		adminGroup.GET("/updateUserBasicInfo", controllers.AdministratorUpdateUserNameOrPassword)
		adminGroup.GET("/updateUserRole", controllers.AdministratorUpdateUserRole)
		adminGroup.GET("/listAllUsers", controllers.AdministratorListAllUsers)
		// zone & department ops
		adminGroup.GET("/createZone", controllers.AdministratorCreateZone)
		adminGroup.GET("/createDepartment", controllers.AdministratorCreateDepartment)

		// todo: not tested
		// admin assigning ops
		adminGroup.GET("/assignDepartmentToZone", controllers.AdministratorAssignDepartmentToZone)
		adminGroup.GET("/assignUserToDepartment", controllers.AdministratorAssignUserToDepartment)
		adminGroup.GET("/assignUserToZone", controllers.AdministratorAssignUserToZone)
		adminGroup.GET("/assignDirectorToZone", controllers.AdministratorAssignDirectorToZone)
		adminGroup.GET("/assignManagerToDepartment", controllers.AdministratorAssignManagerToDepartment)

		// todo: not tested
		// read system log
		adminGroup.GET("/readSystemLog", controllers.AdministratorQuerySystemLog)
	}

	saleGroup := route.Group(api_version+"/sale", middleware.UserRoleAuthMiddleware([]string{"销售代表", "销售经理", "销售总监"}))
	{
		saleGroup.GET("/createCustomer", controllers.SaleCreateCustomer)
		saleGroup.GET("/updateCustomer", controllers.SaleUpdateCustomer)
		saleGroup.GET("/listCustomers", controllers.SaleListCustomers)
	}
}
