// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var (
	ErrUtilObjIsNil    = errors.New("object is nil")
	ErrWrongTimeFormat = errors.New("wrong time format")
)

// Format the string of time.
// @param format, format of the string with YY MM DD hh mm ss.
// @param t, the time object to format.
// @param builder, a string builder. if not nil, only build string, not return
// @return string, the format string of time.
// @return error, the error
func FormatTimeStr(format string, t time.Time, builder *strings.Builder) (string, error) {
	bNeedReturnString := (builder == nil)
	if bNeedReturnString {
		builder = &strings.Builder{}
	}

	startIdx := -1
	YYIdx := strings.Index(format, "YY")
	if YYIdx >= 0 {
		startIdx = YYIdx
	}

	MMIdx := strings.Index(format, "MM")
	if MMIdx >= 0 && startIdx == -1 {
		startIdx = MMIdx
	}

	DDIdx := strings.Index(format, "DD")
	if DDIdx >= 0 && startIdx == -1 {
		startIdx = DDIdx
	}

	hhIdx := strings.Index(format, "hh")
	if hhIdx >= 0 && startIdx == -1 {
		startIdx = hhIdx
	}

	mmIdx := strings.Index(format, "mm")
	if mmIdx >= 0 && startIdx == -1 {
		startIdx = mmIdx
	}

	ssIdx := strings.Index(format, "ss")
	if ssIdx >= 0 && startIdx == -1 {
		startIdx = ssIdx
	}

	if startIdx < 0 {
		return "", ErrWrongTimeFormat
	}

	if startIdx > 0 {
		builder.WriteString(format[:startIdx])
	}

	// year
	if YYIdx >= 0 {
		yy := t.Year()
		FormatInt(int64(yy), 4, false, builder)
		startIdx = YYIdx + 2
	}

	// month
	if MMIdx >= 0 {
		if startIdx < MMIdx {
			builder.WriteString(format[startIdx:MMIdx])
		}

		mm := int(t.Month())
		FormatInt(int64(mm), 2, true, builder)
		startIdx = MMIdx + 2
	}

	// day
	if DDIdx >= 0 {
		if startIdx < DDIdx {
			builder.WriteString(format[startIdx:DDIdx])
		}

		dd := t.Day()
		FormatInt(int64(dd), 2, true, builder)
		startIdx = DDIdx + 2
	}

	// hour
	if hhIdx >= 0 {
		if startIdx < hhIdx {
			builder.WriteString(format[startIdx:hhIdx])
		}

		h := t.Hour()
		FormatInt(int64(h), 2, true, builder)
		startIdx = hhIdx + 2
	}

	// minute
	if mmIdx >= 0 {
		if startIdx < mmIdx {
			builder.WriteString(format[startIdx:mmIdx])
		}

		m := t.Minute()
		FormatInt(int64(m), 2, true, builder)
		startIdx = mmIdx + 2
	}

	// second
	if ssIdx >= 0 {
		if startIdx < ssIdx {
			builder.WriteString(format[startIdx:ssIdx])
		}

		s := t.Second()
		FormatInt(int64(s), 2, true, builder)
		startIdx = ssIdx + 2
	}

	if startIdx < len(format) {
		builder.WriteString(format[startIdx:])
	}

	if bNeedReturnString {
		return builder.String(), nil
	}

	return "", nil
}

// Format integer.
// @param num, the number to format.
// @param maxLength, the max length in decimal.
// @param bFillZero, is fill zero to prefix.
// @param builder, a string builder. if not nil, only build string, not return
// @return string, the format string of time.
func FormatInt(num int64, maxLength uint32, bFillZero bool, builder *strings.Builder) string {
	bNeedReturnString := (builder == nil)
	if bNeedReturnString {
		builder = &strings.Builder{}
	}

	absNum := uint64(num)
	if num < 0 {
		builder.WriteRune('-')
		absNum = uint64(math.Abs(float64(num)))
	}

	return FormatUint(absNum, maxLength, bFillZero, builder)
}

// Format unsigned integer.
// @param num, the number to format.
// @param maxLength, the max length in decimal.
// @param bFillZero, is fill zero to prefix.
// @param builder, a string builder. if not nil, only build string, not return
// @return string, the format string of time.
func FormatUint(num uint64, maxLength uint32, bFillZero bool, builder *strings.Builder) string {
	bNeedReturnString := (builder == nil)
	if bNeedReturnString {
		builder = &strings.Builder{}
	}

	bStart := false
	for i := int(maxLength) - 1; i >= 0; i-- {
		divisor := PowerOfTen(uint32(i))
		quotient := num / divisor
		if quotient > 0 || bStart || bFillZero {
			builder.WriteRune(rune(0x30 + quotient))
		}

		if !bStart {
			bStart = (quotient > 0)
		}

		num = num % divisor
	}

	if bNeedReturnString {
		return builder.String()
	}

	return ""
}

// Get the power of ten.
// @param power, the power.
// @return uint64, the result
func PowerOfTen(power uint32) uint64 {
	result := uint64(1)
	for i := 0; i < int(power); i++ {
		result *= 10
	}

	return result
}

// Protect run, if panic, it will recover.
// @param entry, the function of danger code.
func ProtectRun(entry func()) {
	if entry == nil {
		return
	}

	defer func() {
		err := recover()
		if err != nil {
			log := NewLogger("panic")

			switch err.(type) {
			case runtime.Error:
				log.E("runtime error:", err)
			default:
				log.E("error:", err)
			}
		}
	}()

	entry()
}

// Run danger code.
// @param entry, the function of danger code.
// @param bDebugMode, true will run in normal mode, false will run in protect mode.
func RunDangerCode(entry func(), bDebugMode bool) {
	if entry == nil {
		return
	}

	if !bDebugMode {
		ProtectRun(entry)
	} else {
		entry()
	}
}

// Get the size of file.
// @param path, the file path.
// @return int64, the size.
// @return error, error.
func GetFileSize(path string) (int64, error) {
	fs, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return fs.Size(), nil
}

// Is the file exist.
// @param path, the file path.
// @return bool, if error is nil, true mean exist, false mean not exist.
// @return error, error.
func IsFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Load config from a json file.
// @param v, config object.
// @param path, path of the json file.
// @param decodeCb, a callback function to decode the content of the file.
// @return error, error.
func LoadJsonConf(v interface{}, path string, decodeCb func(data []byte) ([]byte, error)) error {
	// read file
	filePtr, err := os.Open(path)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(filePtr)
	if err != nil {
		return err
	}

	// decode
	if decodeCb != nil {
		d, err = decodeCb(d)
		if err != nil {
			return err
		}
	}

	// json unmarshal
	err = json.Unmarshal(d, v)
	return err
}

func GetClassReflectName(obj interface{}) (string, error) {
	if obj == nil {
		return "", ErrUtilObjIsNil
	}

	t := reflect.TypeOf(obj)
	t = t.Elem()
	name := GetReflectNameByType(t)
	return name, nil
}

func GetReflectNameByType(t reflect.Type) string {
	path := t.PkgPath()
	name := path + "." + t.Name()
	return name
}

func GetFullPackageName(classReflectName string) string {
	idx := strings.LastIndex(classReflectName, ".")
	if idx < 0 {
		return ""
	}

	packName := classReflectName[:idx]
	return packName
}

func GetFilePackageName(fullPackName string) string {
	idx := strings.LastIndex(fullPackName, "/")
	if idx < 0 {
		return fullPackName
	}

	return fullPackName[idx+1:]
}

func GetClassName(classReflectName string) string {
	idx := strings.LastIndex(classReflectName, ".")
	if idx < 0 {
		return classReflectName
	}

	return classReflectName[idx+1:]
}

func GetFilePackageClassName(classReflectName string) string {
	idx := strings.LastIndex(classReflectName, "/")
	if idx < 0 {
		return classReflectName
	}

	return classReflectName[idx+1:]
}

func Daemon(program string, args []string, restartDelay uint16, shutdownFile string) error {
	for {
		cmd := exec.Command(program, args...)
		err := cmd.Start()
		if err != nil {
			return err
		}

		err = cmd.Wait()
		if err != nil {
			return err
		}

		ok, _ := IsFileExist(shutdownFile)
		if ok {
			break
		}

		if restartDelay > 0 {
			t := time.After(time.Duration(restartDelay) * time.Second)
			<-t
		}
	}

	return nil
}
