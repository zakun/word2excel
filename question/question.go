/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-05-13 10:33:03
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-06 13:45:40
 * @FilePath: \word2excel\question\question.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: htcommon.Pcs://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package question

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"example.io/common"
	"example.io/logger"
)

var DicType = map[int][]string{
	1: {"单项选择题", "选择填空题", "词汇与结构", "阅读程序题", "交际用语", "词汇语法", "阅读理解", "单选题", "完形填空"},
	2: {"多项选择题", "阅读理解题", "连线题", "多选题", "不定项选择题"},
	3: {"判断题", "判断题（判断正误，并说明理由）"},
	4: {"简答题", "编程题", "翻译", "概念题", "绘图题", "计算题", "名词解释", "设计题", "识图题", "应用题"},
	5: {"填空题", "词汇语法"},
	6: {"写作题", "写作"},
	7: {"论述题", "论述", "论文", "分析题", "逻辑公式翻译", "业务处理题", "综合题"},
}

var qs = NewQuestion()
var mapQs = make(map[int]Question, 30)     //当前 试卷 试题
var AllQuestions = make([]Question, 0, 30) // 当前文档内所有试卷试题  - 单个文档内可能有多套试卷

type Question struct {
	No       int      `json:"no"`
	TypeNo   int      `json:"type"`
	TypeName string   `json:"typeName"`
	Title    string   `json:"title"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
	Analysis string   `json:"analysis"`
	state    string
}

func NewQuestion() *Question {
	return &Question{}
}

func (qs *Question) ParseType(matched []string) {
	qs.state = "qs_type"
	typeName := matched[2]

	for key, item := range DicType {
		if found := slices.Index(item, typeName); found != -1 {
			qs.TypeNo = key
			break
		}
	}

	qs.TypeName = typeName
}

func (qs *Question) ParseTitle(matched []string) {
	qs.reset()
	qs.state = "qs_title"

	no, err := strconv.Atoi(matched[1])
	common.PC(err)
	qs.No = no
	qs.Title = matched[2]

	// 判断题固定选项
	if qs.TypeNo == 3 {
		qs.Options = []string{"正确", "错误"}
	}

	qs.AddQuestion()
}

func (qs *Question) ParseOptions(matched []string) {
	qs.state = "qs_option"
	optionsNo := matched[1]
	optionText := matched[3]
	newOption := optionsNo + ". " + optionText

	options := qs.Options
	options = append(options, newOption)
	qs.Options = options

	// common.PF("current qs: %v, %v", qs.Options, qs.No)
	qs.AddQuestion()
}

func (qs *Question) ParseAnswer(matched []string) {
	qs.state = "qs_answer"
	answer := matched[2]
	if qs.TypeNo == 5 {
		qs.Options = strings.Split(answer, "@@")
	} else if qs.TypeNo == 3 {
		correctAns := []string{"√", "T", "t", "正确", "对"}
		if strings.Contains(strings.Join(correctAns, "|"), answer) {
			qs.Answer = "正确"
		} else {
			qs.Answer = "错误"
		}
	} else {
		qs.Answer = answer
	}

	qs.AddQuestion()
}

func (qs *Question) ParseAnalysis(matched []string) {
	qs.state = "qs_analysis"
	analysis := matched[2]

	qs.Analysis = analysis

	qs.AddQuestion()
}

func (qs *Question) ParsePaperEnd(matched []string) {
	qs.state = "qs_paper_end"

	for i := 1; i <= len(mapQs); i++ {
		AllQuestions = append(AllQuestions, mapQs[i])
	}

	common.PF("试卷结束: %v - 试题长度：%v", matched[0], len(mapQs))

	logger.Info(common.Sprintf("试卷结束： %v - 试题长度: %v", matched[0], len(mapQs)))
	// 初始化当前试题 map
	mapQs = make(map[int]Question, 30)
}

func ParseContent(text string) {
	patterns := []string{
		`^([一二三四五六七八九十]+)、\s*(.*)`,
		`^(\d+)\.\s*(.*)`,
		`^([A-Z])(\.\s*)(.*)`,
		`^(答案：\s*)(.*)`,
		`^(解析：\s*)(.*)`,
		`^(%试卷结束%)$`,
	}
	isMatched := false

	for index, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		common.PC(err)

		textMatched := reg.FindStringSubmatch(text)
		if textMatched != nil {
			isMatched = true
			switch index {
			case 0:
				qs.ParseType(textMatched)
			case 1:
				qs.ParseTitle(textMatched)
			case 2:
				qs.ParseOptions(textMatched)
			case 3:
				qs.ParseAnswer(textMatched)
			case 4:
				qs.ParseAnalysis(textMatched)
			case 5:
				qs.ParsePaperEnd(textMatched)
			}
			break
		}
	}

	if !isMatched {
		qs.AddContent(text)
	}

}

// 初始化 全局变量
func InitWrod() {
	mapQs = make(map[int]Question, 30)
	AllQuestions = make([]Question, 0, 30)
}

func (qs Question) AddQuestion() {
	if qs.No != 0 {
		mapQs[qs.No] = qs
	} else {
		common.PC(errors.New("试题编号为空"))
	}
}

func (qs *Question) AddContent(text string) {
	if qs.state == "qs_analysis" {
		qs.Analysis += "\n" + text
		qs.AddQuestion()
	}
}

func (qs *Question) reset() {
	qs.No = 0
	qs.Title = ""
	qs.Options = make([]string, 0, 8)
	qs.Answer = ""
	qs.Analysis = ""
}
