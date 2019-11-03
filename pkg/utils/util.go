package utils

import "crypto/md5"

func MD5(str string) string {
	res := md5.New()
	res.Write([]byte(str))
	return res.EncodeToString(res.Sum(nil))
}

//将数字从十进制转为62进制
func Transport(num int) string {
	bytes := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	var res string

	for num > 0 {
		remainder := num % 62
		res += string(bytes[remainder])
		// if num < 62 {
		//     break
		// }
		num /= 62
	}
	return res
}
