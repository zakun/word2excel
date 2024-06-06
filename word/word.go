/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-06 14:05:14
 * @FilePath: \go_demo\helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: htcommon.Pcs://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package word

import (
	"os"

	"example.io/common"
	"example.io/question"
	"github.com/fumiama/go-docx"
)

func ParseContent(name string) []question.Question {

	file, err := os.Open(name)
	common.PC(err)
	defer func() {
		file.Close()
	}()

	fileinfo, err := file.Stat()
	common.PC(err)

	size := fileinfo.Size()
	doc, err := docx.Parse(file, size)
	common.PC(err)

	// 初始化 question 全局变量
	question.InitWrod()

	for _, item := range doc.Document.Body.Items {
		switch item.(type) {
		case *docx.Paragraph, *docx.Table:
			text := common.Sprintf("%s", item)
			if text != "" {
				question.ParseContent(text)
			}
		}
	}
	return question.AllQuestions
}
