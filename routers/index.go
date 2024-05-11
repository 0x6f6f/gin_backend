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
	route.GET(api_version+"/updateUserProfile", controllers.UserUpdateProfile)

	// todo: not tested
	// stats
	route.GET(api_version+"/getSalerPerformance", controllers.GetSalerPerformance)
	route.GET(api_version+"/getDepartmentPerformance", controllers.GetDepartmentPerformance)
	route.GET(api_version+"/getZonePerformance", controllers.GetZonePerformance)
	route.GET(api_version+"/getLoanAnalysis", controllers.LoanAnalysis)

	// todo: not tested
	// get methods for department and zone
	route.GET(api_version+"/getDepartments", controllers.GetDepartments)
	route.GET(api_version+"/getZones", controllers.GetZones)
	route.GET(api_version+"/getDepartmentByID", controllers.GetDepartmentByID)
	route.GET(api_version+"/getZoneByID", controllers.GetZoneByID)

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
		// todo: not tested
		// 管理客户
		saleGroup.GET("/createCustomer", controllers.SaleCreateCustomer)
		saleGroup.GET("/updateCustomer", controllers.SaleUpdateCustomer)
		saleGroup.GET("/listCustomers", controllers.SaleListCustomers)
		saleGroup.GET("/migrateCustomer", controllers.SaleMigrateCustomer)
		saleGroup.GET("/getPublicSeaCustomerList", controllers.SaleGetPublicSeaCustomerList)
		// todo: not tested
		// 管理工作日志
		saleGroup.GET("/createWorkLog", controllers.SaleCreateWorkLog)
		// 提交合同
		saleGroup.GET("/submitContract", controllers.SaleSubmitContract)
	}

	finanaceGroup := route.Group(api_version+"/finance", middleware.UserRoleAuthMiddleware([]string{"金融专员", "金融经理"}))
	{
		finanaceGroup.GET("/updateContractStatus", controllers.FinanaceUpdateContractStatus)
		finanaceGroup.GET("/updateContractAmount", controllers.FinanaceUpdateContractAmount)
	}

	contractAccessGroup := route.Group(api_version+"/contract", middleware.UserRoleAuthMiddleware([]string{"销售代表", "销售经理", "销售总监", "总经理", "金融经理", "会计"}))
	{
		// todo: not tested
		// 获取合同列表
		contractAccessGroup.GET("/getContractList", controllers.GetContractList)
		// 获取合同详情
		contractAccessGroup.GET("/getContractDetail", controllers.GetContractDetail)
	}
}
