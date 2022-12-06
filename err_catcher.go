// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"fmt"
	"sync"
)

type ErrCatcher struct {
	className string
}

func NewErrCatcher(className string) *ErrCatcher {
	return &ErrCatcher{
		className: className,
	}
}

// Catch error.
// @param methodName, the name of method which catch the error.
// @param errRef, reference of an error.
func (c *ErrCatcher) Catch(methodName string, errRef *error) {
	if errRef == nil || *errRef == nil {
		return
	}

	errCatcherInst.CatchError(c.className, methodName, *errRef)
}

// Throw error.
// @param methodName, the name of method which throw the error.
// @param err, an error.
// @return error, return the same error of param err.
func (c *ErrCatcher) Throw(methodName string, err error) error {
	if err == nil {
		return nil
	}

	errCatcherInst.ThrowError(c.className, methodName, err)
	return err
}

func (c *ErrCatcher) TryCodeFunc(methodName string, f func() (int32, error)) (int32, error) {
	code, err := f()
	if err != nil {
		c.Throw(methodName, err)
	}

	return code, err
}

func (c *ErrCatcher) TryFunc(methodName string, f func() error) error {
	err := f()
	if err != nil {
		c.Throw(methodName, err)
	}

	return err
}

// Throw error with defer.
// @param methodName, the name of method which throw the error.
// @param errRef, reference of an error.
func (c *ErrCatcher) DeferThrow(methodName string, errRef *error) {
	if errRef == nil || *errRef == nil {
		return
	}

	errCatcherInst.ThrowError(c.className, methodName, *errRef)
}

type InvokeInfo struct {
	className  string
	methodName string
}

func NewInvokeInfo(className string, methodName string) *InvokeInfo {
	return &InvokeInfo{
		className:  className,
		methodName: methodName,
	}
}

type InvokeStack = []*InvokeInfo

type errCatcher struct {
	mapErr2InvokeStack map[error]InvokeStack
	lck                *sync.Mutex
	// logger             *Logger
}

var errCatcherInst = &errCatcher{
	mapErr2InvokeStack: make(map[error]InvokeStack),
	lck:                &sync.Mutex{},
	// logger:             NewLogger("ErrCatcher"),
}

func (c *errCatcher) ThrowError(className string, methodName string, err error) {
	if err == nil {
		return
	}

	c.pushError(className, methodName, err)
}

func (c *errCatcher) CatchError(className string, methodName string, err error) {
	if err == nil {
		return
	}

	loggerInst.E("ErrCatcher", "Catch Error !!!")

	logs := make([]string, 0)
	logs = c.beginPrintError(err, logs)

	stack, ok := c.popError(err)
	if ok {
		logs = c.printInvokeStack(className, methodName, stack, logs)
	} else {
		log := fmt.Sprint("[S] ", className, ".", methodName, "()\n")
		logs = append(logs, log)
	}

	logs = c.endPrintError(logs)
	loggerInst.Detail(LOG_LV_ERROR, logs)
}

func (c *errCatcher) printInvokeStack(className string, methodName string, stack []*InvokeInfo, logs []string) []string {
	log := ""
	blanks := " "
	mark := "|__ "
	for i, info := range stack {
		if i == 0 {
			log = fmt.Sprint("[S]", blanks, info.className, ".", info.methodName, "()\n")
		} else {
			log = fmt.Sprint("[S]", blanks, mark, info.className, ".", info.methodName, "()\n")
		}

		logs = append(logs, log)
		blanks += "  "
	}

	log = fmt.Sprint("[S]", blanks, mark, className, ".", methodName, "()\n")
	logs = append(logs, log)
	return logs
}

func (c *errCatcher) beginPrintError(err error, logs []string) []string {
	logs = append(logs, "[E] ====================================================\n")

	log := fmt.Sprint("[M] ** ERROR: ", err.Error(), " **", "\n")
	logs = append(logs, log)

	logs = append(logs, "\n")
	return logs
}

func (c *errCatcher) endPrintError(logs []string) []string {
	logs = append(logs, "[S]\n")
	logs = append(logs, "[E] ====================================================\n")
	logs = append(logs, "\n")
	return logs
}

func (c *errCatcher) pushError(className string, methodName string, err error) {
	c.lck.Lock()
	defer c.lck.Unlock()

	stack, ok := c.mapErr2InvokeStack[err]
	if !ok {
		stack = make(InvokeStack, 0)
	}

	info := NewInvokeInfo(className, methodName)
	c.mapErr2InvokeStack[err] = append(stack, info)
}

func (c *errCatcher) popError(err error) (InvokeStack, bool) {
	c.lck.Lock()
	defer c.lck.Unlock()

	stack, ok := c.mapErr2InvokeStack[err]
	if ok {
		delete(c.mapErr2InvokeStack, err)
	}

	return stack, ok
}
