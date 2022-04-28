// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	ErrUtilObjIsNil = errors.New("object is nil")
)

// Get the string of full time.
// @param t, the time object to format.
// @param format, format of the string.
// @return the format string of full time.
func GetFullTimeString(t time.Time, format string) string {
	// t := time.Now()

	yy := t.Year()
	yyStr := strconv.Itoa(yy)

	mm := int(t.Month())
	mmStr := strconv.Itoa(mm)
	if mm < 10 {
		mmStr = "0" + mmStr
	}

	dd := t.Day()
	ddStr := strconv.Itoa(dd)
	if dd < 10 {
		ddStr = "0" + ddStr
	}

	h := t.Hour()
	hStr := strconv.Itoa(h)
	if h < 10 {
		hStr = "0" + hStr
	}

	m := t.Minute()
	mStr := strconv.Itoa(m)
	if m < 10 {
		mStr = "0" + mStr
	}

	s := t.Second()
	sStr := strconv.Itoa(s)
	if s < 10 {
		sStr = "0" + sStr
	}

	return fmt.Sprintf(format, yyStr, mmStr, ddStr, hStr, mStr, sStr)
}

// Get the string of date.
// @param t, the time object to format.
// @param format, format of the string.
// @return the format string of date.
func GetDateString(t time.Time, format string) string {
	// t := time.Now()

	yy := t.Year()
	yyStr := strconv.Itoa(yy)

	mm := int(t.Month())
	mmStr := strconv.Itoa(mm)
	if mm < 10 {
		mmStr = "0" + mmStr
	}

	dd := t.Day()
	ddStr := strconv.Itoa(dd)
	if dd < 10 {
		ddStr = "0" + ddStr
	}

	return fmt.Sprintf(format, yyStr, mmStr, ddStr)
}

// Get the string of time.
// @param t, the time object to format.
// @param format, format of the string.
// @return the format string of time.
func GetTimeString(t time.Time, format string) string {
	// t := time.Now()

	h := t.Hour()
	hStr := strconv.Itoa(h)
	if h < 10 {
		hStr = "0" + hStr
	}

	m := t.Minute()
	mStr := strconv.Itoa(m)
	if m < 10 {
		mStr = "0" + mStr
	}

	s := t.Second()
	sStr := strconv.Itoa(s)
	if s < 10 {
		sStr = "0" + sStr
	}

	return fmt.Sprintf(format, hStr, mStr, sStr)
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
	path := t.PkgPath()
	name := path + "." + t.Name()
	return name, nil
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
