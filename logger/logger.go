// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)

var (
	LogChan = make(chan string, 100)
)

func GetLogFileName() string {
	return fmt.Sprintf("logs/log_%s.log", GetFormatDate())
}

func GetFormatDate() string {
	time := time.Now()
	year := time.Year()
	month := time.Month()
	day := time.Day()

	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func GetFormatTime() string {
	time := time.Now()
	hour := time.Local().Hour()
	minute := time.Minute()
	sec := time.Second()

	var second string

	if sec < 10 {
		second = fmt.Sprintf("0%d", sec)
	} else {
		second = fmt.Sprintf("%d", sec)
	}

	return fmt.Sprintf("%d:%d:%s", hour, minute, second)
}

func InitLogFile() {
	_, logsFolderStatErr := os.Stat("logs")
	if os.IsNotExist(logsFolderStatErr) {
		os.Mkdir("logs", 0755)
	}

	_, logFileStatErr := os.Stat(GetLogFileName())
	if os.IsNotExist(logFileStatErr) {
		os.Create(GetLogFileName())
	}
}

func WriteToFile(message string) {
	file, err := os.OpenFile(GetLogFileName(), os.O_APPEND|os.O_WRONLY, 0644)
	if os.IsNotExist(err) {
		InitLogFile()
	}
	defer file.Close()

	file.WriteString(message)
}

func Info(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=blue;> INFO </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][INFO] - %s\n", GetFormatTime(), msg))
}

func Warning(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=yellow;> WARN </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][WARN] - %s\n", GetFormatTime(), msg))
}

func Error(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=red;> ERROR </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][ERROR] - %s\n", GetFormatTime(), msg))
}

func LogRequest(c *gin.Context) {
	statusCode := strconv.Itoa(c.Writer.Status())
	method := strings.Trim(c.Request.Method, " ")
	remoteAddr := c.ClientIP()
	url := c.Request.URL.Path

	statusCodeLog := fmt.Sprintf("<bg=green;> %s </>", statusCode)
	methodLog := fmt.Sprintf("<bg=blue;> %s     </>", method)

	if strings.HasPrefix(statusCode, "4") {
		statusCodeLog = fmt.Sprintf("<bg=yellow;> %s </>", statusCode)
	} else if strings.HasPrefix(statusCode, "3") {
		statusCodeLog = fmt.Sprintf("<bg=blue;> %s </>", statusCode)
	}

	switch method {
	case "POST":
		methodLog = fmt.Sprintf("<bg=yellow;> %s    </>", method)
	case "DELETE":
		methodLog = fmt.Sprintf("<bg=red;> %s   </>", method)
	case "OPTIONS":
		methodLog = fmt.Sprintf("<bg=gray;> %s </>", method)
	}

	msg := fmt.Sprintf("|%s| %15s |%s %s", statusCodeLog, remoteAddr, methodLog, url)
	color.Printf("<fg=gray;>[%s]</> <bg=blue;> INFO </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][INFO] - | %s | %s | %s %s\n", GetFormatTime(), statusCode, remoteAddr, method, url))
}
