package base62

import (
	"fmt"
	"strings"
)

// 62 进制的转换模块

// 10进制转换62进制
// const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// IntToBase62 将十进制整数转换为 62 进制字符串表示
func IntToBase62(num uint64) string {
	if num == 0 {
		return "0"
	}

	result := []byte{}
	for num > 0 {
		remainder := num % 62
		result = append([]byte{base62Chars[remainder]}, result...) // 前插字符
		num /= 62
	}
	return string(result)
}

// Base62ToInt 将 62 进制字符串转换为十进制整数
func Base62ToInt(s string) (uint64, error) {
	var result uint64 = 0
	for _, ch := range s {
		index := strings.IndexRune(base62Chars, ch)
		if index == -1 {
			return 0, fmt.Errorf("invalid character: %c", ch)
		}
		result = result*62 + uint64(index)
	}
	return result, nil
}
