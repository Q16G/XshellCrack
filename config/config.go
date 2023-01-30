package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func Get_Context(file string) ([]string, error) { //(string, error)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		err2 := errors.New("not found this file")
		return []string{}, err2
	}
	s := strings.Split(string(b), "\r\n")
	return s, nil
}
func Errlog_print(str string) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	color.Red("%v   %v", timeStr, str)
	os.Exit(0)
}
func Checkerror(err error) {
	if err != nil {
		Errlog_print(err.Error())
	}
}
func Loginfo_print(str string) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	color.Green("%v   %v", timeStr, str)
}
