/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-05-13 10:33:03
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-21 10:07:28
 * @FilePath: \word2excel\question\question.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: htcommon.Pcs://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package question

type Q_STATE int

const (
	Q_TYPE Q_STATE = 1 << iota
	Q_TITLE
	Q_OPTIONS
	Q_ANSWER
	Q_ANALYSIS
	Q_PAPER_END
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

type Question struct {
	No       int      `json:"no"`
	TypeNo   int      `json:"type"`
	TypeName string   `json:"typeName"`
	Title    string   `json:"title"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
	Analysis string   `json:"analysis"`
	state    Q_STATE
}

func NewQuestion() *Question {
	return &Question{}
}

func (qs *Question) reset() {
	qs.No = 0
	qs.Title = ""
	qs.Options = make([]string, 0, 8)
	qs.Answer = ""
	qs.Analysis = ""
}
