package utils

import (
	"MyProject/Short_Url/models"
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	md5String := md5.New()
	md5String.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
	// temp := ""
	// var sum, sumToAdd int
	// i := 0
	// for _, v := range md5String.Sum(nil) {
	// 	if i != 0 && i%3 == 0 {
	// 		temp += strconv.Itoa(sum)
	// 		sum = 0
	// 	} else {
	// 		sumToAdd = int(v)
	// 		sum += sumToAdd
	// 	}

	// 	i++
	// }

	// res, err := strconv.ParseInt(string(temp), 10, 64)
	// if err != nil {

	// 	fmt.Printf("err: %s\n", err)
	// 	return -1

	// }

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

func CodeGenerator(url string, userId int) (string, error) {
	urlCode := models.UrlCode
	urlId, err := urlCode.AddUrl(url, userId)
	if err != nil {
		return "", err
	}

	shortCode := Transport(urlId)

	err = urlCode.UpdateCode(urlId, shortCode)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}
