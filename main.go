/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-05-30 17:27:37
 * @FilePath: \helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"errors"
	"io/fs"

	"path/filepath"

	"example.io/common"
	"example.io/excel"
	"example.io/word"
)

func main() {

	wordDir := filepath.Clean("./runtime/word/")
	if ok := common.IsExistDir(wordDir); !ok {
		common.PC(errors.New("word 目录不存在：" + wordDir))
	}

	common.PF("word path: %v", wordDir)
	filepath.WalkDir(wordDir, func(name string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			extname := filepath.Ext(name)
			if extname == ".docx" {
				// common.PF("word file: %v", name)
				arrQuestion := word.ParseContent(name)

				// common.PF("questions: %v", arrQuestion)
				// 生成excel文件
				if len(arrQuestion) > 0 {
					excel.GenerateExcelFile(arrQuestion, name)
				}
			}

		}
		return err
	})
}
