/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-08-02 16:07:06
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
	defer stop()

	// 解析命令行参数
	parseFlag()

	path := filepath.Clean("./runtime/word/")
	if ok := common.IsExistDir(path); !ok {
		common.Throw_panic(errors.New("word 目录不存在：" + path))
	}

	var res = make(chan Result, maxN)
	var docx = make(chan string, maxN)

	time_start := time.Now()
	logger.Info("==文件处理开始: %v==", time_start.Format("2006-01-02 15:04:05"))

	// 遍历目录
	go walkDir(path, docx)

	// 解析world
	go parseWords(docx, res)

	// 生成excel
	var wg sync.WaitGroup
	for v := range res {
		wg.Add(1)
		go func(r Result) {
			defer wg.Done()
			excel.GenerateExcelFile(r.Data, r.Name, r.No)
		}(v)
	}
	wg.Wait()

	// excel.GenerateExcelFile(r.Data, r.Name, r.No)
	time_end := time.Now()
	logger.Info("==文件处理结束：%v, %v==", time_end.Format("2006-01-02 15:04:05"), time_end.Sub(time_start).Seconds())

}

func parseWords(docx <-chan string, res chan<- Result) {
	defer close(res)

	var wg sync.WaitGroup
	var limit = make(chan int, maxN) // 最大并发数控制
	fileNo := 0

	for name := range docx {
		fileNo++
		wordFile := word.NewWord(fileNo, name, template)

		limit <- 1
		wg.Add(1)
		// word文件解析
		go func() {
			defer func() {
				wg.Done()
				<-limit
			}()

			arrQuestion, ret := wordFile.ParseContent()
			if ret.Status && len(arrQuestion) > 0 {
				// 解析结果
				res <- Result{
					No:   wordFile.No,
					Name: wordFile.Name,
					Data: arrQuestion,
				}
			} else {
				logger.Info("=解析失败：%v, 试题长度：%v", ret.Msg, len(arrQuestion))
			}
		}()
	}
	// 等待
	wg.Wait()
}

func walkDir(path string, docx chan<- string) {
	defer close(docx)

	err := filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			logger.Error("目录遍历错误：%v", err)
			return err
		}

		if !entry.IsDir() {
			extname := filepath.Ext(path)
			if extname == ".docx" {
				docx <- path
			}
		}
		return nil
	})
	common.Throw_panic(err)
}

func parseFlag() {
	flag.StringVar(&template, "t", "", "请输入待解析的模板类型,当前仅支持值：one, two, three; 分别对应三种模板类型")
	flag.IntVar(&maxN, "n", 2, "请输入最大并发数")
	flag.Parse()

	if maxN <= 0 || maxN > runtime.NumCPU() {
		msg := fmt.Sprintf("并发数需大于0小于cup核数-%v", runtime.NumCPU())
		common.Throw_panic(errors.New(msg))
	}

	if template == "" {
		common.Throw_panic(errors.New("请输入待解析的模板类型"))
	}
}

func stop() {
	logger.Close()
}
