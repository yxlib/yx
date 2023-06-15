// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

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
	lck                *FastLock
	// lck                *sync.Mutex
	// logger             *Logger
}

var errCatcherInst = &errCatcher{
	mapErr2InvokeStack: make(map[error]InvokeStack),
	lck:                NewFastLock(),
	// lck:                &sync.Mutex{},
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

	logs := make([][]interface{}, 0)
	logs = c.beginPrintError(err, logs)

	stack, ok := c.popError(err)
	if ok {
		logs = c.printInvokeStack(className, methodName, stack, logs)
	} else {
		logs = append(logs, LogArgs("[S] ", className, ".", methodName, "()"))
	}

	logs = c.endPrintError(logs)
	loggerInst.Detail(LOG_LV_ERROR, logs)
}

func (c *errCatcher) printInvokeStack(className string, methodName string, stack []*InvokeInfo, logs [][]interface{}) [][]interface{} {
	var log []interface{} = nil
	blanks := " "
	mark := "|__ "
	for i, info := range stack {
		if i == 0 {
			log = LogArgs("[S]", blanks, info.className, ".", info.methodName, "()")
		} else {
			log = LogArgs("[S]", blanks, mark, info.className, ".", info.methodName, "()")
		}

		logs = append(logs, log)
		blanks += "  "
	}

	logs = append(logs, LogArgs("[S]", blanks, mark, className, ".", methodName, "()"))
	return logs
}

func (c *errCatcher) beginPrintError(err error, logs [][]interface{}) [][]interface{} {
	logs = append(logs, LogArgs("[E] ===================================================="))
	logs = append(logs, LogArgs("[M] ** ERROR: ", err.Error(), " **"))
	logs = append(logs, LogArgs("[E]"))
	return logs
}

func (c *errCatcher) endPrintError(logs [][]interface{}) [][]interface{} {
	logs = append(logs, LogArgs("[S]"))
	logs = append(logs, LogArgs("[E] ===================================================="))
	return logs
}

func (c *errCatcher) pushError(className string, methodName string, err error) {
	// c.lck.Lock()
	if c.lck.TryLock(0) != nil {
		return
	}

	defer c.lck.Unlock()

	stack, ok := c.mapErr2InvokeStack[err]
	if !ok {
		stack = make(InvokeStack, 0)
	}

	info := NewInvokeInfo(className, methodName)
	c.mapErr2InvokeStack[err] = append(stack, info)
}

func (c *errCatcher) popError(err error) (InvokeStack, bool) {
	// c.lck.Lock()
	if c.lck.TryLock(0) != nil {
		return nil, false
	}

	defer c.lck.Unlock()

	stack, ok := c.mapErr2InvokeStack[err]
	if ok {
		delete(c.mapErr2InvokeStack, err)
	}

	return stack, ok
}
