/*
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-06 10:09:40
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-06-24 10:17:51
 * @FilePath: \word2excel\logger\logger.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger
var logfile *os.File

func init() {
	logpath := filepath.Clean("./runtime")
	filename := time.Now().Format("2006-01-02") + "_app.log"

	var err error
	logfile, err = os.OpenFile(logpath+string(os.PathSeparator)+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	if logger == nil {
		logger = log.New(logfile, "", log.LstdFlags)
	}
}

func Info(msg string, params ...any) {
	logger.SetPrefix(" [INFO] ")

	logger.Printf(msg, params...)
	fmt.Printf(msg+"\n", params...)
}

func Close() {
	logfile.Close()
}
