package utils

//将数字从十进制转为62进制
func Transport(num int64) string {
	bytes := []byte("47N9ABdefghC0123DEF8GHuIJpKPQjRSToUWXYZabciklOmnqrst56vwxVyzLM")

	var res string

	for num > 0 {
		remainder := num % 62
		res += string(bytes[remainder])
		num /= 62
	}
	return "http://shortUrl.com/" + res
}
