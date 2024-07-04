/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-20 14:15:06
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-07-02 16:05:02
 * @FilePath: \word2excel\question\rules.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package rules

import (
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"word2excel.io/logger"
	"word2excel.io/question"
)

type RuleThree struct {
	patterns              []string
	FileName              string
	questionNo            int
	size                  int
	currentQuestion       *question.Question
	currentPaperQuestions map[int]question.Question
	AllPaperQuestions     []question.Question
	LastLine              string
}

func NewRuleThreeInstance(config ...any) PaperRule {
	size := 30
	if len(config) > 0 {
		if v, ok := config[0].(int); ok {
			size = v
		}
	}

	return &RuleThree{
		patterns: []string{
			`^(选择一项)(或多项)?([:：]?)$`,
			`^([a-zA-Z])(\.\s*)(.*)`,
			`^(正确答案是)([:：]?\s*)(.*)`,
		},
		FileName:              "",
		questionNo:            0,
		size:                  size,
		currentQuestion:       question.NewQuestion(),
		currentPaperQuestions: make(map[int]question.Question, size),
		AllPaperQuestions:     nil,
		LastLine:              "",
	}
}

func (r *RuleThree) StartParse(text string, name string) {
	r.FileName = name
	isMatched := false

	for index, pattern := range r.patterns {
		reg, _ := regexp.Compile(pattern)

		textMatched := reg.FindStringSubmatch(text)
		if textMatched != nil {
			isMatched = true
			switch index {
			case 0:
				r.ParseType(textMatched)
				r.ParseTitle(textMatched)
			case 1:
				r.ParseOptions(textMatched)
			case 2:
				r.ParseAnswer(textMatched)
				r.ParseAnalysis(textMatched)
			}
		}
	} // end for

	if !isMatched {
		r.AddContent(text)
	}

	r.LastLine = text
}

func (r *RuleThree) ParseType(matched []string) {
	r.currentQuestion.State = question.Q_TYPE

	if matched[2] != "" {
		r.currentQuestion.TypeNo = 2
		r.currentQuestion.TypeName = "多项选择题"
	} else {
		r.currentQuestion.TypeNo = 1
		r.currentQuestion.TypeName = "单项选择题"
	}
}

func (r *RuleThree) ParseTitle(matched []string) {
	r.currentQuestion.Reset()
	r.currentQuestion.State = question.Q_TITLE

	r.questionNo++

	r.currentQuestion.No = r.questionNo

	if r.LastLine != "" {
		p1 := `([（\(]\s*[\)）])`
		r1, _ := regexp.Compile(p1)
		title := r1.ReplaceAllString(r.LastLine, "( )")
		r.currentQuestion.Title = title
	}

	r.AddQuestion(nil)
}

func (r *RuleThree) ParseOptions(matched []string) {
	r.currentQuestion.State = question.Q_OPTIONS

	// optionsNo := matched[1]
	optionText := matched[3]
	// newOption := strings.ToUpper(optionsNo) + ". " + optionText

	options := r.currentQuestion.Options
	options = append(options, optionText)
	r.currentQuestion.Options = options

	r.AddQuestion(nil)
}

func (r *RuleThree) ParseAnswer(matched []string) {
	r.currentQuestion.State = question.Q_ANSWER
	answer := matched[3]

	if strings.Contains(answer, "“对”") || strings.Contains(answer, "“错”") {
		// 判断
		r.currentQuestion.TypeNo = 3
		r.currentQuestion.TypeName = "判断题"
		r.currentQuestion.Options = []string{"正确", "错误"}

		if !strings.Contains(answer, "“对”") {
			r.currentQuestion.Answer = "错误"
		} else {
			r.currentQuestion.Answer = "正确"
		}

	} else if slices.Index([]int{1, 2}, r.currentQuestion.TypeNo) != -1 {
		// 单选，多选
		s1 := strings.Split(answer, ",")
		if len(s1) > 0 {
			// logger.Info("answers: %q, len: %v", s1, len(s1))
			tmp := ""
			for _, v := range s1 {
				for oi, option := range r.currentQuestion.Options {
					if strings.Contains(option, v) {
						tmp += string(byte('A' + oi))
					}
				}
			}
			r.currentQuestion.Answer = tmp
		}
	}

	r.AddQuestion(nil)
}

func (r *RuleThree) ParseAnalysis(matched []string) {
	r.currentQuestion.State = question.Q_ANALYSIS
	analysis := matched[3]

	r.currentQuestion.Analysis = analysis

	r.AddQuestion(nil)
}

func (r *RuleThree) ParsePaperEnd(matched []string) {
	r.currentQuestion.State = question.Q_PAPER_END

	var mkeys []int
	for key := range r.currentPaperQuestions {
		mkeys = append(mkeys, key)
	}
	sort.Ints(mkeys)

	var s1 []question.Question
	for _, v := range mkeys {
		s1 = append(s1, r.currentPaperQuestions[v])
	}

	if r.AllPaperQuestions == nil {
		r.AllPaperQuestions = s1
	} else {
		r.AllPaperQuestions = append(r.AllPaperQuestions, s1...)
	}

	// 清空当前数据
	r.currentPaperQuestions = make(map[int]question.Question, 30)
}

func (r *RuleThree) AddQuestion(qs *question.Question) {
	if qs != nil && qs.No != 0 {
		r.currentPaperQuestions[qs.No] = *qs
	} else if r.currentQuestion.No != 0 {
		r.currentPaperQuestions[r.currentQuestion.No] = *r.currentQuestion
	} else {
		// logger.Info("==error==>试题编号不能为空, file: %v, question: %+v", r.FileName, r.currentQuestion)
		// common.Throw_panic(errors.New(msg))
		logger.Error("试题解析错误, 不符合试题规则：%v， %v", r.FileName, r.currentQuestion)
	}
}

func (r *RuleThree) AddContent(text string) {
	// 未被匹配到的数据行都会到这里, 处理换行数据
	// if r.currentQuestion.State == question.Q_ANALYSIS {
	// 	r.currentQuestion.Analysis += "\n" + text
	// 	r.AddQuestion(nil)
	// }
}

func (r *RuleThree) GetAllQuestions() []question.Question {
	if len(r.currentPaperQuestions) > 0 {
		// word文本结束时，最后一行没有 结束标记时， 移除currentPaperQuestions中的数据到AllPaperQuestions中
		r.ParsePaperEnd(nil)
	}

	fileName := filepath.Base(r.FileName)
	logger.Info("=解析完毕= %v, 试题长度: %v", fileName, len(r.AllPaperQuestions))
	return r.AllPaperQuestions
}
