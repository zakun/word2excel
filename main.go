/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-25 13:11:47
 * @FilePath: \helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"errors"
	"io/fs"
	"sync"

	"path/filepath"

	"example.io/common"
	"example.io/excel"
	"example.io/logger"
	"example.io/word"
)

func main() {
	defer func() {
		logger.Close()
	}()

	wordDir := filepath.Clean("./runtime/word/")
	if ok := common.IsExistDir(wordDir); !ok {
		common.Throw_panic(errors.New("word 目录不存在：" + wordDir))
	}

	var fileNo int
	var taskLimit = make(chan int, 6) // 最大并发数控制
	var wg sync.WaitGroup

	filepath.WalkDir(wordDir, func(name string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			extname := filepath.Ext(name)
			if extname == ".docx" {
				fileNo++
				logger.Info("Word 文件: %v# %v", fileNo, name)

				taskLimit <- fileNo
				wg.Add(1)

				go func(c chan int) {
					defer wg.Done()

					arrQuestion := word.ParseContent(name)
					// 生成excel文件
					if len(arrQuestion) > 0 {
						excel.GenerateExcelFile(arrQuestion, name)
					} else {
						logger.Info("为解析出试题：Word 文件: %v# %v", fileNo, name)
					}

					<-c
				}(taskLimit)
			}

		}
		return err
	})

	wg.Wait()
}
