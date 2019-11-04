package utils

import (
	"testing"
)

func Test_generator(t *testing.T) {
	url := "https://juejin.im/post/5cea14756fb9a07ee565fb5c"
	num := MD5(url)
	t.Logf("str is: %v", num)

	res := Transport(num)
	t.Logf("res is :%s", res)
}
