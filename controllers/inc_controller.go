package controllers

import (
	"gin-boilerplate/helpers"
	"gin-boilerplate/infra/database"
	"gin-boilerplate/models"
	"gin-boilerplate/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		if registerForm.DepartmentID == 0 {
			response := Response{
				Code:    http.StatusBadRequest,
				Message: "DepartmentID is required",
			}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		user, err = repository.CreateUser(database.DB,
			registerForm.Username,
			registerForm.Password,
			registerForm.DepartmentID,
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
	response := Response{
		Code:    http.StatusOK,
		Message: "Register successful",
		Data:    user,
	}
	ctx.JSON(http.StatusOK, response)

}

// 登录
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
	response := Response{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data:    user,
	}
	ctx.JSON(http.StatusOK, response)
}

// 部门CRUD
func DepartmentCreate(ctx *gin.Context) {

}
