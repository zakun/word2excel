/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-20 14:15:06
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-21 10:23:00
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

type Rule struct {
	QuestionNo            int
	size                  int
	CurrentPaperQuestions map[int]Question
	AllPaperQuestions     []Question
}

func GetRuleInstance(size int) *Rule {
	return &Rule{
		QuestionNo:            0,
		size:                  size,
		CurrentPaperQuestions: make(map[int]Question, size),
		AllPaperQuestions:     nil,
	}
}

func (r *Rule) StartParse(text string, qs *Question) {
	isMatched := false

	for index, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		common.PC(err)

		textMatched := reg.FindStringSubmatch(text)
		if textMatched != nil {
			isMatched = true
			switch index {
			case 0:
				r.ParseType(textMatched, qs)
			case 1:
				r.ParseTitle(textMatched, qs)
			case 2:
				r.ParseOptions(textMatched, qs)
			case 3:
				r.ParseAnswer(textMatched, qs)
			case 4:
				r.ParseAnalysis(textMatched, qs)
			case 5:
				r.ParsePaperEnd(textMatched, qs)
			}
			break
		}
	}

	if !isMatched {
		r.AddContent(text, qs)
	}
}

func (r *Rule) ParseType(matched []string, qs *Question) {
	qs.state = Q_TYPE
	typeName := matched[2]

	for key, item := range DicType {
		if found := slices.Index(item, typeName); found != -1 {
			qs.TypeNo = key
			break
		}
	}

	qs.TypeName = typeName
}

func (r *Rule) ParseTitle(matched []string, qs *Question) {
	qs.reset()
	qs.state = Q_TITLE

	// no, err := strconv.Atoi(matched[1])
	// common.PC(err)
	r.QuestionNo++

	qs.No = r.QuestionNo
	qs.Title = matched[2]

	// 判断题固定选项
	if qs.TypeNo == 3 {
		qs.Options = []string{"正确", "错误"}
	}

	r.AddQuestion(qs)
}

func (r Rule) ParseOptions(matched []string, qs *Question) {
	qs.state = Q_OPTIONS
	optionsNo := matched[1]
	optionText := matched[3]
	newOption := optionsNo + ". " + optionText

	options := qs.Options
	options = append(options, newOption)
	qs.Options = options

	// common.PF("current qs: %v, %v", qs.Options, qs.No)
	r.AddQuestion(qs)
}

func (r *Rule) ParseAnswer(matched []string, qs *Question) {
	qs.state = Q_ANSWER
	answer := matched[2]
	if qs.TypeNo == 5 {
		// 填空题, 填空题的答案时保存在选项中
		qs.Options = strings.Split(answer, "@@")
	} else if qs.TypeNo == 3 {
		// 判断题
		correctAns := []string{"√", "T", "t", "正确", "对"}
		if strings.Contains(strings.Join(correctAns, "|"), answer) {
			qs.Answer = "正确"
		} else {
			qs.Answer = "错误"
		}
	} else {
		qs.Answer = answer
	}

	r.AddQuestion(qs)
}

func (r *Rule) ParseAnalysis(matched []string, qs *Question) {
	qs.state = Q_ANALYSIS
	analysis := matched[2]

	qs.Analysis = analysis

	r.AddQuestion(qs)
}

func (r *Rule) ParsePaperEnd(matched []string, qs *Question) {
	qs.state = Q_PAPER_END

	var mkeys []int
	for key := range r.CurrentPaperQuestions {
		mkeys = append(mkeys, key)
	}
	sort.Ints(mkeys)
	// common.PF("current map keys count: %v, %v", len(mkeys), mkeys)

	var s1 []Question
	for _, v := range mkeys {
		s1 = append(s1, r.CurrentPaperQuestions[v])
	}

	if r.AllPaperQuestions == nil {
		r.AllPaperQuestions = s1
	} else {
		r.AllPaperQuestions = append(r.AllPaperQuestions, s1...)
	}

	// common.PF("current all question: %v, %v ", len(r.AllPaperQuestions), r.AllPaperQuestions)

	common.PF("试卷结束: %v - 试题长度：%v", matched[0], len(r.CurrentPaperQuestions))

	logger.Info(common.Sprintf("试卷结束： %v - 试题长度: %v", matched[0], len(r.CurrentPaperQuestions)))

	// common.PF("== all questions == \n %v, %v", len(r.AllPaperQuestions), r.AllPaperQuestions)
	r.CurrentPaperQuestions = make(map[int]Question)
}

func (r *Rule) AddQuestion(qs *Question) {
	if qs.No != 0 {
		r.CurrentPaperQuestions[qs.No] = *qs
	} else {
		common.PC(errors.New("试题编号为空"))
	}
}

func (r *Rule) AddContent(text string, qs *Question) {
	// 未被匹配到的数据行都会到这里
	if qs.state == Q_ANALYSIS {
		qs.Analysis += "\n" + text
		r.AddQuestion(qs)
	}
}
