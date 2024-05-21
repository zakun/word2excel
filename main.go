/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-05-20 10:23:54
 * @FilePath: \helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"example.io/common"
	"example.io/excel"
	"example.io/word"
)

func main() {
	baseDir, _ := os.Getwd()
	wordDir := filepath.Clean(baseDir + "/runtime/word/")

	common.Pf("word path: %v", wordDir)
	filepath.WalkDir(wordDir, func(name string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			extname := filepath.Ext(name)
			if extname == ".docx" {
				arrQuestion := word.ParseContent(name)

				// common.Pf("questions: %v", arrQuestion)
				// 生成excel文件
				if len(arrQuestion) > 0 {
					excel.GenerateExcelFile(arrQuestion, name)
				}
			}

		}
		return err
	})
}
