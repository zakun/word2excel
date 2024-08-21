/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2022-09-06 13:44:45
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-08-21 15:14:53
 * @FilePath: \helloworld\hello.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"errors"
	"strings"
	"sync"
	"time"

	"path/filepath"

	"word2excel.io/common"
	"word2excel.io/excel"
	"word2excel.io/logger"
	"word2excel.io/question"
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

	var res = make(chan Result, maxN)
	var cfiles = make(chan string, maxN)

	var dirName, ext string
	if strings.Contains(strings.Join([]string{"one", "two", "three"}, "|"), template) {
		dirName = "word"
		ext = ".docx"
	} else {
		dirName = "html"
		ext = ".html"
	}
	path := filepath.Clean("./runtime/" + dirName + "/")
	if ok := common.IsExistDir(path); !ok {
		common.Throw_panic(errors.New(dirName + " 目录不存在：" + path))
	}

	time_start := time.Now()
	logger.Info("==文件处理开始: %v==", time_start.Format("2006-01-02 15:04:05"))

	// 遍历目录
	go walkDir(path, ext, cfiles)

	go parseTempFile(cfiles, res)

	// 生成excel
	var wg sync.WaitGroup
	limit := make(chan int, maxN)
	for v := range res {
		limit <- 1
		wg.Add(1)
		go func(r Result) {
			defer func() {
				wg.Done()
				<-limit
			}()
			excel.GenerateExcelFile(r.Data, r.Name, r.No)
		}(v)
	}
	wg.Wait()

	// excel.GenerateExcelFile(r.Data, r.Name, r.No)
	time_end := time.Now()
	logger.Info("==文件处理结束：%v, %v==", time_end.Format("2006-01-02 15:04:05"), time_end.Sub(time_start).Seconds())
}
