package models

import (
	"time"

	"gorm.io/gorm"
)

/*数据库架构设计*/

type RoleID uint

// 定义用户身份的枚举值
const (
	GENERAL_MANAGER      RoleID = iota // 总经理
	SYSTEM_ADMINISTRATOR               //系统管理员
	SALES_REPRESENTATIVE               //销售代表
	SALES_MANAGER                      //销售经理
	SALES_DIRECTOR                     //销售总监
	ACCOUNTANT                         //会计
	FINANCE_SPECIALIST                 //金融专员
	FINANCE_MANAGER                    //金融经理
	DEFAULT                            //默认权限
)

// 枚举值到名称的映射
var RoleNameMap = map[RoleID]string{
	GENERAL_MANAGER:      "总经理",
	SYSTEM_ADMINISTRATOR: "系统管理员",
	SALES_REPRESENTATIVE: "销售代表",
	SALES_MANAGER:        "销售经理",
	SALES_DIRECTOR:       "销售总监",
	ACCOUNTANT:           "会计",
	FINANCE_SPECIALIST:   "金融专员",
	FINANCE_MANAGER:      "金融经理",
	DEFAULT:              "默认权限",
}

// 名称到枚举值的映射
var RoleStrToEnumMap = map[string]RoleID{
	"总经理":   GENERAL_MANAGER,
	"系统管理员": SYSTEM_ADMINISTRATOR,
	"销售代表":  SALES_REPRESENTATIVE,
	"销售经理":  SALES_MANAGER,
	"销售总监":  SALES_DIRECTOR,
	"会计":    ACCOUNTANT,
	"金融专员":  FINANCE_SPECIALIST,
	"金融经理":  FINANCE_MANAGER,
	"默认权限":  DEFAULT,
}

type Gender uint

// 定义性别的枚举值
const (
	FEMALE Gender = iota //女性
	MALE                 //男性
)

var GenderNameMap = map[Gender]string{
	FEMALE: "女",
	MALE:   "男",
}

var GenderStrToEnumMap = map[string]Gender{
	"女": FEMALE,
	"男": MALE,
}

// 公司的员工账户
type User struct {
	gorm.Model
	UserName     string      `gorm:"unique"`   // 用户名（唯一）
	PasswordHash string      `gorm:"not null"` // hash后的密码
	RoleID       RoleID      `gorm:"not null"` // 用户身份
	UserProfile  UserProfile // 关联的 UserProfile 实体
	DepartmentID *uint       // 所属部门ID
	ZoneID       *uint       // 所属战区ID
	WorkLogs     []WorkLog   //工作日志
}

// 用户详细信息
type UserProfile struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	Name    string // 姓名
	Age     uint   // 年龄
	Gender  Gender // 性别
	Address string // 地址
	Phone   string // 电话
}

// 销售部门
type Department struct {
	gorm.Model
	Name    string `gorm:"unique"`
	User    []User `gorm:"foreignKey:DepartmentID"` //包含全部销售人员
	ZoneID  *uint  //所属战区ID
	ManagerID *uint   //部门销售经理
}

// 销售部战区
type Zone struct {
	gorm.Model
	Name        string       `gorm:"unique"`
	Departments []Department `gorm:"foreignKey:ZoneID"` //包含全部销售部门
	DirectorID    *uint         //战区销售总监
}

// 销售人员的工作日志
type WorkLog struct {
	gorm.Model
	UserID     uint      `gorm:"not null"`
	Calls      int       // 电话拨打次数
	ValidCalls int       // 有效电话拨打次数
	Visits     int       // 面谈客户次数
	Contracts  int       // 签订合同次数
	Date       time.Time // 工作日志日期
}

// 贷款客户
type Customer struct {
	gorm.Model
	Name    string `gorm:"not null"` // 姓名
	Phone   string `gorm:"not null"` // 电话
	Age     uint   // 年龄
	Gender  Gender // 性别
	Address string // 地址
	// 贷款意向，起始为10，每天减1，为0时表示不再有贷款意向，需要移入客户公海
	// 如果用户贷款，将其设置为100
	LoanIntent    int        `gorm:"not null"`
	IsInPublicSea bool       `gorm:"not null"`              // 是否在客户公海
	Contracts     []Contract `gorm:"foreignKey:CustomerID"` // 有关的贷款合同
	// 如果全为nil表示当前客户在公海
	SalerID      *uint // 当前接触的销售代表ID
	DepartmentID *uint // 当前所属部门ID
	ZoneID       *uint // 当前所属战区ID
}

type ContractStatus uint

// 定义贷款合同的状态
const (
	NEW       ContractStatus = iota // 新建
	APPROVING                       // 审批中
	APPROVED                        // 已批准
	REJECTED                        // 已拒绝
)

var ContractStatusNameMap = map[ContractStatus]string{
    NEW:       "新建",
	APPROVING: "审批中",
	APPROVED:   "已批准",
	REJECTED:   "已拒绝",
}

var ContractStatusStrToEnumMap = map[string]ContractStatus{
    "新建":       NEW,
	"审批中": APPROVING,
	"已批准":   APPROVED,
	"已拒绝":   REJECTED,
}

// 贷款详情
type Contract struct {
	gorm.Model
	// 贷款信息
	Amount           float64        // 贷款金额
	ServiceFee       float64        // 服务费
	Status           ContractStatus // 贷款状态
	ContractDocument string         // 合同文档
	FinancialProduct string         // 金融产品
	BankDocuments    string         // 银行文件
	BankAmount       float64        // 银行金额(实际金额)
	Image            []byte         // 合同图片
	// 相关人员
	CustomerID   uint // 贷款客户ID
	SalerID      uint // 销售人员ID
	FinanceID    uint // 财务专员ID
	AccountantID uint // 会计ID
	DepartmentID uint // 所属部门ID
	ZoneID       uint // 所属战区ID
}

// 系统日志
type SystemLog struct {
	gorm.Model
	UserID uint   `gorm:"not null"`           // 关联的用户ID
	Action string `gorm:"type:text;not null"` // 日志动作或消息
}
