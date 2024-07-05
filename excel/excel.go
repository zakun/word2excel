/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-05-17 15:35:01
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-07-05 11:05:37
 * @FilePath: \word2excel\excel\excel.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package excel

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
	"word2excel.io/common"
	"word2excel.io/logger"
	"word2excel.io/question"
)

var fileTitleField = []string{"title", "type", "typeName", "oA", "oB", "oC", "oD", "oE", "oF", "oG", "oH", "oI", "oJ", "oK", "answer", "analysis", "material", "score", "chapter", "keyword", "frequency", "facility", "ability", "year"}

var fileTitleName = []string{"题干", "题型", "题型名称", "A项内容", "B项内容", "C项内容", "D项内容", "E项内容", "F项内容", "G项内容", "H项内容", "I项内容", "J项内容", "K项内容", "答案", "解析", "材料ID", "分值", "章节", "知识点", "考频", "难易度", "考核能力", "真题年份(多个英文逗号,隔开，单个也要加上参照如下)"}

var dicQuestionType = map[int]string{
	1: "单选题",
	2: "多选题",
	3: "判断题",
	4: "简答题",
	5: "填空题",
	6: "写作题",
	7: "论述题",
}

var rowNo = 1
var sheetName = "Sheet1"

func GenerateExcelFile(data []question.Question, name string, no int) {
	defer func() {
		logger.Info("=Excel file: #%v, %v", no, name)
	}()

	rowNo = 1
	excelDir := "./runtime/excel"
	if ok := common.IsExistDir(excelDir); !ok {
		err := os.Mkdir(excelDir, 0755)
		common.Throw_panic(err)
	}

	baseName := filepath.Base(name)
	excelName := baseName[:len(baseName)-5] + ".xlsx"
	excelPathName := filepath.Clean(excelDir + "/" + excelName)

	// common.PF("生成Excel文件: %v\n", excelPathName)
	// logger.Info("Excel 文件: %v", excelPathName)

	f := excelize.NewFile()
	defer func() {
		f.Close()
	}()

	// write file
	WriteHeader(f)

	WriteBody(f, data)

	err := f.SaveAs(excelPathName)
	common.Throw_panic(err)
}

func WriteHeader(f *excelize.File) {
	d := make([]interface{}, len(fileTitleName))
	for i, name := range fileTitleName {
		d[i] = name
	}

	startCell, err := excelize.CoordinatesToCellName(1, rowNo)
	common.Throw_panic(err)

	err = f.SetSheetRow(sheetName, startCell, &d)
	common.Throw_panic(err)

	rowNo += 1
}

func WriteBody(f *excelize.File, data []question.Question) {
	if len(data) > 0 {
		for _, item := range data {
			formatData := FormatRow(item)

			startCell, err := excelize.CoordinatesToCellName(1, rowNo)
			common.Throw_panic(err)

			err = f.SetSheetRow(sheetName, startCell, &formatData)
			common.Throw_panic(err)

			rowNo += 1
		}
	}
}

func FormatRow(question question.Question) []interface{} {
	data := make([]interface{}, len(fileTitleField))

	// "题干", "题型", "题型名称", "A项内容", "B项内容", "C项内容", "D项内容", "E项内容", "F项内容", "G项内容", "H项内容", "I项内容", "J项内容", "K项内容", "答案", "解析", "材料ID", "分值", "章节", "知识点", "考频", "难易度", "考核能力", "真题年份(多个英文逗号,隔开，单个也要加上参照如下)"
	var typeTxt string
	var ok bool
	if typeTxt, ok = dicQuestionType[question.TypeNo]; !ok {
		common.Throw_panic(errors.New("题型名称无法识别"))
	}

	data[0] = question.Title
	data[1] = typeTxt
	data[2] = question.TypeName

	for i := 3; i < 14; i++ {
		if len(question.Options) > 0 && i < len(question.Options)+3 {
			data[i] = question.Options[i-3]
		} else {
			data[i] = ""
		}
	}

	data[14] = question.Answer
	data[15] = question.Analysis

	for di := 16; di < len(data); di++ {
		data[di] = ""
	}
	// common.PF("data: %v", data)
	// os.Exit(1)
	return data
}
