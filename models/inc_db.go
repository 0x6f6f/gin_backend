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
)

type Gender uint

// 定义性别的枚举值
const (
	FEMALE Gender = iota //女性
	MALE                 //男性
)

// 公司的员工账户
type User struct {
	gorm.Model
	UserName     string      `gorm:"unique"`   // 用户名唯一
	PasswordHash string      `gorm:"not null"` // hash后的密码
	RoleID       RoleID      `gorm:"not null"` // 用户身份
	UserProfile  UserProfile `gorm:"not null"` // 关联的 UserProfile 实体
}

// 用户详细信息
type UserProfile struct {
	gorm.Model
	UserID  uint `gorm:"not null"`
	Name    string
	Age     uint
	Gender  Gender
	Address string
	Phone   string
}

// 销售部门
type SalesDepartment struct {
	gorm.Model
	Name   string  `gorm:"not null"`
	Salers []Saler `gorm:"foreignKey:SalesDepartmentID"`
	ZoneID uint    `gorm:"not null"`
	Zone   Zone    `gorm:"foreignKey:ZoneID"`
}

// 销售部战区
type Zone struct {
	gorm.Model
	Name             string            `gorm:"not null"`
	SalesDepartments []SalesDepartment `gorm:"foreignKey:ZoneID"`
}

// 销售人员
type Saler struct {
	gorm.Model
	RoleID            RoleID    `gorm:"not null"`          // 销售人员的角色(销售代表、销售经理、销售总监)
	UserID            uint      `gorm:"not null"`          // 销售人员的账户ID
	SalesDepartmentID uint      `gorm:"not null"`          // 销售人员所在的销售部门ID
	ZoneID            uint      `gorm:"not null"`          // 销售人员所在的战区ID
	WorkLogs          []WorkLog `gorm:"foreignKey:UserID"` // 关联的 WorkLog 实体
}

// 销售人员的工作日志
type WorkLog struct {
	gorm.Model
	UserID     uint      `gorm:"not null"`
	User       User      `gorm:"foreignKey:UserID"`
	Calls      int       // 电话拨打次数
	ValidCalls int       // 有效电话拨打次数
	Visits     int       // 面谈客户次数
	Contracts  int       // 签订合同次数
	Date       time.Time // 工作日志日期
}

// 财务专员
type FinanceSpecialist struct {
	gorm.Model
	RoleID RoleID `gorm:"not null"` // 财务人员的角色(金融专员、金融经理)
	UserID uint   `gorm:"not null"` // 销售人员的账户ID
}

// 会计
type Accountant struct {
	gorm.Model
	RoleID RoleID `gorm:"not null"` // 会计的角色
	UserID uint   `gorm:"not null"` // 会计的账户ID
}

// 贷款客户
type Customer struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Phone string `gorm:"not null"`
	// 贷款意向，起始为10，每天减1，为0时表示不再有贷款意向，需要移入客户公海
	// 如果用户贷款，将其设置为100
	LoanIntent    int        `gorm:"not null"`
	IsInPublicSea bool       `gorm:"not null"`              // 是否在客户公海
	Contracts     []Contract `gorm:"foreignKey:CustomerID"` // 有关的贷款合同
	SalerID       uint       // 当前接触的销售代表ID
}

type ContractStatus uint

// 定义贷款合同的状态
const (
	NEW       ContractStatus = iota // 新建
	APPROVING                       // 审批中
	APPROVED                        // 已批准
	REJECTED                        // 已拒绝
)

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
	SpecialistID uint // 财务专员ID
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
