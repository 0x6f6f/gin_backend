package helpers

func IsValidUsername(username string) bool {
	if len(username) < 1 || len(username) > 20 {
		return false
	}

	for _, char := range username {
		// 检查字符是否是英文字母、数字或汉字
		if !isAlphanumeric(char) && !isChinese(char) {
			return false
		}
	}
	return true
}

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 16 {
		return false
	}

	for _, char := range password {
		if !isAlphanumeric(char) && !isPunctuation(char) {
			return false
		}
	}

	return true
}

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func isChinese(char rune) bool {
	// 检查字符是否是汉字
	// 根据Unicode范围判断
	return (char >= '\u4e00' && char <= '\u9fff')
}

func isPunctuation(char rune) bool {
	punctuations := []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '-', '=', '[', ']', '{', '}', '|', '\\', ';', ':', '\'', '"', '<', '>', ',', '.', '/', '?'}
	for _, p := range punctuations {
		if char == p {
			return true
		}
	}
	return false
}
