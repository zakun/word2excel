/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-25 12:54:47
 * @FilePath: \go_demo\helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: htcommon.Throw_panics://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package word

import (
	"encoding/json"
	"os"

	"example.io/common"
	"example.io/question"
	"example.io/question/factory"
	"github.com/fumiama/go-docx"
)

func ParseContent(name string) []question.Question {

	file, err := os.Open(name)
	common.Throw_panic(err)
	defer func() {
		file.Close()
	}()

	fileinfo, err := file.Stat()
	common.Throw_panic(err)

	size := fileinfo.Size()
	doc, err := docx.Parse(file, size)
	common.Throw_panic(err)

	rule, err := factory.GetRuleInstance("three", 30)
	common.Throw_panic(err)

	// qs := question.NewQuestion()
	for _, item := range doc.Document.Body.Items {
		switch item.(type) {
		case *docx.Paragraph, *docx.Table:
			text := common.Sprintf("%s", item)
			if text != "" {
				// question.ParseContent(text)
				rule.StartParse(text, name)
			}
		}
	}
	// common.DD("json: %v", ToJson(rule.GetAllQuestions()))
	return rule.GetAllQuestions()
}

func ToJson(data []question.Question) string {
	str, err := json.Marshal(data)
	common.Throw_panic(err)

	return string(str)
}
