package utils

// ComparePassword
// @Description: 比较密码
// @param requestPassword
// @param dbPassword
// @param salt
// @return bool
func ComparePassword(requestPassword string, dbPassword string, salt string) bool {
	if MD5V([]byte(requestPassword+salt)) != dbPassword {
		return false
	}
	return true
}
