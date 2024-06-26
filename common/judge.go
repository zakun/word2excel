/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-26 10:55:04
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-26 13:08:07
 * @FilePath: \word2excel\common\judge.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package common

import "fmt"

type Judge struct {
	Status bool
	Code   string
	Msg    string
	Err    error
	Params []any
}

func Success() Judge {
	return Judge{
		Status: true,
		Msg:    "success",
		Err:    nil,
		Params: nil,
	}
}

func Fail(msg string, a ...any) Judge {
	return Judge{
		Status: false,
		Msg:    fmt.Sprintf(msg, a...),
		Err:    fmt.Errorf(msg, a...),
		Params: a,
	}
}

func Judging(err error, msg string, a ...any) Judge {
	if err == nil {
		return Success()
	} else {
		msg := msg + ", " + err.Error()
		return Fail(msg, a...)
	}
}

func (j *Judge) SetJudgeCode(c string) {
	j.Code = c
}
