package config

import "github.com/spf13/viper"

func JWTSecret() string {
	return viper.GetString("JWT_SECRET")
}

func JWTConfiguration() (string, string, uint, uint) {
	jwtSecret := viper.GetString("JWT_SECRET")
	jwtAlgorithm := viper.GetString("JWT_ALGORITHM")
	jwtExpirationTime := viper.GetInt64("JWT_ACCESS_TOKEN_EXPIRE_MINUTES")
	jwtRefreshTime := viper.GetInt64("JWT_REFRESH_TOKEN_EXPIRE_MINUTES")
	return jwtSecret, jwtAlgorithm, uint(jwtExpirationTime), uint(jwtRefreshTime)
}
