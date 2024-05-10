package controllers

import (
	"gin-boilerplate/helpers"
	"gin-boilerplate/infra/database"
	"gin-boilerplate/models"
	"gin-boilerplate/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户注册控制器，注册完毕后返回jwt令牌和用户对象
func UserRegister(ctx *gin.Context) {
	var registerForm RegisterForm
	if err := ctx.ShouldBind(&registerForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid register form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// 创建新用户
	var user *models.User
	var err error

	// 判断用户名密码是否合规
	if !helpers.IsValidUsername(registerForm.Username) {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid username, only 1-20 numbers/alphabets/chinese characters allowed",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if !helpers.IsValidPassword(registerForm.Password) {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid password, 8-16 characters, only numbers and alphabets allowed",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if registerForm.Role == models.RoleNameMap[models.SYSTEM_ADMINISTRATOR] {
		// 如果是系统管理员，则创建系统管理员
		user, err = repository.CreateSystemManager(database.DB,
			registerForm.Username,
			registerForm.Password,
		)
	} else {
		// 否则创建普通用户
		user, err = repository.CreateUser(database.DB,
			registerForm.Username,
			registerForm.Password,
		)
	}

	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// 注册成功
	access_token, refresh_token, err := helpers.GenerateToken(*user)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate jwt token: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{
		Code:    http.StatusOK,
		Message: "Register successful",
		Data: map[string]interface{}{
			"user":          user,
			"access_token":  access_token,
			"refresh_token": refresh_token,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// 登录控制器，登录成功后返回jwt令牌和用户对象
func UserLogin(ctx *gin.Context) {
	// 获取用户名和密码
	var loginForm LoginForm
	if err := ctx.ShouldBind(&loginForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid login form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// 验证用户名和密码
	user, err := repository.Login(database.DB, loginForm.Username, loginForm.Password)
	if err != nil {
		response := Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
		}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	// 登录成功
	access_token, refresh_token, err := helpers.GenerateToken(*user)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate jwt token: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data: map[string]interface{}{
			"user":          user,
			"access_token":  access_token,
			"refresh_token": refresh_token,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// 更新用户信息
func UserUpdateProfile(ctx *gin.Context) {
	// 获取用户ID和更新信息
	var updateForm UpdateUserProfileForm
	if err := ctx.ShouldBind(&updateForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid update form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	// 更新用户信息
	user, err := repository.UpdateUserProfile(
		database.DB,
		updateForm.UserID,
		updateForm.Name,
		updateForm.Age,
		models.GenderStrToEnumMap[updateForm.Gender],
		updateForm.Address,
		updateForm.Phone,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	// 更新成功
	response := Response{
		Code:    http.StatusOK,
		Message: "Update successful",
		Data:    user,
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorUpdateUserNameOrPassword(ctx *gin.Context) {
	// 获取用户ID和更新信息
	var updateForm UpdateUserNameOrPasswordForm
	if err := ctx.ShouldBind(&updateForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid update form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// 更新用户信息
	user, err := repository.UpdateUserNameOrPassword(
		database.DB,
		updateForm.SystemManagerID,
		updateForm.UserID,
		updateForm.Username,
		updateForm.Password,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// 更新成功
	response := Response{
		Code:    http.StatusOK,
		Message: "Update successful",
		Data:    user,
	}
	ctx.JSON(http.StatusOK, response)
}

// 管理员更新其他信息
func AdministratorUpdateUserRole(ctx *gin.Context) {
	// 获取用户ID和更新信息
	var updateForm UpdateUserRoleForm
	if err := ctx.ShouldBind(&updateForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid update form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	// 更新用户信息
	user, err := repository.UpdateUserRole(
		database.DB,
		updateForm.SystemManagerID,
		updateForm.UserID,
		models.RoleStrToEnumMap[updateForm.Role],
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// 更新成功
	response := Response{
		Code:    http.StatusOK,
		Message: "Update successful",
		Data:    user,
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorListAllUsers(ctx *gin.Context) {
	var listForm ListAllUsersFrom
	if err := ctx.ShouldBind(&listForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid list form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	users, err := repository.GetUserList(database.DB, listForm.SystemManagerID)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list users: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{
		Code:    http.StatusOK,
		Message: "List successful",
		Data:    users,
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorCreateZone(ctx *gin.Context) {
	var createForm CreateZoneForm
	if err := ctx.ShouldBind(&createForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid create form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	zone, err := repository.CreateZone(
		database.DB,
		createForm.SystemManagerID,
		createForm.Name,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create zone: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{
		Code:    http.StatusOK,
		Message: "Create successful",
		Data:    zone,
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorCreateDepartment(ctx *gin.Context) {
	var createForm CreateDepartmentForm
	if err := ctx.ShouldBind(&createForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid create form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var department *models.Department
	var err error

	if createForm.Type == "销售部" {
		department, err = repository.CreateSalesDepartment(
			database.DB,
			createForm.SystemManagerID,
			createForm.Name,
			createForm.ZoneID,
		)
	} else if createForm.Type == "金融部" {
		department, err = repository.CreateFinanceDepartment(
			database.DB,
			createForm.SystemManagerID,
			createForm.Name,
		)
	} else {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid department type",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create department: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Create successful",
		Data:    department,
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorAssignDepartmentToZone(ctx *gin.Context) {
	var assignForm AssignDepartmentToZoneForm
	if err := ctx.ShouldBind(&assignForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid assign form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := repository.AssignDepartmentToZone(
		database.DB,
		assignForm.SystemManagerID,
		assignForm.DepartmentID,
		assignForm.ZoneID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to assign department to zone: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Assign successful",
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorAssignUserToDepartment(ctx *gin.Context) {
	var assignForm AssignUserToDepartmentForm
	if err := ctx.ShouldBind(&assignForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid assign form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := repository.AssignUserToDepartment(
		database.DB,
		assignForm.SystemManagerID,
		assignForm.UserID,
		assignForm.DepartmentID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to assign user to department: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Assign successful",
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorAssignUserToZone(ctx *gin.Context) {
	var assignForm AssignUserToZoneForm
	if err := ctx.ShouldBind(&assignForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid assign form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := repository.AssignUserToZone(
		database.DB,
		assignForm.SystemManagerID,
		assignForm.UserID,
		assignForm.ZoneID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to assign user to zone: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Assign successful",
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorAssignDirectorToZone(ctx *gin.Context) {
	var assignForm AssignDirectorToZoneForm
	if err := ctx.ShouldBind(&assignForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid assign form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := repository.AssignDirectorToZone(
		database.DB,
		assignForm.SystemManagerID,
		assignForm.UserID,
		assignForm.ZoneID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to assign director to zone: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Assign successful",
	}
	ctx.JSON(http.StatusOK, response)
}

func AdministratorAssignManagerToDepartment(ctx *gin.Context) {
	var assignForm AssignManagerToDepartmentForm
	if err := ctx.ShouldBind(&assignForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid assign form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := repository.AssignManagerToDepartment(
		database.DB,
		assignForm.SystemManagerID,
		assignForm.UserID,
		assignForm.DepartmentID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to assign manager to department: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Assign successful",
	}
	ctx.JSON(http.StatusOK, response)
}

// 管理员系统日志查询
func AdministratorQuerySystemLog(ctx *gin.Context) {
	var queryForm QuerySystemLogForm
	if err := ctx.ShouldBind(&queryForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid query form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	systemLogs, err := repository.GetSystemLogList(database.DB, queryForm.SystemManagerID)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to query system log: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Query successful",
		Data:    systemLogs,
	}
	ctx.JSON(http.StatusOK, response)
}

// 销售部api控制器
func SaleCreateCustomer(ctx *gin.Context) {
	var createForm CreateCustomerForm
	if err := ctx.ShouldBind(&createForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid create form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	customer, err := repository.CreateCustomer(
		database.DB,
		createForm.UserID,
		createForm.CustomerName,
		createForm.CustomerPhone,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create customer: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Create successful",
		Data:    customer,
	}
	ctx.JSON(http.StatusOK, response)
}

func SaleUpdateCustomer(ctx *gin.Context) {
	var updateForm UpdateCustomerForm
	if err := ctx.ShouldBind(&updateForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid update form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	updated_customer, err := repository.UpdateCustomer(
		database.DB,
		updateForm.UserID,
		updateForm.CustomerID,
		updateForm.CustomerName,
		updateForm.CustomerPhone,
		updateForm.CustomerAge,
		models.GenderStrToEnumMap[updateForm.CustomerGender],
		updateForm.CustomerAddress,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update customer: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Update successful",
		Data:    updated_customer,
	}
	ctx.JSON(http.StatusOK, response)
}

// todo: repo对应的功能还没写
func SaleListCustomers(ctx *gin.Context) {
	var listForm ListCustomersForm
	if err := ctx.ShouldBind(&listForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid list form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	customers, err := repository.ListCustomer(
		database.DB,
		listForm.UserID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list customers: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "List successful",
		Data:    customers,
	}
	ctx.JSON(http.StatusOK, response)

}

// todo: repo对应的功能还没写
func SaleMigrateCustomer(ctx *gin.Context) {
	var migrateForm MigrateCustomerForm
	if err := ctx.ShouldBind(&migrateForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid migrate form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// 根据user身份和migrateForm中的customerID进行迁移操作
	migrated_customer, err := repository.MigrateCustomer(
		database.DB,
		migrateForm.UserID,
		migrateForm.NewSalerID,
		migrateForm.CustomerID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to migrate customer: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Migrate successful",
		Data:    migrated_customer,
	}
	ctx.JSON(http.StatusOK, response)
}

func SaleGetPublicSeaCustomerList(ctx *gin.Context) {
	var getForm GetPublicSeaCustomerListForm
	if err := ctx.ShouldBind(&getForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid get form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	customers, err := repository.GetPublicSeaCustomerList(
		database.DB,
		getForm.UserID,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get public sea customer list: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Get public sea customer list successful",
		Data:    customers,
	}
	ctx.JSON(http.StatusOK, response)
}

// 该控制器用于记录一日的工作情况
func SaleCreateWorkLog(ctx *gin.Context) {
	var createForm CreateWorkLogForm
	if err := ctx.ShouldBind(&createForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid create form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	workLog, err := repository.CreateWorkLog(
		database.DB,
		createForm.UserID,
		createForm.Calls,
		createForm.ValidCalls,
		createForm.Visits,
		createForm.Contracts,
		createForm.Date,
	)
	if err != nil {
		response := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create work log: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Message: "Create work log successful",
		Data:    workLog,
	}
	ctx.JSON(http.StatusOK, response)
}

func SaleSubmitContract(ctx *gin.Context) {
    var submitForm SubmitContractForm
    if err := ctx.ShouldBind(&submitForm); err != nil {
        response := Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid submit form",
        }
        ctx.JSON(http.StatusBadRequest, response)
        return
    }

    contract, err := repository.SubmitContract(
        database.DB,
        submitForm.UserID,
        submitForm.CustomerID,
        submitForm.FinanceID,
		submitForm.AccountantID,
		submitForm.Amount,
		submitForm.ServiceFee,
		submitForm.BankAmount,
		submitForm.FinancialProduct,
		submitForm.ContractDocument,
		submitForm.BankDocuments,
    )
	if err != nil {
	    response := Response{
	        Code:    http.StatusInternalServerError,
	        Message: "Failed to submit contract: " + err.Error(),
	    }
	    ctx.JSON(http.StatusInternalServerError, response)
	    return
	}

	response := Response{
	    Code:    http.StatusOK,
	    Message: "Submit contract successful",
	    Data:    contract,
	}
	ctx.JSON(http.StatusOK, response)
}

func GetContractList(ctx *gin.Context) {
	var getForm GetContractListForm
	if err := ctx.ShouldBind(&getForm); err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid get form",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	contracts, err := repository.GetContractListByUser(
		database.DB,
		getForm.UserID,
	)
	if err != nil {
	    response := Response{
	        Code:    http.StatusInternalServerError,
	        Message: "Failed to get contract list: " + err.Error(),
	    }
	    ctx.JSON(http.StatusInternalServerError, response)
	    return
	}

	response := Response{
	    Code:    http.StatusOK,
	    Message: "Get contract list successful",
	    Data:    contracts,
	}
	ctx.JSON(http.StatusOK, response)
}

func GetContractDetail(ctx *gin.Context) {
    var getForm GetContractDetailForm
    if err := ctx.ShouldBind(&getForm); err != nil {
        response := Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid get form",
        }
        ctx.JSON(http.StatusBadRequest, response)
        return
    }

    contract, err := repository.GetContract(
        database.DB,
        getForm.UserID,
        getForm.ContractID,
    )
    if err != nil {
        response := Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get contract detail: " + err.Error(),
        }
        ctx.JSON(http.StatusInternalServerError, response)
        return
    }

	response := Response{
	    Code:    http.StatusOK,
	    Message: "Get contract detail successful",
	    Data:    contract,
	}
	ctx.JSON(http.StatusOK, response)
}

// GetSalerPerformance 获取销售人员的业绩
func GetSalerPerformance(ctx *gin.Context) {
    var getForm GetSalerPerformanceForm
    if err := ctx.ShouldBind(&getForm); err != nil {
        response := Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid get form",
        }
        ctx.JSON(http.StatusBadRequest, response)
        return
    }

    performance, err := repository.GetSalerPerformance(
        database.DB,
        getForm.UserID,
		getForm.SalerID,
        getForm.StartDate,
        getForm.EndDate,
    )
	if err != nil {
	    response := Response{
	        Code:    http.StatusInternalServerError,
	        Message: "Failed to get saler performance: " + err.Error(),
	    }
	    ctx.JSON(http.StatusInternalServerError, response)
	    return
	}

	response := Response{
	    Code:    http.StatusOK,
	    Message: "Get saler performance successful",
	    Data:    performance,
	}
	ctx.JSON(http.StatusOK, response)
}

// GetDepartmentPerformance 获取部门业绩
func GetDepartmentPerformance(ctx *gin.Context) {
    var getForm GetDepartmentPerformanceForm
    if err := ctx.ShouldBind(&getForm); err != nil {
        response := Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid get form",
        }
		ctx.JSON(http.StatusBadRequest, response)
		return
    }
    performance, err := repository.GetDepartmentPerformance(
	    database.DB,
		getForm.UserID,
		getForm.DepartmentID,
		getForm.StartDate,
		getForm.EndDate,

	)
	if err != nil {
		response := Response{
		    Code:    http.StatusInternalServerError,
		    Message: "Failed to get department performance: " + err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := Response{
	    Code:    http.StatusOK,
	    Message: "Get department performance successful",
	    Data:    performance,
	}
	ctx.JSON(http.StatusOK, response)
}

// GetZonePerformance 获取战区业绩
func GetZonePerformance(ctx *gin.Context) {
    var getForm GetZonePerformanceForm
    if err := ctx.ShouldBind(&getForm); err != nil {
        response := Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid get form",
        }
		ctx.JSON(http.StatusBadRequest, response)
		return

    }

    performance, err := repository.GetZonePerformance(
	    database.DB,
		getForm.UserID,
		getForm.ZoneID,
		getForm.StartDate,
		getForm.EndDate,
	)
	if err != nil {
	    response := Response{
	        Code:    http.StatusInternalServerError,
	        Message: "Failed to get zone performance: " + err.Error(),
	    }
	    ctx.JSON(http.StatusInternalServerError, response)
	    return
	}

	response := Response{
	    Code:    http.StatusOK,
	    Message: "Get zone performance successful",
	    Data:    performance,
	}
	ctx.JSON(http.StatusOK, response)
}

func LoanAnalysis(ctx *gin.Context) {
    var getForm GetLoanAnalysisForm
    if err := ctx.ShouldBind(&getForm); err != nil {
        response := Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid get form",
        }
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	totalAmount, count, averageAmount, err := repository.LoanAnalysis(
	    database.DB,
		getForm.UserID,
	)
	if err != nil {
	    response := Response{
	        Code:    http.StatusInternalServerError,
	        Message: "Failed to get loan analysis: " + err.Error(),
	    }
	    ctx.JSON(http.StatusInternalServerError, response)
	    return
	}

	response := Response{
	    Code:    http.StatusOK,
	    Message: "Get loan analysis successful",
	    Data: gin.H{
	        "total_amount": totalAmount,
	        "count":         count,
	        "average_amount": averageAmount,
	    },
	}
	ctx.JSON(http.StatusOK, response)
}