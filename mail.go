// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"errors"

	"gopkg.in/gomail.v2"
)

var (
	ErrMailIsNil = errors.New("mail is nil")
)

type Mail struct {
	User      string
	AliasName string
	Pwd       string
	Host      string
	Port      int
	To        []string
	Subject   string
	Body      string
}

// Send a mail.
// @param m, the mail to send.
// @return error, error.
func SendMail(m *Mail) error {
	if m == nil {
		return ErrMailIsNil
	}

	msg := gomail.NewMessage()

	if len(m.AliasName) != 0 {
		msg.SetHeader("From", msg.FormatAddress(m.User, m.AliasName))
	} else {
		msg.SetHeader("From", m.User)
	}

	msg.SetHeader("To", m.To...)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.Body)

	d := gomail.NewDialer(m.Host, m.Port, m.User, m.Pwd)
	err := d.DialAndSend(msg)
	return err
}
