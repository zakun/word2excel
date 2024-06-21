/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-20 14:15:06
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-21 14:12:59
 * @FilePath: \word2excel\question\rules.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package question

import (
	"errors"
	"regexp"
	"slices"
	"sort"
	"strings"

	"example.io/common"
	"example.io/logger"
)

var (
	patterns = []string{
		`^([一二三四五六七八九十]+)、\s*(.*)`,
		`^(\d+)\.\s*(.*)`,
		`^([A-Z])(\.\s*)(.*)`,
		`^(答案：\s*)(.*)`,
		`^(解析：\s*)(.*)`,
		`^(%试卷结束%)$`,
	}
)

type RuleTwo struct {
	questionNo            int
	size                  int
	currentQuestion       *Question
	currentPaperQuestions map[int]Question
	AllPaperQuestions     []Question
}

func NewRuleTwoInstance(config ...any) *RuleTwo {
	size := 30
	if len(config) > 0 {
		if v, ok := config[0].(int); ok {
			size = v
		}
	}

	return &RuleTwo{
		questionNo:            0,
		size:                  size,
		currentQuestion:       NewQuestion(),
		currentPaperQuestions: make(map[int]Question, size),
		AllPaperQuestions:     nil,
	}
}

func (r *RuleTwo) StartParse(text string) {
	isMatched := false

	for index, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		common.PC(err)

		textMatched := reg.FindStringSubmatch(text)
		if textMatched != nil {
			isMatched = true
			switch index {
			case 0:
				r.ParseType(textMatched)
			case 1:
				r.ParseTitle(textMatched)
			case 2:
				r.ParseOptions(textMatched)
			case 3:
				r.ParseAnswer(textMatched)
			case 4:
				r.ParseAnalysis(textMatched)
			case 5:
				r.ParsePaperEnd(textMatched)
			}
			break
		}
	}

	if !isMatched {
		r.AddContent(text)
	}
}

func (r *RuleTwo) ParseType(matched []string) {
	r.currentQuestion.state = Q_TYPE
	typeName := matched[2]

	for key, item := range DicType {
		if found := slices.Index(item, typeName); found != -1 {
			r.currentQuestion.TypeNo = key
			break
		}
	}

	r.currentQuestion.TypeName = typeName
}

func (r *RuleTwo) ParseTitle(matched []string) {
	r.currentQuestion.reset()
	r.currentQuestion.state = Q_TITLE

	r.questionNo++

	r.currentQuestion.No = r.questionNo
	r.currentQuestion.Title = matched[2]

	// 判断题固定选项
	if r.currentQuestion.TypeNo == 3 {
		r.currentQuestion.Options = []string{"正确", "错误"}
	}

	r.AddQuestion(nil)
}

func (r *RuleTwo) ParseOptions(matched []string) {
	r.currentQuestion.state = Q_OPTIONS
	optionsNo := matched[1]
	optionText := matched[3]
	newOption := optionsNo + ". " + optionText

	options := r.currentQuestion.Options
	options = append(options, newOption)
	r.currentQuestion.Options = options

	r.AddQuestion(nil)
}

func (r *RuleTwo) ParseAnswer(matched []string) {
	r.currentQuestion.state = Q_ANSWER
	answer := matched[2]
	if r.currentQuestion.TypeNo == 5 {
		// 填空题, 填空题的答案保存在选项中
		r.currentQuestion.Options = strings.Split(answer, "@@")
	} else if r.currentQuestion.TypeNo == 3 {
		// 判断题
		correctAns := []string{"√", "T", "t", "正确", "对"}
		if strings.Contains(strings.Join(correctAns, "|"), answer) {
			r.currentQuestion.Answer = "正确"
		} else {
			r.currentQuestion.Answer = "错误"
		}
	} else {
		r.currentQuestion.Answer = answer
	}

	r.AddQuestion(nil)
}

func (r *RuleTwo) ParseAnalysis(matched []string) {
	r.currentQuestion.state = Q_ANALYSIS
	analysis := matched[2]

	r.currentQuestion.Analysis = analysis

	r.AddQuestion(nil)
}

func (r *RuleTwo) ParsePaperEnd(matched []string) {
	r.currentQuestion.state = Q_PAPER_END

	var mkeys []int
	for key := range r.currentPaperQuestions {
		mkeys = append(mkeys, key)
	}
	sort.Ints(mkeys)

	var s1 []Question
	for _, v := range mkeys {
		s1 = append(s1, r.currentPaperQuestions[v])
	}

	if r.AllPaperQuestions == nil {
		r.AllPaperQuestions = s1
	} else {
		r.AllPaperQuestions = append(r.AllPaperQuestions, s1...)
	}

	// common.PF("current all question: %v, %v ", len(r.AllPaperQuestions), r.AllPaperQuestions)

	common.PF("试卷结束: %v - 试题长度：%v", matched[0], len(r.currentPaperQuestions))

	logger.Info(common.Sprintf("试卷结束： %v - 试题长度: %v", matched[0], len(r.currentPaperQuestions)))

	// common.PF("== all questions == \n %v, %v", len(r.AllPaperQuestions), r.AllPaperQuestions)
	r.currentPaperQuestions = make(map[int]Question, 30)
}

func (r *RuleTwo) AddQuestion(qs *Question) {
	if qs != nil && qs.No != 0 {
		r.currentPaperQuestions[qs.No] = *qs
	} else if r.currentQuestion.No != 0 {
		r.currentPaperQuestions[r.currentQuestion.No] = *r.currentQuestion
	} else {
		common.PC(errors.New("试题编号不能为空"))
	}
}

func (r *RuleTwo) AddContent(text string) {
	// 未被匹配到的数据行都会到这里
	if r.currentQuestion.state == Q_ANALYSIS {
		r.currentQuestion.Analysis += "\n" + text
		r.AddQuestion(nil)
	}
}

func (r RuleTwo) GetAllQuestions() []Question {
	return r.AllPaperQuestions
}