/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-07-02 15:44:55
 * @FilePath: \go_demo\helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: htcommon.Throw_panics://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package word

import (
	"encoding/json"
	"os"

	"github.com/fumiama/go-docx"
	"word2excel.io/common"
	"word2excel.io/question"
	"word2excel.io/question/rules"
)

type Word struct {
	No       int
	Name     string
	Template string
}

func NewWord(no int, name string, template string) *Word {
	return &Word{
		No:       no,
		Name:     name,
		Template: template,
	}
}

func (w Word) ParseContent() ([]question.Question, common.Judge) {

	file, err := os.Open(w.Name)
	if err != nil {
		return nil, common.Judging(err, "文件打开失败：%v", w.Name)
	}

	defer func() {
		file.Close()

		// 移除处理过的文件
		// err := os.Remove(w.Name)
		// if err != nil {
		// 	logger.Info("删除失败：%v, err: %v", w.Name, err)
		// }
	}()

	fileinfo, _ := file.Stat()
	size := fileinfo.Size()

	// word文件解析
	doc, err := docx.Parse(file, size)
	if err != nil {
		return nil, common.Judging(err, "word文件解析错误：%v", w.Name)
	}

	// 模板规则解析
	rule, err := rules.GetRuleInstance(w.Template, 30)
	if err != nil {
		return nil, common.Judging(err, "模板规则解析错误：%v", w.Template)
	}

	for _, item := range doc.Document.Body.Items {
		switch item.(type) {
		case *docx.Paragraph, *docx.Table:
			text := common.Sprintf("%s", item)
			if text != "" {
				// question.ParseContent(text)
				rule.StartParse(text, w.Name)
			}
		}
	}
	// common.DD("json: %v", ToJson(rule.GetAllQuestions()))
	return rule.GetAllQuestions(), common.Success()
}

func ToJson(data []question.Question) string {
	str, err := json.Marshal(data)
	common.Throw_panic(err)

	return string(str)
}
