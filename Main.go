package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
)

var (
	WriterChan chan bool
	isActive   = false
	LogDir     = ""
)

func init() {
	_, err := exec.Command("xdotool", "getwindowfocus", "getwindowname").Output()
	if err != nil {
		fmt.Println("Failed to execute command \"xdotool getwindowfocus getwindowname\". Do you have xdotool?")
		fmt.Println(err)
		os.Exit(-1)
	}

	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Please provide only path of log directory.")
		os.Exit(-1)
	}

	if !Exists(args[0]) {
		fmt.Println("Specified path doesn't exist.")
		os.Exit(-1)
	}
	LogDir = flag.Args()[0]
}

func main() {
	fmt.Println("Running...")
	go Writer()

	for {
		EvChan := hook.Start()

		<-EvChan

		isActive = true
		hook.End()

		<-WriterChan
	}

}

func Writer() {
	for {
		time.Sleep(time.Minute)
		var (
			now     time.Time = time.Now()
			dir     string    = fmt.Sprintf("%s/%d/%d/", LogDir, now.Year(), now.Month())
			fdir    string    = dir + fmt.Sprintf("%d.log", now.Day())
			LogText string    = fmt.Sprintf("%02d:%02d %s \"%s\"", now.Hour(), now.Minute(), BoolToStr(isActive), GetCurrentWindow())
			file    *os.File  = nil
			err     error     = nil
		)

		if !Exists(dir) {
			err = os.MkdirAll(dir, 0777)
			HandleError(err)
		}

		if Exists(fdir) {
			file, err = os.OpenFile(fdir, os.O_WRONLY|os.O_APPEND, 0777)
		} else {
			file, err = os.Create(fdir)
		}
		HandleError(err)

		fmt.Println(LogText)
		_, err = fmt.Fprintln(file, LogText)
		HandleError(err)

		err = file.Close()
		HandleError(err)

		isActive = false
		WriterChan <- true
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetCurrentWindow() string {
	out, err := exec.Command("xdotool", "getwindowfocus", "getwindowname").Output()
	HandleError(err)

	return strings.ReplaceAll(string(out), "\n", "")
}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func BoolToStr(b bool) string {
	if b {
		return "true "
	} else {
		return "false"
	}
}
