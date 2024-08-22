/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-08-19 11:10:19
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-08-22 10:55:37
 * @FilePath: \word2excel\functions.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"

	"word2excel.io/common"
	"word2excel.io/html"
	"word2excel.io/logger"
	"word2excel.io/question"
	"word2excel.io/word"
)

type TempFile interface {
	ParseTemplateContent() ([]question.Question, common.Judge)
}

func parseTempFile(cfile <-chan string, res chan<- Result) {
	defer close(res)

	var wg sync.WaitGroup
	var limit = make(chan int, maxN) // 最大并发数控制
	fileNo := 0

	for name := range cfile {
		fileNo++
		tempName := name

		var tempFile TempFile
		tempFile = html.NewHtml(fileNo, tempName, template)
		if strings.Contains(strings.Join([]string{"one", "two", "three"}, "|"), template) {
			tempFile = word.NewWord(fileNo, tempName, template)
		}

		limit <- 1
		wg.Add(1)
		// 模板文件 解析
		go func() {
			defer func() {
				wg.Done()
				<-limit
			}()

			arrQuestion, ret := tempFile.ParseTemplateContent()

			// jsons := question.ToJson(arrQuestion)
			// common.DD("json txt: %v", jsons)
			// os.Exit(0)
			if ret.Status && len(arrQuestion) > 0 {
				// 解析结果
				res <- Result{
					No:   fileNo,
					Name: tempName,
					Data: arrQuestion,
				}
			} else {
				logger.Info("=解析失败: %v, 试题长度：%v", ret.Msg, len(arrQuestion))
			}
		}()
	}
	// 等待
	wg.Wait()
}

func walkDir(path string, ext string, docx chan<- string) {
	defer close(docx)

	err := filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			logger.Error("目录遍历错误：%v", err)
			return err
		}

		if !entry.IsDir() {
			extname := filepath.Ext(path)
			if strings.Contains(extname, ext) {
				docx <- path
			}
		}
		return nil
	})
	common.Throw_panic(err)
}

func parseFlag() {
	flag.StringVar(&template, "t", "", "请输入待解析的模板类型,当前仅支持值：one, two, three, html; 分别对应四种模板类型")
	flag.IntVar(&maxN, "n", 2, "请输入最大并发数")
	flag.Parse()

	if maxN <= 0 || maxN > runtime.NumCPU() {
		msg := fmt.Sprintf("并发数需大于0小于cup核数-%v", runtime.NumCPU())
		common.Throw_panic(errors.New(msg))
	}

	allTemplateType := append(dTemplate["word"], dTemplate["html"]...)
	if !slices.Contains(allTemplateType, template) {
		common.Throw_panic(errors.New("请输入待解析的模板类型: " + strings.Join(allTemplateType, ",")))
	}
}

func stop() {
	logger.Close()
}
