package repository

import (
	"fmt"
	"gin-boilerplate/helpers"
	"gin-boilerplate/models"
	"time"

	"gorm.io/gorm"
)

/*用户管理*/
// CreateUser 新建User用户
// 普通用户注册账户
func CreateUser(db *gorm.DB, userName, password string, departmentID uint) (*models.User, error) {
	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := models.User{UserName: userName, PasswordHash: passwordHash, RoleID: models.DEFAULT, DepartmentID: &departmentID}
	err = db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	logAction(db, user.ID, "新建用户")
	return &user, nil
}

// DeleteUser 删除User用户
// 只有系统管理员才能注销账户
func DeleteUser(db *gorm.DB, systemManagerID, userID uint) error {
	if err := db.Delete(&models.User{}, userID).Error; err != nil {
		return err
	}
	logAction(db, systemManagerID, fmt.Sprintf("注销用户: %d", userID))
	return nil
}

// Login User用户登录，验证用户名和密码
func Login(db *gorm.DB, userName, password string) (*models.User, error) {
	var user models.User
	err := db.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		logAction(db, 0, fmt.Sprintf("用户名: %s 错误，登录失败", userName))
		return nil, err
	}
	err = helpers.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		logAction(db, user.ID, "密码错误，登录失败")
		return nil, err
	}
	logAction(db, user.ID, fmt.Sprintf("用户名: %s, 角色%s, 登录成功", userName, models.RoleNameMap[user.RoleID]))
	//返回user实体
	return &user, nil
}

// GetUserByUserName 根据用户名获取用户信息，成功返回User实体
func GetUserByUserName(db *gorm.DB, userName string) (*models.User, error) {
	var user models.User
	err := db.Preload("UserProfile").Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据用户ID获取用户信息，成功返回User实体
func GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	var user models.User
	err := db.Preload("UserProfile").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新User账户信息
// 系统管理员或用户修改用户名或密码（用户修改账户信息传入systemManagerID=0）
func UpdateUserNameOrPassword(db *gorm.DB, systemManagerID, userID uint, userName, password string) (*models.User, error) {
	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return nil, err
	}
	err = db.Model(&models.User{}).Where("id = ?", userID).Updates(models.User{
		UserName:     userName,
		PasswordHash: passwordHash,
	}).Error
	if err != nil {
		return nil, err
	}
	if systemManagerID == 0 {
		logAction(db, userID, "更改账户名或密码")
	} else {
		logAction(db, systemManagerID, fmt.Sprintf("更改用户: %d 账户名或密码", userID))
	}
	return GetUserByID(db, userID)
}

// UpdateUserRole 更新用户角色
// 只有系统管理员才能修改用户角色
func UpdateUserRole(db *gorm.DB, systemManagerID, userID uint, roleID models.RoleID) (*models.User, error) {
	err := db.Model(&models.User{}).Where("id = ?", userID).Update("role_id", roleID).Error
	if err != nil {
		return nil, err
	}
	logAction(db, systemManagerID, fmt.Sprintf("更改用户: %d 角色为: %d", userID, roleID))
	return GetUserByID(db, userID)
}

// UpdateUserProfile 更新用户的个人信息，成功返回nil，否则返回error
// 用户需要更新自己的个人信息
func UpdateUserProfile(db *gorm.DB, userID uint, name string, age uint,
	gender models.Gender, address, phone string) (*models.User, error) {
	// 创建或更新用户详细信息
	err := db.Where(models.UserProfile{UserID: userID}).Assign(models.UserProfile{
		Name:    name,
		Age:     age,
		Gender:  gender,
		Address: address,
		Phone:   phone,
	}).FirstOrCreate(&models.UserProfile{}).Error
	if err != nil {
		return nil, err
	}
	logAction(db, userID, "更新用户个人信息")
	return GetUserByID(db, userID)
}

/*系统管理员*/

// CreateSystemManager 新建系统管理员
// 仅用于系统管理员注册账户
func CreateSystemManager(db *gorm.DB, userName, password string) (*models.User, error) {
	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := models.User{UserName: userName, PasswordHash: passwordHash, RoleID: models.SYSTEM_ADMINISTRATOR}
	err = db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	logAction(db, user.ID, "新建系统管理员")
	return &user, nil
}

// GetUserList 用户列表查询
// 系统管理员可以查看所有用户
func GetUserList(db *gorm.DB, systemManagerID uint) ([]models.User, error) {
	var users []models.User
	err := db.Preload("UserProfile").Find(&users).Error
	if err != nil {
		return nil, err
	}
	logAction(db, systemManagerID, "查看用户列表")
	return users, nil
}

// CreateZone 新建销售战区
// 系统管理员可以进行战区注册
func CreateZone(db *gorm.DB, systemManagerID uint, name string) (*models.Zone, error) {
	err := db.Create(&models.Zone{Name: name}).Error
	if err != nil {
		return nil, err
	}
	logAction(db, systemManagerID, fmt.Sprintf("新建销售战区: %s", name))
	var zone models.Zone
	err = db.Where("name = ?", name).First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

// CreateSalesDepartment 新建销售部门
// 系统管理员可以进行部门注册
func CreateSalesDepartment(db *gorm.DB, systemManagerID uint, name string, zoneID uint) (*models.Department, error) {
	err := db.Create(&models.Department{Name: name, ZoneID: zoneID}).Error
	if err != nil {
		return nil, err
	}
	logAction(db, systemManagerID, fmt.Sprintf("新建销售部门: %s", name))
	var department models.Department
	err = db.Where("name = ?", name).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// CreateFinanceDepartment 新建金融部门
// 系统管理员可以进行部门注册
func CreateFinanceDepartment(db *gorm.DB, systemManagerID uint, name string) (*models.Department, error) {
	err := db.Create(&models.Department{Name: name}).Error
	if err != nil {
		return nil, err
	}
	logAction(db, systemManagerID, fmt.Sprintf("新建金融部门: %s", name))
	var department models.Department
	err = db.Where("name = ?", name).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// AssignDepartmentToZone 分配部门到战区
// 系统管理员可以分配部门到战区
func AssignDepartmentToZone(db *gorm.DB, systemManagerID, departmentID, zoneID uint) error {
	//获取部门
	var department models.Department
	if err := db.Where("id = ?", departmentID).First(&department).Error; err != nil {
		return err
	}
	//获取战区
	var zone models.Zone
	if err := db.Where("id = ?", zoneID).First(&zone).Error; err != nil {
		return err
	}
	//更新部门所属战区
	if err := db.Model(&models.Department{}).Where("id = ?", departmentID).Update("zone_id", zoneID).Error; err != nil {
		return err
	}
	//更新战区内包含的部门
	if err := db.Model(&models.Zone{}).Where("id = ?", zoneID).Association("Departments").Append(&department); err != nil {
		return err
	}
	logAction(db, systemManagerID, fmt.Sprintf("分配部门: %d 到战区: %d", departmentID, zoneID))
	return nil
}

// AssignUserToDepartment 分配用户到部门
// 系统管理员可以分配用户到部门
func AssignUserToDepartment(db *gorm.DB, systemManagerID, userID, departmentID uint) error {
	//获取用户
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	//获取部门
	var department models.Department
	if err := db.Where("id = ?", departmentID).First(&department).Error; err != nil {
		return err
	}
	//更新用户所属部门
	if err := db.Model(&models.User{}).Where("id = ?", userID).Update("department_id", departmentID).Error; err != nil {
		return err
	}
	//更新部门内包含的用户
	if err := db.Model(&models.Department{}).Where("id = ?", departmentID).Association("models.User").Append(&user); err != nil {
		return err
	}
	logAction(db, systemManagerID, fmt.Sprintf("分配用户: %d 到部门: %d", userID, departmentID))
	return nil
}

// AssignUserToZone 分配用户到战区
// 系统管理员可以分配用户到战区
func AssignUserToZone(db *gorm.DB, systemManagerID, userID, zoneID uint) error {
	//获取用户
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	//获取战区
	var zone models.Zone
	if err := db.Where("id = ?", zoneID).First(&zone).Error; err != nil {
		return err
	}
	//更新用户所属战区
	if err := db.Model(&models.User{}).Where("id = ?", userID).Update("zone_id", zoneID).Error; err != nil {
		return err
	}
	logAction(db, systemManagerID, fmt.Sprintf("分配用户: %d 到战区: %d", userID, zoneID))
	return nil
}

// GetSystemLogList 日志查询
// 系统管理员可以查看所有日志
func GetSystemLogList(db *gorm.DB, systemManagerID uint) ([]models.SystemLog, error) {
	var logs []models.SystemLog
	err := db.Find(&logs).Error
	if err != nil {
		return nil, err
	}
	logAction(db, systemManagerID, "查看日志列表")
	return logs, nil
}

/*客户管理*/

// CreateCustomer 销售人员新建客户信息
func CreateCustomer(db *gorm.DB, userID uint, name, phone string) error {
	customer := models.Customer{Name: name, Phone: phone, LoanIntent: 10, IsInPublicSea: false, SalerID: userID}
	err := db.Create(&customer).Error
	if err != nil {
		return err
	}
	logAction(db, userID, fmt.Sprintf("新建客户: %d 信息", customer.ID))
	return nil
}

// UpdateCustomer 销售人员更新客户基本信息
func UpdateCustomer(db *gorm.DB, userID, customerID uint, name, phone string, age uint, gender models.Gender, address string) error {
	err := db.Model(&models.Customer{}).Where("id = ?", customerID).Updates(models.Customer{
		Name:    name,
		Phone:   phone,
		Age:     age,
		Gender:  gender,
		Address: address,
	}).Error
	if err != nil {
		return err
	}
	logAction(db, userID, fmt.Sprintf("更新客户: %d 信息", customerID))
	return nil
}

// GetCustomer 查询客户信息
// 销售人员可以查看自己的客户信息，销售部长可以查看部门内的客户信息，
// 销售总监可以查看战区内的客户信息，总经理可以查看所有客户信息
func GetCustomer(db *gorm.DB, userID, customerID uint) (models.Customer, error) {
	var customer models.Customer
	err := db.Where("id = ?", customerID).First(&customer).Error
	if err != nil {
		return models.Customer{}, err
	}
	logAction(db, userID, fmt.Sprintf("查看客户: %d 信息", customerID))
	return customer, nil
}

// GetPublicSeaCustomerList 查询公海客户列表
// 所有人可以查看公海客户列表
func GetPublicSeaCustomerList(db *gorm.DB, userID uint) ([]models.Customer, error) {
	var customers []models.Customer
	err := db.Where("is_in_public_sea = ?", true).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	logAction(db, userID, "查看公海客户列表")
	return customers, nil
}

// MigrateCustomer 迁移客户
// 总经理可以跨战区迁移，销售总监可以战区内迁移，销售部长可以部门内迁移
func MigrateCustomer(db *gorm.DB, userID, newSalerID, customerID uint) error {
	// 更新客户的销售人员归属
	if err := db.Model(&models.Customer{}).Where("id = ?", customerID).Update("saler_id", newSalerID).Error; err != nil {
		return err
	}
	logAction(db, userID, fmt.Sprintf("迁移了客户：%d 到销售人员：%d", customerID, newSalerID))
	return nil
}

// AutoUpdateCustomerLoanIntent 系统每天自动更新客户贷款意向
func AutoUpdateCustomerLoanIntent(db *gorm.DB) error {
	// 使用事务确保整个操作的一致性
	return db.Transaction(func(tx *gorm.DB) error {
		// 对贷款意向小于等于10且大于0的客户，每天减1
		result := tx.Model(&models.Customer{}).
			Where("loan_intent > 0 AND loan_intent <= 10").
			Update("loan_intent", gorm.Expr("loan_intent - 1"))
		if result.Error != nil {
			return result.Error
		}
		// 记录操作影响的行数
		if result.RowsAffected > 0 {
			logAction(db, 0, fmt.Sprintf("自动更新了 %d 个客户的贷款意向", result.RowsAffected))
		}
		return nil
	})
}

// AutoMigrateCustomerToPublicSea 将客户自动迁移到公海
func AutoMigrateCustomerToPublicSea(db *gorm.DB) error {
	// 使用事务确保更新操作的原子性
	return db.Transaction(func(tx *gorm.DB) error {
		// 将贷款意向为0的客户移入公海，同时更新SalerID为0，表示不再有销售人员负责
		result := tx.Model(&models.Customer{}).Where("loan_intent = 0").Updates(map[string]interface{}{
			"is_in_public_sea": true,
			"saler_id":         0,
		})
		if result.Error != nil {
			return result.Error
		}
		// 记录迁移操作的日志
		if result.RowsAffected > 0 {
			logAction(db, 0, fmt.Sprintf("自动迁移了 %d 个客户到公海", result.RowsAffected))
		}
		return nil
	})
}

/*销售代表记录工作日志*/

// CreateWorkLog 销售人员记录工作日志
func CreateWorkLog(db *gorm.DB, userID uint, calls, validCalls, visits, contracts int, date time.Time) error {
	workLog := models.WorkLog{UserID: userID, Calls: calls, ValidCalls: validCalls, Visits: visits, Contracts: contracts, Date: date}
	err := db.Create(&workLog).Error
	if err != nil {
		return err
	}
	logAction(db, userID, "记录工作日志")
	return nil
}

/*合同管理*/

// SubmitContract 销售人员提交合同
func SubmitContract(db *gorm.DB, salerID, customerID, finanaceID, accountantID uint,
	amount, serviceFee, bankAmount float64,
	financialProduct, contractDocument, bankDocuments string) error {
	// 获取销售人员信息
	var saler models.User
	if err := db.Where("id = ?", salerID).First(&saler).Error; err != nil {
		return err
	}
	contract := models.Contract{
		Amount:           amount,
		ServiceFee:       serviceFee,
		Status:           models.NEW,
		ContractDocument: contractDocument,
		FinancialProduct: financialProduct,
		BankDocuments:    bankDocuments,
		BankAmount:       bankAmount,
		CustomerID:       customerID,
		SalerID:          salerID,
		FinanceID:        finanaceID,
		AccountantID:     accountantID,
		DepartmentID:     *(saler.DepartmentID),
		ZoneID:           *(saler.ZoneID),
	}
	if err := db.Create(&contract).Error; err != nil {
		return err
	}
	logAction(db, salerID, fmt.Sprintf("提交合同: %d", contract.ID))
	return nil
}

// UpdateContractStatus 更新合同状态
// 金融专员/经理可以更新合同状态为审批中，审批通过或审批拒绝
func UpdateContractStatus(db *gorm.DB, userID, contractID uint, status models.ContractStatus) error {
	// 更新合同状态
	if err := db.Model(&models.Contract{}).Where("id = ?", contractID).Update("status", status).Error; err != nil {
		return err
	}
	logAction(db, userID, fmt.Sprintf("更新了合同: %d 状态为: %d", contractID, status))
	return nil
}

// UpdateContractAmount 更新合同金额信息
// 会计可以更新合同金额信息
func UpdateContractAmount(db *gorm.DB, userID, contractID uint, amount, serviceFee, bankAmount float64) error {
	// 更新合同金额信息
	if err := db.Model(&models.Contract{}).Where("id = ?", contractID).Updates(models.Contract{
		Amount:     amount,
		ServiceFee: serviceFee,
		BankAmount: bankAmount,
	}).Error; err != nil {
		return err
	}
	logAction(db, userID, fmt.Sprintf("更新了合同: %d 金额信息", contractID))
	return nil
}

// GetContractListBySalerID 查询销售人员的合同列表
// 销售人员可以查看自己的合同列表
func GetContractListBySalerID(db *gorm.DB, userID, salerID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	if err := db.Where("saler_id = ?", salerID).Find(&contracts).Error; err != nil {
		return nil, err
	}
	logAction(db, userID, fmt.Sprintf("查看了销售人员: %d 合同列表", salerID))
	return contracts, nil
}

// GetContractListByDepartmentID 查询销售部门的合同列表
// 销售经理可以查看部门内的合同列表
func GetContractListByDepartmentID(db *gorm.DB, userID, departmentID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	if err := db.Where("department_id = ?", departmentID).Find(&contracts).Error; err != nil {
		return nil, err
	}
	logAction(db, userID, fmt.Sprintf("查看了部门: %d 合同列表", departmentID))
	return contracts, nil
}

// GetContractListByZoneID 查询销售战区的合同列表
// 销售总监可以查看战区内的合同列表
func GetContractListByZoneID(db *gorm.DB, userID, zoneID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	if err := db.Where("zone_id = ?", zoneID).Find(&contracts).Error; err != nil {
		return nil, err
	}
	logAction(db, userID, fmt.Sprintf("查看了战区: %d 合同列表", zoneID))
	return contracts, nil
}

// GetContractList 查询所有合同列表
// 总经理/金融经理/会计可以查看所有合同列表
func GetContractList(db *gorm.DB, userID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	if err := db.Find(&contracts).Error; err != nil {
		return nil, err
	}
	logAction(db, userID, "查看了所有合同列表")
	return contracts, nil
}

// GetContract 查询合同信息
// 销售人员/金融专员可以查看自己的合同信息，销售经理可以查看部门内的合同信息，
// 销售总监可以查看战区内的合同信息，总经理/金融经理/会计可以查看所有合同信息
func GetContract(db *gorm.DB, userID, contractID uint) (models.Contract, error) {
	// 获取contractID的合同信息
	var contract models.Contract
	if err := db.Where("id = ?", contractID).First(&contract).Error; err != nil {
		return models.Contract{}, err
	}
	logAction(db, userID, fmt.Sprintf("查看了合同: %d 信息", contractID))
	return contract, nil
}

/*业绩与报表*/

// GetSalerPerformance 销售代表业绩查询
func GetSalerPerformance(db *gorm.DB, userID, salerID uint, startDate, endDate time.Time) (float64, error) {
	var totalAmount float64
	// 累加指定销售代表、指定时间范围内的合同金额
	err := db.Model(&models.Contract{}).
		Where("saler_id = ? AND created_at >= ? AND created_at <= ?", salerID, startDate, endDate).
		Select("sum(amount) as total_amount").
		Row().Scan(&totalAmount)
	if err != nil {
		return 0, err
	}
	logAction(db, userID, fmt.Sprintf("查看了销售人员: %d 的业绩", salerID))

	return totalAmount, nil
}

// GetDepartmentPerformance 销售部门业绩查询
func GetDepartmentPerformance(db *gorm.DB, userID, departmentID uint, startDate, endDate time.Time) (float64, error) {
	var totalAmount float64
	// 累加指定销售部门、指定时间范围内的合同金额
	err := db.Model(&models.Contract{}).
		Where("department_id = ? AND created_at >= ? AND created_at <= ?", departmentID, startDate, endDate).
		Select("sum(amount) as total_amount").
		Row().Scan(&totalAmount)
	if err != nil {
		return 0, err
	}
	logAction(db, userID, fmt.Sprintf("查看了部门: %d 的业绩", departmentID))
	return totalAmount, nil
}

// GetZonePerformance 销售战区业绩查询
func GetZonePerformance(db *gorm.DB, userID, zoneID uint, startDate, endDate time.Time) (float64, error) {
	var totalAmount float64
	// 累加指定销售战区、指定时间范围内的合同金额
	err := db.Model(&models.Contract{}).
		Where("zone_id = ? AND created_at >= ? AND created_at <= ?", zoneID, startDate, endDate).
		Select("sum(amount) as total_amount").
		Row().Scan(&totalAmount)
	if err != nil {
		return 0, err
	}
	logAction(db, userID, fmt.Sprintf("查看了战区: %d 的业绩", zoneID))
	return totalAmount, nil
}

// LoanAnalysis 贷款业务分析，返回总贷款额，贷款产品数量，平均贷款额
func LoanAnalysis(db *gorm.DB, userID uint) (float64, int, float64, error) {
	var totalAmount float64
	var count int
	// 查询总贷款额
	err := db.Model(&models.Contract{}).Select("sum(amount) as total_amount, count(*) as count").
		Row().Scan(&totalAmount, &count)
	if err != nil {
		return 0, 0, 0, err
	}
	// 计算平均贷款额
	averageAmount := totalAmount / float64(count)
	logAction(db, userID, "查看了贷款业务分析")
	return totalAmount, count, averageAmount, nil
}

// logAction 记录系统日志
func logAction(db *gorm.DB, userID uint, action string) error {
	// 创建SystemLog实例
	logEntry := models.SystemLog{
		UserID: userID,
		Action: action,
	}
	// 保存到数据库
	if err := db.Create(&logEntry).Error; err != nil {
		return err
	}
	return nil
}
