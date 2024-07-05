/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-07-05 13:04:07
 * @FilePath: \helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"runtime"
	"sync"
	"time"

	"path/filepath"

	"word2excel.io/common"
	"word2excel.io/excel"
	"word2excel.io/logger"
	"word2excel.io/question"
	"word2excel.io/word"
)

var (
	template string // 模板类型
	maxN     int    // 最大并发数
)

type Result struct {
	No   int
	Name string
	Data []question.Question
}

func main() {
	flag.StringVar(&template, "t", "", "请输入待解析的模板类型,当前仅支持值：two、three; two: 模板2类型，three: 模板3类型")
	flag.IntVar(&maxN, "n", 2, "请输入最大并发数")
	flag.Parse()

	defer func() {
		logger.Close()
	}()

	if maxN <= 0 || maxN > runtime.NumCPU() {
		msg := fmt.Sprintf("并发数需大于0小于cup核数-%v", runtime.NumCPU())
		common.Throw_panic(errors.New(msg))
	}

	if template == "" {
		common.Throw_panic(errors.New("请输入待解析的模板类型"))
	}

	workDir := filepath.Clean("./runtime/word/")
	if ok := common.IsExistDir(workDir); !ok {
		common.Throw_panic(errors.New("word 目录不存在：" + workDir))
	}

	var fileNo int
	var limit = make(chan int, maxN) // 最大并发数控制
	var res = make(chan Result, maxN)

	time_start := time.Now()
	logger.Info("==文件处理开始: %v==", time_start.Format("2006-01-02 15:04:05"))
	go func(path string) {
		// 目录遍历结束，关闭解析 word 结果 channel
		defer close(res)

		var wg sync.WaitGroup
		err := filepath.WalkDir(workDir, func(name string, entry fs.DirEntry, err error) error {
			if err != nil {
				logger.Error("目录遍历错误：%v", err)
				return err
			}

			if !entry.IsDir() {
				extname := filepath.Ext(name)
				if extname == ".docx" {
					fileNo++
					wordFile := word.NewWord(fileNo, name, template)
					// logger.Info("Word 文件: %v# %v", fileNo, name)

					// 解析word文件
					limit <- fileNo
					wg.Add(1)
					go func(wordFile *word.Word, fileNo int) {
						defer func() {
							wg.Done()
							<-limit
						}()

						arrQuestion, ret := wordFile.ParseContent()
						if ret.Status && len(arrQuestion) > 0 {
							// 解析结果
							res <- Result{
								No:   fileNo,
								Name: name,
								Data: arrQuestion,
							}
						} else {
							logger.Info("=解析失败：%v, 试题长度：%v", ret.Msg, len(arrQuestion))
						}
					}(wordFile, fileNo)
				}

			}
			return nil
		})
		// 目录遍历 异常
		if err != nil {
			common.Throw_panic(err)
		}

		// 所有解析word goroutine 结束
		wg.Wait()
	}(workDir)

	var wg2 sync.WaitGroup
	for v := range res {
		wg2.Add(1)
		go func(r Result) {
			defer wg2.Done()
			excel.GenerateExcelFile(r.Data, r.Name, r.No)
		}(v)
	}
	wg2.Wait()

	// excel.GenerateExcelFile(r.Data, r.Name, r.No)
	time_end := time.Now()
	logger.Info("==文件处理结束：%v, %v==", time_end.Format("2006-01-02 15:04:05"), time_end.Sub(time_start).Seconds())

}
