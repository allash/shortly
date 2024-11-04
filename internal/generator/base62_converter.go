package generator

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(num int64) string {
	if num == 0 {
		return "0"
	}

	result := ""
	base := int64(len(base62Chars))

	for num > 0 {
		remainder := num % base
		result = string(base62Chars[remainder]) + result 
		num /= base
	}

	return result
}