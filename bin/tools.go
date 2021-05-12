package structMaker

import (
	"fmt"
	"os"
	"strings"
)

// 判断文件|目录是否存在
func CheckFileIsExist(file string) bool {
	var exist = true
	if _, err := os.Stat(file); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 打开文件
func OpenFile(filename string) (file *os.File, err error) {
	if CheckFileIsExist(filename) {
		//如果文件存在
		//打开文件
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	} else {
		//创建文件
		file, err = os.Create(filename)
	}
	return file, err
}

// 驼峰格式
func HumpFormat(str string) string {
	var res string
	ary := strings.Split(str, "_")
	for _, v := range ary {
		res += Capitalize(v)
	}
	return res
}

// 首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
