// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import "errors"

var (
	ErrNoSpaceToRead = errors.New("no space to read data")
	ErrBuffEmpty     = errors.New("buffer empty")
	ErrReachEnd      = errors.New("buffer reach end")
)

type SimpleBuffer struct {
	dataBuff    []byte
	writeOffset uint32
	readOffset  uint32
}

func NewSimpleBuffer(capacity int) *SimpleBuffer {
	return &SimpleBuffer{
		dataBuff:    make([]byte, capacity),
		writeOffset: 0,
		readOffset:  0,
	}
}

// Get the size of write buffer.
// @return int, the size of write buffer.
func (buff *SimpleBuffer) GetWriteBuffSize() int {
	return len(buff.dataBuff) - int(buff.writeOffset)
}

// Get the write buffer.
// @return []byte, the write buffer.
func (buff *SimpleBuffer) GetWriteBuff() []byte {
	return buff.dataBuff[buff.writeOffset:]
}

// Get the length of data.
// @return int, the length of data.
func (buff *SimpleBuffer) GetDataLen() int {
	if buff.isNoData() || buff.isReachEnd() {
		return 0
	}

	return int(buff.writeOffset - buff.readOffset)
}

// Get the data.
// @return []byte, the data.
func (buff *SimpleBuffer) GetData() []byte {
	return buff.dataBuff[buff.readOffset:buff.writeOffset]
}

// Move the data to offset = 0.
func (buff *SimpleBuffer) MoveDataToBegin() {
	if buff.readOffset > 0 {
		copy(buff.dataBuff, buff.dataBuff[buff.readOffset:buff.writeOffset])
		buff.writeOffset = buff.writeOffset - buff.readOffset
		buff.readOffset = 0
	}
}

// Update the write offset.
// @param addLen, the length to add.
func (buff *SimpleBuffer) UpdateWriteOffset(addLen uint32) {
	if addLen == 0 {
		return
	}

	capacity := len(buff.dataBuff)
	buff.writeOffset += addLen
	if int(buff.writeOffset) > capacity {
		buff.writeOffset = uint32(capacity)
	}
}

// Skip the data.
// @param skipLen, the length to skip.
func (buff *SimpleBuffer) Skip(skipLen uint32) {
	if skipLen == 0 {
		return
	}

	buff.updateReadOffset(skipLen)
}

// Simulate reading data, only get the data, the read offset is not update.
// @param b, the dest buffer.
// @return int, the length of data has read
// @return error, error
func (buff *SimpleBuffer) SimulateRead(b []byte) (n int, err error) {
	return buff.readImpl(b, false)
}

// Read data, get the data, and update the read offset.
// @param b, the dest buffer.
// @return int, the length of data has read
// @return error, error
func (buff *SimpleBuffer) Read(b []byte) (n int, err error) {
	return buff.readImpl(b, true)
}

func (buff *SimpleBuffer) readImpl(b []byte, bUpdateOffset bool) (n int, err error) {
	if len(b) == 0 {
		return 0, ErrNoSpaceToRead
	}

	if buff.isNoData() {
		return 0, ErrBuffEmpty
	}

	if buff.isReachEnd() {
		return 0, ErrReachEnd
	}

	copyLen := copy(b, buff.dataBuff[buff.readOffset:buff.writeOffset])
	if bUpdateOffset {
		buff.updateReadOffset(uint32(copyLen))
	}

	return copyLen, nil
}

func (buff *SimpleBuffer) updateReadOffset(addLen uint32) {
	buff.readOffset += addLen
	if buff.isReachEnd() {
		buff.reset()
	}
}

func (buff *SimpleBuffer) isNoData() bool {
	return buff.writeOffset == 0
}

func (buff *SimpleBuffer) isReachEnd() bool {
	return buff.readOffset >= buff.writeOffset
}

func (buff *SimpleBuffer) reset() {
	buff.writeOffset = 0
	buff.readOffset = 0
}
