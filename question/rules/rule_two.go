/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-20 14:15:06
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-07-02 15:44:12
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

type RuleTwo struct {
	patterns              []string
	FileName              string
	questionNo            int
	size                  int
	currentQuestion       *question.Question
	currentPaperQuestions map[int]question.Question
	AllPaperQuestions     []question.Question
}

func NewRuleTwoInstance(config ...any) PaperRule {
	size := 30
	if len(config) > 0 {
		if v, ok := config[0].(int); ok {
			size = v
		}
	}

	return &RuleTwo{
		patterns: []string{
			`^([一二三四五六七八九十]+)、\s*(.*)`,
			`^(\d+)\.\s*(.*)`,
			`^([A-Z])(\.\s*)(.*)`,
			`^(答案：\s*)(.*)`,
			`^(解析：\s*)(.*)`,
			`^(%试卷结束%)$`,
		},
		FileName:              "",
		questionNo:            0,
		size:                  size,
		currentQuestion:       question.NewQuestion(),
		currentPaperQuestions: make(map[int]question.Question, size),
		AllPaperQuestions:     nil,
	}
}

func (r *RuleTwo) StartParse(text string, name string) {
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
	r.currentQuestion.State = question.Q_TYPE
	typeName := matched[2]

	for key, item := range question.DicType {
		if found := slices.Index(item, typeName); found != -1 {
			r.currentQuestion.TypeNo = key
			break
		}
	}

	r.currentQuestion.TypeName = typeName
}

func (r *RuleTwo) ParseTitle(matched []string) {
	r.currentQuestion.Reset()
	r.currentQuestion.State = question.Q_TITLE

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
	r.currentQuestion.State = question.Q_OPTIONS
	// optionsNo := matched[1]
	optionText := matched[3]
	// newOption := optionsNo + ". " + optionText

	options := r.currentQuestion.Options
	options = append(options, optionText)
	r.currentQuestion.Options = options

	r.AddQuestion(nil)
}

func (r *RuleTwo) ParseAnswer(matched []string) {
	r.currentQuestion.State = question.Q_ANSWER
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
	r.currentQuestion.State = question.Q_ANALYSIS
	analysis := matched[2]

	r.currentQuestion.Analysis = analysis

	r.AddQuestion(nil)
}

func (r *RuleTwo) ParsePaperEnd(matched []string) {
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

	r.currentPaperQuestions = make(map[int]question.Question, 30)
}

func (r *RuleTwo) AddQuestion(qs *question.Question) {
	if qs != nil && qs.No != 0 {
		r.currentPaperQuestions[qs.No] = *qs
	} else if r.currentQuestion.No != 0 {
		r.currentPaperQuestions[r.currentQuestion.No] = *r.currentQuestion
	} else {
		logger.Error("试题解析错误, 不符合试题规则：%v， %v", r.FileName, r.currentQuestion)
	}
}

func (r *RuleTwo) AddContent(text string) {
	// 未被匹配到的数据行都会到这里, 处理换行数据
	if r.currentQuestion.State == question.Q_ANALYSIS {
		r.currentQuestion.Analysis += "\n" + text
		r.AddQuestion(nil)
	}
}

func (r *RuleTwo) GetAllQuestions() []question.Question {
	if len(r.currentPaperQuestions) > 0 {
		// word文本结束时，最后一行没有 结束标记时， 移除currentPaperQuestions中的数据到AllPaperQuestions中
		r.ParsePaperEnd(nil)
	}

	fileName := filepath.Base(r.FileName)
	logger.Info("=解析完毕= %v, 试题长度: %v", fileName, len(r.AllPaperQuestions))
	return r.AllPaperQuestions
}
