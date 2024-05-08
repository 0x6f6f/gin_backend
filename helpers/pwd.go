package helpers

import "golang.org/x/crypto/bcrypt"

/*密码加密与验证*/

// 生成密码的哈希值
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// 使用string函数将[]byte转换为string
	return string(hash), nil
}

// 验证密码是否正确
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
