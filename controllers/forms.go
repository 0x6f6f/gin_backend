package controllers

import "time"

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

/*
传参时Role参数(中文)和RoldID(uint枚举类型)的映射关系：

- GENERAL_MANAGER:      "总经理",

- SYSTEM_ADMINISTRATOR: "系统管理员",

- SALES_REPRESENTATIVE: "销售代表",

- SALES_MANAGER:        "销售经理",

- SALES_DIRECTOR:       "销售总监",

- ACCOUNTANT:           "会计",

- FINANCE_SPECIALIST:   "金融专员",

- FINANCE_MANAGER:      "金融经理",

- DEFAULT:              "默认权限",
*/
type RegisterForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Role     string `form:"role"`
}

// 更新用户信息
type UpdateUserNameOrPasswordForm struct {
	SystemManagerID uint `form:"system_manager_id"`
	UserID          uint `form:"user_id"`
	// to update
	Username string `form:"username"`
	Password string `form:"password"`
}

type UpdateUserRoleForm struct {
	SystemManagerID uint `form:"system_manager_id"`
	UserID          uint `form:"user_id"`
	// to update
	Role string `form:"role"`
}

type ListAllUsersFrom struct {
	SystemManagerID uint `form:"system_manager_id"`
}

/*
传参时Gender参数(中文)和GenderID(uint枚举类型)的映射关系：

- MALE:   "男",

- FEMALE: "女",
*/
type UpdateUserProfileForm struct {
	UserID uint `form:"user_id"`
	// to update
	Name    string `form:"name"`
	Age     uint   `form:"age"`
	Gender  string `form:"gender"`
	Address string `form:"address"`
	Phone   string `form:"phone"`
}

type CreateZoneForm struct {
	SystemManagerID uint   `form:"system_manager_id"`
	Name            string `form:"name"`
}

/*
Type 参数可以在以下值中选择：

- “销售部”

- “金融部”
*/
type CreateDepartmentForm struct {
	SystemManagerID uint   `form:"system_manager_id"`
	Name            string `form:"name"`
	Type            string `form:"type"`
	ZoneID          *uint  `form:"zone_id"`
}

type AssignDepartmentToZoneForm struct {
	SystemManagerID uint `form:"system_manager_id"`
	DepartmentID    uint `form:"department_id"`
	ZoneID          uint `form:"zone_id"`
}

type AssignUserToDepartmentForm struct {
	SystemManagerID uint `form:"system_manager_id"`
	UserID          uint `form:"user_id"`
	DepartmentID    uint `form:"department_id"`
}

type AssignUserToZoneForm struct {
	SystemManagerID uint `form:"system_manager_id"`
	UserID          uint `form:"user_id"`
	ZoneID          uint `form:"zone_id"`
}

type AssignDirectorToZoneForm struct {
	SystemManagerID uint `form:"system_manager_id"`
	UserID          uint `form:"user_id"`
	ZoneID          uint `form:"zone_id"`
}

type AssignManagerToDepartmentForm struct {
	SystemManagerID uint `form:"system_manager_id"`
	UserID          uint `form:"user_id"`
	DepartmentID    uint `form:"department_id"`
}

type QuerySystemLogForm struct {
	SystemManagerID uint `form:"system_manager_id"`
}

type CreateCustomerForm struct {
	UserID        uint   `form:"user_id"`
	CustomerName  string `form:"customer_name"`
	CustomerPhone string `form:"customer_phone"`
}

/*
CustomerGender同样可以使用以下值：

- MALE:   "男",

- FEMALE: "女"
*/
type UpdateCustomerForm struct {
	UserID          uint   `form:"user_id"`
	CustomerID      uint   `form:"customer_id"`
	CustomerName    string `form:"customer_name"`
	CustomerPhone   string `form:"customer_phone"`
	CustomerAge     uint   `form:"customer_age"`
	CustomerGender  string `form:"customer_gender"`
	CustomerAddress string `form:"customer_address"`
}

type ListCustomersForm struct {
	UserID uint `form:"user_id"`
}

type MigrateCustomerForm struct {
	UserID     uint `form:"user_id"`
	NewSalerID uint `form:"new_saler_id"`
	CustomerID uint `form:"customer_id"`
}

type GetPublicSeaCustomerListForm struct {
	UserID uint `form:"user_id"`
}

/*
创建工作日志，时间字段的格式是标准的RFC3339格式
（例如：2022-01-01T12:34:56Z）
*/
type CreateWorkLogForm struct {
	UserID     uint      `form:"user_id"`
	Calls      int       `form:"calls"`
	ValidCalls int       `form:"valid_calls"`
	Visits     int       `form:"visits"`
	Contracts  int       `form:"contracts"`
	Date       time.Time `form:"date"`
}

type SubmitContractForm struct {
    UserID     uint      `form:"user_id"`
	CustomerID uint      `form:"customer_id"`
	FinanceID  uint      `form:"finance_id"`
	AccountantID uint      `form:"accountant_id"`
	Amount    float64   `form:"amount"`
	ServiceFee float64   `form:"service_fee"`
	BankAmount  float64   `form:"bank_amount"`
	FinancialProduct string    `form:"financial_product"`
	// mock, not implemented. just upload str.
	ContractDocument string    `form:"contract_document"`
	BankDocuments string    `form:"bank_documents"`
}

type GetContractListForm struct {
	UserID uint `form:"user_id"`
}

type GetContractDetailForm struct {
    UserID uint `form:"user_id"`
	ContractID uint `form:"contract_id"`
}

type GetSalerPerformanceForm struct {
    UserID uint `form:"user_id"`
	SalerID uint `form:"saler_id"`
	StartDate time.Time `form:"start_date"`
	EndDate time.Time `form:"end_date"`
}

type GetDepartmentPerformanceForm struct {
    UserID uint `form:"user_id"`
	DepartmentID uint `form:"department_id"`
	StartDate time.Time `form:"start_date"`
	EndDate time.Time `form:"end_date"`
}

type GetZonePerformanceForm struct {
    UserID uint `form:"user_id"`
	ZoneID uint `form:"zone_id"`
	StartDate time.Time `form:"start_date"`
	EndDate time.Time `form:"end_date"`
}

type GetLoanAnalysisForm struct {
    UserID uint `form:"user_id"`
}