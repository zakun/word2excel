/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-28 13:58:57
 * @FilePath: \helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"errors"
	"flag"
	"io/fs"
	"sync"
	"time"

	"path/filepath"

	"word2excel.io/common"
	"word2excel.io/excel"
	"word2excel.io/logger"
	"word2excel.io/word"
)

var (
	template string // 模板类型
	maxN     int    // 最大并发数
)

func main() {
	flag.StringVar(&template, "t", "", "请输入待解析的模板类型,当前仅支持值：two、three; two: 模板2类型，three: 模板3类型")
	flag.IntVar(&maxN, "n", 2, "请输入最大并发数")
	flag.Parse()

	defer func() {
		logger.Close()
	}()

	if maxN <= 0 {
		common.Throw_panic(errors.New("并发数至少大于0"))
	}

	if template == "" {
		common.Throw_panic(errors.New("请输入待解析的模板类型"))
	}

	wordDir := filepath.Clean("./runtime/word/")
	if ok := common.IsExistDir(wordDir); !ok {
		common.Throw_panic(errors.New("word 目录不存在：" + wordDir))
	}

	var fileNo int
	var taskLimit = make(chan int, maxN) // 最大并发数控制
	var wg sync.WaitGroup

	time_start := time.Now()
	logger.Info("==文件处理开始: %v==", time_start.Format("2006-01-02 15:04:05"))
	filepath.WalkDir(wordDir, func(name string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			extname := filepath.Ext(name)
			if extname == ".docx" {
				fileNo++
				wordFile := word.NewWord(fileNo, name, template)
				logger.Info("Word 文件: %v# %v", fileNo, name)

				taskLimit <- fileNo
				wg.Add(1)
				go func(c chan int) {
					defer wg.Done()

					arrQuestion, ret := wordFile.ParseContent()
					if ret.Status && len(arrQuestion) > 0 {
						// 生成Excel文件
						excel.GenerateExcelFile(arrQuestion, name)
					} else {
						logger.Info(ret.Msg)
					}

					<-c
				}(taskLimit)
			}

		}
		return err
	})

	wg.Wait()
	time_end := time.Now()
	logger.Info("==文件处理结束：%v, %v==", time_end.Format("2006-01-02 15:04:05"), time_end.Sub(time_start).Seconds())
}
