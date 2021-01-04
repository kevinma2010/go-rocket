package utils

import (
	"strconv"
	"strings"
)

const (
	// VersionEquality 版本相等
	VersionEquality = 0
	// VersionGreater 版本大于
	VersionGreater = 1
	// VersionLess 版本小于
	VersionLess = -1
)

// CompareVersion 对比版本号
func CompareVersion(curr, base string) int {
	currArr := strings.Split(curr, ".")
	baseArr := strings.Split(base, ".")

	for i := len(currArr); i < 4; i++ {
		currArr = append(currArr, "0")
	}
	for i := len(baseArr); i < 4; i++ {
		baseArr = append(baseArr, "0")
	}
	for i := 0; i < 4; i++ {
		version1, _ := strconv.Atoi(currArr[i])
		version2, _ := strconv.Atoi(baseArr[i])
		if version1 == version2 {
			continue
		} else if version1 > version2 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}
