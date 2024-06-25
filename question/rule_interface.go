/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-20 14:15:06
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-21 13:03:42
 * @FilePath: \word2excel\question\rules.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package question

type Rule interface {
	StartParse(string, string)

	ParseType([]string)
	ParseTitle([]string)
	ParseOptions([]string)
	ParseAnswer([]string)
	ParseAnalysis([]string)
	ParsePaperEnd([]string)

	AddQuestion(*Question)
	AddContent(string)

	GetAllQuestions() []Question
}
