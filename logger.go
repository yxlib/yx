// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	// "syscall"
)

const (
	LOG_BATCH_DUMP_COUNT       = 100
	LOG_MAX_CACHE_SIZE         = 64 * 1024
	LOG_DEFAULT_DUMP_SIZE      = 32 * 1024 * 1024
	LOG_DEFAULT_DUMP_THRESHOLD = 32 * 1024
	LOG_DEFAULT_DUMP_INTV      = 100
)

const LOG_DEBUG_SWITCH_FILE = "debug.sf"

type LogLv = int

const (
	LOG_LV_DEBUG LogLv = 0
	LOG_LV_INFO  LogLv = 1
	LOG_LV_WARN  LogLv = 2
	LOG_LV_ERROR LogLv = 3
)

//========================
//    global method
//========================
func StartLogger() {
	go loggerInst.loop()
}

func StopLogger() {
	loggerInst.stop()
}

// Start dump log.
// @param file, the relative/full path of a file.
// @param dumpFileSize, max size of the dump file.
// @param dumpThreshold, max count of logs in buffer to cause dump.
// @param dumpIntervalMs, dump interval in millisecond.
func StartDumpLog(file string, dumpFileSize int, dumpThreshold int, dumpIntervalMs uint32) {
	loggerInst.startDump(file, dumpFileSize, dumpThreshold, dumpIntervalMs)
}

// Start dump log by default params.
// @param file, the relative/full path of a file.
func StartDumpLogDefault(file string) {
	loggerInst.startDump(file, LOG_DEFAULT_DUMP_SIZE, LOG_DEFAULT_DUMP_THRESHOLD, LOG_DEFAULT_DUMP_INTV)
}

// Stop dump log.
func StopDumpLog() {
	loggerInst.stopDump()
}

// Set log level.
// @param lv, the level to begin print.
func SetLogLevel(lv LogLv) {
	loggerInst.SetLevel(lv)
}

// Set to PowerShell mode.
// func SetPowerShellMode() {
// 	loggerInst.SetPowerShellMode()
// }

// Set a new print func to instead of the default one.
// eg: PowerShell print.
//
// var (
// 	kernel32                *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
// 	SetConsoleTextAttribute *syscall.LazyProc = nil
// 	CloseHandle             *syscall.LazyProc = nil
// 	BgColor   int   = 0x50
// 	FontColor Color = Color{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}
// )
//
// type Color struct {
// 	black        int
// 	blue         int
// 	green        int
// 	cyan         int
// 	red          int
// 	purple       int
// 	yellow       int
// 	light_gray   int
// 	gray         int
// 	light_blue   int
// 	light_green  int
// 	light_cyan   int
// 	light_red    int
// 	light_purple int
// 	light_yellow int
// 	white        int
// }
//
// func PowerShellPrint(lv yx.LogLv, logStr string) {
// 	if SetConsoleTextAttribute == nil {
// 		SetConsoleTextAttribute = kernel32.NewProc(`SetConsoleTextAttribute`)
// 	}
//
// 	if CloseHandle == nil {
// 		CloseHandle = kernel32.NewProc(`CloseHandle`)
// 	}
//
// 	color := FontColor.white
// 	if lv == yx.LOG_LV_ERROR {
// 		color = FontColor.light_red
// 	} else if lv == yx.LOG_LV_WARN {
// 		color = FontColor.light_yellow
// 	}
//
// 	handle, _, _ := SetConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(BgColor|color))
// 	fmt.Print(logStr)
// 	CloseHandle.Call(handle)
// }
func SetPrintFunc(printFunc func(lv LogLv, logStr string)) {
	loggerInst.SetPrintFunc(printFunc)
}

//========================
//    log config
//========================
type LogConf struct {
	Level int `json:"level"`
	// IsPowerShellRun bool   `json:"power_shell_run"`
	IsDump        bool   `json:"is_dump"`
	DumpPath      string `json:"dump_path"`
	DumpFileSize  int    `json:"dump_file_size"`
	DumpThreshold int    `json:"dump_threshold"`
	DumpInterval  uint32 `json:"dump_interval"`
}

func ConfigLogger(cfg *LogConf, printFunc func(lv LogLv, logStr string)) {
	SetLogLevel(cfg.Level)
	SetPrintFunc(printFunc)
	// if cfg.IsPowerShellRun {
	// 	SetPowerShellMode()
	// }

	if cfg.IsDump {
		StartDumpLog(cfg.DumpPath, cfg.DumpFileSize, cfg.DumpThreshold, cfg.DumpInterval)
	}
}

//========================
//        Logger
//========================
type Logger struct {
	tag string
}

func NewLogger(tag string) *Logger {
	return &Logger{
		tag: tag,
	}
}

// Print debug log.
func (l *Logger) D(a ...interface{}) {
	loggerInst.D(l.tag, a...)
}

// Print infomation log.
func (l *Logger) I(a ...interface{}) {
	loggerInst.I(l.tag, a...)
}

// Print warn log.
func (l *Logger) W(a ...interface{}) {
	loggerInst.W(l.tag, a...)
}

// Print error log.
func (l *Logger) E(a ...interface{}) {
	loggerInst.E(l.tag, a...)
}

// Print detail log.
func (l *Logger) Detail(lv LogLv, s string) {
	loggerInst.Detail(lv, s)
}

// // Print ln.
// func (l *Logger) Ln() {
// 	loggerInst.Ln()
// }

//==============================================
//                   logger
//==============================================
type LogInfo struct {
	Lv     LogLv
	LogStr string
}

type logger struct {
	level           LogLv
	bPowerShellMode bool
	printFunc       func(lv LogLv, logStr string)
	bDumpOpen       bool
	strDumpFile     string
	dumpFileSno     uint64
	dumpFileSize    int
	dumpThreshold   int
	dumpIntervalMs  uint32
	// queLogs         chan string
	lck           *sync.Mutex
	queLogs       []*LogInfo
	writeLogs     []*LogInfo
	evtDumpToFile *Event
	evtStop       *Event
	evtStopSucc   *Event
}

var loggerInst = &logger{
	level:           LOG_LV_DEBUG,
	bPowerShellMode: false,
	printFunc:       nil,
	bDumpOpen:       false,
	strDumpFile:     "",
	dumpFileSno:     0,
	dumpFileSize:    LOG_DEFAULT_DUMP_SIZE,
	dumpThreshold:   LOG_DEFAULT_DUMP_THRESHOLD,
	dumpIntervalMs:  LOG_DEFAULT_DUMP_INTV,
	// queLogs:         make(chan string, MAX_LOG_CACHE_SIZE),
	lck:           &sync.Mutex{},
	queLogs:       nil,
	writeLogs:     nil,
	evtDumpToFile: NewEvent(),
	evtStop:       NewEvent(),
	evtStopSucc:   NewEvent(),
}

func (l *logger) SetLevel(lv LogLv) {
	l.level = lv
}

func (l *logger) SetPowerShellMode() {
	l.bPowerShellMode = true
}

func (l *logger) SetPrintFunc(printFunc func(lv LogLv, logStr string)) {
	l.printFunc = printFunc
}

func (l *logger) D(tag string, a ...interface{}) {
	bExist, _ := IsFileExist(LOG_DEBUG_SWITCH_FILE)
	if !bExist && l.level > LOG_LV_DEBUG {
		return
	}

	l.doLog(LOG_LV_DEBUG, "DEBUG", tag, a...)
}

func (l *logger) I(tag string, a ...interface{}) {
	if l.level > LOG_LV_INFO {
		return
	}

	l.doLog(LOG_LV_INFO, "INFO ", tag, a...)
}

func (l *logger) W(tag string, a ...interface{}) {
	if l.level > LOG_LV_WARN {
		return
	}

	l.doLog(LOG_LV_WARN, "WARN ", tag, a...)
}

func (l *logger) E(tag string, a ...interface{}) {
	l.doLog(LOG_LV_ERROR, "ERROR", tag, a...)
}

func (l *logger) Detail(lv LogLv, s string) {
	l.printLog(lv, s+"\n")
}

func (l *logger) Ln() {
	l.printLog(LOG_LV_INFO, "\n")
}

func (l *logger) doLog(lv LogLv, lvStr string, tag string, a ...interface{}) {
	now := time.Now()
	// timeStr := GetFullTimeString(now, "[%s/%s/%s %s:%s:%s]")
	msg := fmt.Sprint(a...)

	// logStr := ""
	// if !l.bDumpOpen {
	// 	logStr = fmt.Sprint(timeStr, " ", lvStr, " ["+tag+"]  ", msg, "\n")
	// } else {
	// 	logStr = fmt.Sprintln(timeStr, lvStr, "["+tag+"] ", msg)
	// }

	logStr := l.buildLogStr(now, lvStr, tag, msg)
	l.printLog(lv, logStr)
}

func (l *logger) buildLogStr(t time.Time, lvStr string, tag string, msg string) string {
	builder := &strings.Builder{}

	// time
	builder.WriteRune('[')
	FormatTimeStr("YY/MM/DD hh:mm:ss", t, builder)
	builder.WriteRune(']')
	builder.WriteRune(' ')

	// level
	builder.WriteRune('[')
	builder.WriteString(lvStr)
	builder.WriteRune(']')
	builder.WriteRune(' ')

	// tag
	builder.WriteRune('[')
	builder.WriteString(tag)
	builder.WriteRune(']')
	builder.WriteRune(' ')
	builder.WriteRune(' ')

	// msg
	builder.WriteString(msg)
	builder.WriteRune('\n')

	return builder.String()
}

func (l *logger) printLog(lv LogLv, logStr string) {
	l.pushLog(lv, logStr)

	if l.bDumpOpen && l.needDump() {
		l.evtDumpToFile.Send()
	}
}

func (l *logger) pushLog(lv LogLv, log string) {
	l.lck.Lock()
	defer l.lck.Unlock()

	info := &LogInfo{
		Lv:     lv,
		LogStr: log,
	}
	l.queLogs = append(l.queLogs, info)
}

func (l *logger) popLogs() {
	l.lck.Lock()
	defer l.lck.Unlock()

	l.queLogs, l.writeLogs = l.writeLogs, l.queLogs
	// logs := make([]string, len(l.queLogs))
	// copy(logs, l.queLogs)
	// l.queLogs = l.queLogs[0:0]
	// return logs
}

func (l *logger) loop() {
	for {
		bEnd := false
		if !l.bDumpOpen {
			bEnd = l.isStop()
			l.printConsoleLogs()
		} else {
			l.evtDumpToFile.WaitUntilTimeout(l.dumpIntervalMs)
			bEnd = l.isStop() // judge end first, ensure dump all logs before stop dump
			l.dump()
		}

		if bEnd {
			l.evtStopSucc.Send()
			break
		}
	}
}

func (l *logger) stop() {
	l.evtStop.Send()
	l.evtDumpToFile.Send()
	l.evtStopSucc.Wait()
}

func (l *logger) isStop() bool {
	bEnd := false

	select {
	case <-l.evtStop.C:
		bEnd = true

	default:
	}

	return bEnd
}

func (l *logger) printConsoleLogs() {
	l.popLogs()
	if len(l.writeLogs) == 0 {
		<-time.After(time.Millisecond * 10)
		return
	}

	for _, info := range l.writeLogs {
		if l.printFunc != nil {
			l.printFunc(info.Lv, info.LogStr)
		} else {
			l.linuxPrint(info.Lv, info.LogStr)
		}
	}

	l.writeLogs = l.writeLogs[0:0]
}

func (l *logger) linuxPrint(lv LogLv, logStr string) {
	logPrintStr := ""
	if lv == LOG_LV_ERROR {
		logPrintStr = fmt.Sprintf("%c[1;40;31m%s%c[0m", 0x1B, logStr, 0x1B)
	} else if lv == LOG_LV_WARN {
		logPrintStr = fmt.Sprintf("%c[1;40;33m%s%c[0m", 0x1B, logStr, 0x1B)
	} else {
		logPrintStr = logStr
	}
	fmt.Print(logPrintStr)
}

func (l *logger) startDump(file string, dumpFileSize int, dumpThreshold int, dumpIntervalMs uint32) {
	l.strDumpFile = file
	l.dumpFileSize = dumpFileSize
	l.dumpThreshold = dumpThreshold
	l.dumpIntervalMs = dumpIntervalMs
	l.queLogs = make([]*LogInfo, 0, LOG_MAX_CACHE_SIZE)
	l.writeLogs = make([]*LogInfo, 0, LOG_MAX_CACHE_SIZE)
	l.bDumpOpen = true
	// go l.dumpLoop()
}

func (l *logger) stopDump() {
	l.bDumpOpen = false
	// l.evtStop.Send()
	l.evtDumpToFile.Send()
	// l.evtStopSucc.Wait()
	// l.strDumpFile = ""
}

func (l *logger) needDump() bool {
	return len(l.queLogs) >= l.dumpThreshold
}

func (l *logger) dump() {
	l.popLogs()
	if len(l.writeLogs) == 0 {
		return
	}

	cnt, err := l.dumpToFile(l.writeLogs)
	if err != nil {
		l.dumpToBak(l.writeLogs[cnt:])
	}

	l.writeLogs = l.writeLogs[0:0]
}

func (l *logger) dumpToFile(logs []*LogInfo) (int, error) {
	var err error = nil
	totalCnt := len(logs)
	idx := 0
	cnt := int(0)

	for {
		// dump one file
		bNeedRename := false
		cnt, bNeedRename, err = l.dumpOneFile(logs[idx:])
		idx += cnt
		if err != nil {
			break
		}

		// rename
		if bNeedRename {
			renameErr := l.renameDumpFile()
			if renameErr != nil {
				fmt.Println("rename dump file error: ", renameErr)
			}
		}

		// check end
		if idx == totalCnt {
			break
		}
	}

	return idx, err
}

func (l *logger) dumpOneFile(logs []*LogInfo) (int, bool, error) {
	var err error = nil

	// open file
	fileName := l.strDumpFile
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open log dump file error: ", err)
		return 0, false, err
	}

	defer f.Close()

	// dump loop
	totalCnt := len(logs)
	idx := 0
	cnt := int(0)
	bNeedRename := false

	for {
		// batch size
		batchSize := totalCnt - idx
		if batchSize > LOG_BATCH_DUMP_COUNT {
			batchSize = LOG_BATCH_DUMP_COUNT
		}

		// dump
		cnt, err = l.batchDumpToFile(logs[idx:idx+batchSize], f)
		idx += cnt
		if err != nil {
			break
		}

		// check file size
		size, sizeErr := GetFileSize(fileName)
		if sizeErr != nil {
			fmt.Println("GetFileSize error: ", sizeErr)
		} else if size >= int64(l.dumpFileSize) {
			bNeedRename = true
			break
		}

		// check end
		if idx == totalCnt {
			break
		}
	}

	return idx, bNeedRename, err
}

func (l *logger) dumpToBak(logs []*LogInfo) {
	f, err := os.OpenFile("dump.log.bak", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open dump.log.bak error: ", err)
		return
	}

	defer f.Close()

	l.batchDumpToFile(logs, f)
}

func (l *logger) batchDumpToFile(logs []*LogInfo, f *os.File) (int, error) {
	w := bufio.NewWriter(f)
	defer w.Flush()

	loopCnt := len(logs)
	for i := 0; i < loopCnt; i++ {
		_, err := w.WriteString(logs[i].LogStr)
		if err != nil {
			fmt.Println("batchDumpToFile w.WriteString error: ", err)
			return i, err
		}
	}

	return loopCnt, nil
}

func (l *logger) renameDumpFile() error {
	l.dumpFileSno++
	dir := path.Dir(l.strDumpFile)
	name := path.Base(l.strDumpFile)
	ext := path.Ext(name)
	nameOnly := strings.TrimSuffix(name, ext)

	builder := &strings.Builder{}
	builder.WriteString(nameOnly)

	now := time.Now()
	FormatTimeStr("_YYMMDD_hhmmss_", now, builder)
	FormatUint(l.dumpFileSno, 5, false, builder)
	builder.WriteString(ext)
	newName := path.Join(dir, builder.String())

	return os.Rename(l.strDumpFile, newName)
}
