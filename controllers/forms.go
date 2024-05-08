package controllers

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
	Username     string `form:"username"`
	Password     string `form:"password"`
	Role         string `form:"role"`
	DepartmentID uint   `form:"department_id"`
}

type ZoneForm struct {
	Name string `form:"name"`
}

type DepartmentForm struct {
	Name   string `form:"name"`
	ZoneID uint   `form:"zone_id"`
}
