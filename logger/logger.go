package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
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
	color.Printf("<fg=gray;>[%s]</> <bg=blue;> WARN </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][INFO] - %s\n", GetFormatTime(), msg))
}

func Warning(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=yellow;> WARN </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][WARN] - %s\n", GetFormatTime(), msg))
}

func Error(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=red;> WARN </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][ERROR] - %s\n", GetFormatTime(), msg))
}
