/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-20 14:15:06
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-27 16:14:59
 * @FilePath: \word2excel\question\rules.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package rule_one

import (
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"word2excel.io/logger"
	"word2excel.io/question"
)

type RuleOne struct {
	patterns              []string
	FileName              string
	AnswerStart           bool
	size                  int
	currentQuestion       *question.Question
	currentPaperQuestions map[int]question.Question
	AllPaperQuestions     []question.Question
}

func NewRuleOneInstance(config ...any) *RuleOne {
	size := 30
	if len(config) > 0 {
		if v, ok := config[0].(int); ok {
			size = v
		}
	}

	return &RuleOne{
		patterns: []string{
			`^([一二三四五六七八九十]+)、\s*(.*)`,
			`^(\d+)\.\s*(.*)`,
			`^([A-Z])(\.\s*)(.*)`,
			`^%参考答案%`,
			`^(%试卷结束%)$`,
		},
		FileName:              "",
		AnswerStart:           false,
		size:                  size,
		currentQuestion:       question.NewQuestion(),
		currentPaperQuestions: make(map[int]question.Question, size),
		AllPaperQuestions:     nil,
	}
}

func (r *RuleOne) StartParse(text string, name string) {
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
				if !r.AnswerStart {
					r.ParseTitle(textMatched)
				} else {
					r.ParseAnswer(textMatched)
				}
			case 2:
				r.ParseOptions(textMatched)
			case 3:
				r.MarkAnswerStart()
			case 4:
				r.ParsePaperEnd(textMatched)
			}
			break
		}
	}

	if !isMatched {
		r.AddContent(text)
	}
}

func (r *RuleOne) ParseType(matched []string) {
	if !r.AnswerStart {
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
}

func (r *RuleOne) ParseTitle(matched []string) {
	if !r.AnswerStart {
		r.currentQuestion.Reset()
		r.currentQuestion.State = question.Q_TITLE

		r.currentQuestion.No, _ = strconv.Atoi(matched[1])

		title := matched[2]
		if title != "" {
			p1 := `([（\(]\s*[\)）])`
			r1, _ := regexp.Compile(p1)
			title = r1.ReplaceAllString(title, "( )")
		}
		r.currentQuestion.Title = title

		// 判断题固定选项
		if r.currentQuestion.TypeNo == 3 {
			r.currentQuestion.Options = []string{"正确", "错误"}
		}

		r.AddQuestion(nil)
	}

}

func (r *RuleOne) ParseOptions(matched []string) {
	if !r.AnswerStart {
		r.currentQuestion.State = question.Q_OPTIONS

		if slices.Index([]int{1, 2}, r.currentQuestion.TypeNo) != -1 {
			optionsNo := matched[1]
			optionText := matched[3]
			newOption := optionsNo + ". " + optionText

			options := r.currentQuestion.Options
			options = append(options, newOption)
			r.currentQuestion.Options = options

			r.AddQuestion(nil)
		}
	}

}

func (r *RuleOne) MarkAnswerStart() {
	r.AnswerStart = true
}

func (r *RuleOne) ParseAnswer(matched []string) {
	// 解析答案
	if r.AnswerStart {
		answerNo, _ := strconv.Atoi(matched[1])
		q, ok := r.currentPaperQuestions[answerNo]
		if !ok {
			logger.Error("试题解析错误, 答案序号和试题不匹配：%v， %v", r.FileName, r.currentQuestion)
			return
		}
		// 重置当前题目
		r.currentQuestion = &q
		r.currentQuestion.State = question.Q_ANSWER

		answer := matched[2]
		if r.currentQuestion.TypeNo == 5 {
			// 填空题, 填空题的答案保存在选项中
			r.currentQuestion.Options = strings.Split(answer, "@@")
		} else if slices.Index([]int{1, 2, 3}, q.TypeNo) != -1 {
			if r.currentQuestion.TypeNo == 3 {
				// 判断题
				correctAns := []string{"√", "T", "t", "正确", "对"}
				if strings.Contains(strings.Join(correctAns, "|"), answer) {
					r.currentQuestion.Answer = "正确"
				} else {
					r.currentQuestion.Answer = "错误"
				}
			} else {
				// 单选, 多选
				r.currentQuestion.Answer = answer
			}

		} else if slices.Index([]int{4, 6, 7}, r.currentQuestion.TypeNo) != -1 {
			// 简答 写作 论述 没有答案只有解析
			r.currentQuestion.State = question.Q_ANALYSIS
			r.currentQuestion.Analysis = answer
			r.currentQuestion.Answer = ""
		}

		r.AddQuestion(nil)
	}
}

func (r *RuleOne) ParseAnalysis(matched []string) {
	// r.currentQuestion.State = question.Q_ANALYSIS
	// analysis := matched[2]

	// r.currentQuestion.Analysis = analysis

	// r.AddQuestion(nil)
}

func (r *RuleOne) ParsePaperEnd(matched []string) {

	r.currentQuestion.State = question.Q_PAPER_END
	r.AnswerStart = false

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

func (r *RuleOne) AddQuestion(qs *question.Question) {
	if qs != nil && qs.No != 0 {
		r.currentPaperQuestions[qs.No] = *qs
	} else if r.currentQuestion.No != 0 {
		if r.currentQuestion.TypeNo == 0 {
			logger.Error("试题解析错误, 未识别到题型：%v， %v", r.FileName, r.currentQuestion)
		} else {
			r.currentPaperQuestions[r.currentQuestion.No] = *r.currentQuestion
		}
	} else {
		logger.Error("试题解析错误, 不符合试题规则：%v， %v", r.FileName, r.currentQuestion)
	}
}

func (r *RuleOne) AddContent(text string) {
	// 未被匹配到的数据行都会到这里, 处理换行数据
	if r.AnswerStart && r.currentQuestion.State == question.Q_ANALYSIS {
		r.currentQuestion.Analysis += "\n" + text
		r.AddQuestion(nil)
	}
}

func (r *RuleOne) GetAllQuestions() []question.Question {
	if len(r.currentPaperQuestions) > 0 {
		// word文本结束时，最后一行没有 结束标记时， 移除currentPaperQuestions中的数据到AllPaperQuestions中
		r.ParsePaperEnd(nil)
	}

	fileName := filepath.Base(r.FileName)
	logger.Info("=解析完毕= %v, 试题长度: %v", fileName, len(r.AllPaperQuestions))
	return r.AllPaperQuestions
}
