/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-21 10:29:01
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-21 13:58:47
 * @FilePath: \word2excel\question\rule_factory.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package factory

import (
	"errors"

	"example.io/question"
)

func GetRuleInstance(name string, params ...any) (question.Rule, error) {
	switch name {
	case "two":
		// 第二套导题模板
		return question.NewRuleTwoInstance(params...), nil
	}

	return nil, errors.New("规则模板不存在")
}
