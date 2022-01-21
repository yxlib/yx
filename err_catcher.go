// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import "sync"

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

	blanks := " "
	stack, ok := c.popError(err)
	if !ok {
		c.beginPrintError(err)
		loggerInst.Detail(LOG_LV_ERROR, "[S]"+blanks+className+"."+methodName+"()")
		c.endPrintError()
		return
	}

	mark := "|__ "
	for i, info := range stack {
		if i == 0 {
			c.beginPrintError(err)
			loggerInst.Detail(LOG_LV_ERROR, "[S]"+blanks+info.className+"."+info.methodName+"()")
		} else {
			loggerInst.Detail(LOG_LV_ERROR, "[S]"+blanks+mark+info.className+"."+info.methodName+"()")
		}

		blanks += "  "
	}

	loggerInst.Detail(LOG_LV_ERROR, "[S]"+blanks+mark+className+"."+methodName+"()")
	c.endPrintError()
}

func (c *errCatcher) beginPrintError(err error) {
	loggerInst.E("ErrCatcher", "Catch Error !!!")
	loggerInst.Detail(LOG_LV_ERROR, "[E] ====================================================")
	loggerInst.Detail(LOG_LV_ERROR, "[M] ** ERROR: "+err.Error()+" **")
	loggerInst.Ln()
}

func (c *errCatcher) endPrintError() {
	loggerInst.Detail(LOG_LV_ERROR, "[S]")
	loggerInst.Detail(LOG_LV_ERROR, "[E] ====================================================")
	loggerInst.Ln()
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
