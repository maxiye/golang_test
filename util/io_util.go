package util // 大写以被外部访问

import "io/ioutil"

// 大写以被外部访问
func IoReadFile(file string) string {
	if content, err := ioutil.ReadFile(file); err != nil {
		return string(content)
	}
	return ""
}
