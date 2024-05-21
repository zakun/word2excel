/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-05-14 11:00:37
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-05-21 14:29:03
 * @FilePath: \word2excel\common\func.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package common

import (
	"fmt"
	"os"
)

func PF(f string, v ...any) {
	fmt.Printf(f+"\n", v...)
}

func PC(err error) {
	if err != nil {
		panic(err)
	}
}

func IsExistDir(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
