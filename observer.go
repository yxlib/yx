// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import "errors"

var (
	ErrObsMsgEmpty = errors.New("msgName is empty")
)

type NotifyMsg struct {
	Name   string
	Params []interface{}
}

type Observer interface {
	// Callback when dispatcher notify message.
	// @param msg, the message notifying.
	OnNotify(msg *NotifyMsg)
}

type Dispatcher interface {
	// Add an observer to the dispatcher.
	// @param o, the observer to be added.
	AddObserver(o Observer)

	// Remove an observer from the dispatcher.
	// @param o, the observer to be removed.
	RemoveObserver(o Observer)

	// Notify message.
	// @param params, the params of the message.
	Notify(params ...interface{})
}

//=====================================================
//                    BaseDispatcher
//=====================================================
type BaseDispatcher struct {
	msgName     string
	observers   *ObjectSet
	lckObserver *FastLock
}

func NewBaseDispatcher(msgName string) *BaseDispatcher {
	return &BaseDispatcher{
		msgName:     msgName,
		observers:   NewObjectSet(),
		lckObserver: NewFastLock(),
	}
}

func (d *BaseDispatcher) GetMsgName() string {
	return d.msgName
}

func (d *BaseDispatcher) AddObserver(o Observer) {
	if o == nil {
		return
	}

	// d.lckObserver.Lock()
	if d.lckObserver.TryLock(0) != nil {
		return
	}

	defer d.lckObserver.Unlock()

	d.observers.Add(o)

	// for _, obs := range d.observers {
	// 	if o == obs {
	// 		return
	// 	}
	// }

	// d.observers = append(d.observers, o)
}

func (d *BaseDispatcher) RemoveObserver(o Observer) {
	if o == nil {
		return
	}

	// d.lckObserver.Lock()
	if d.lckObserver.TryLock(0) != nil {
		return
	}

	defer d.lckObserver.Unlock()

	d.observers.Remove(o)
	// for i, obs := range d.observers {
	// 	if o == obs {
	// 		d.observers = append(d.observers[:i], d.observers[i+1:]...)
	// 		return
	// 	}
	// }
}

func (d *BaseDispatcher) Notify(params ...interface{}) {

}

func (d *BaseDispatcher) notifyImpl(params ...interface{}) {
	observers := d.cloneObservers()
	for _, obj := range observers {
		obs := obj.(Observer)
		bRm, err := d.isObserverRemove(obs)
		if err != nil || bRm {
			continue
		}

		msgParams := make([]interface{}, 0, len(params))
		msgParams = append(msgParams, params...)

		msg := &NotifyMsg{
			Name:   d.msgName,
			Params: msgParams,
		}

		obs.OnNotify(msg)
	}
}

func (d *BaseDispatcher) cloneObservers() []interface{} {
	// d.lckObserver.RLock()
	if d.lckObserver.TryLock(0) != nil {
		return []interface{}{}
	}

	defer d.lckObserver.Unlock()

	return d.observers.GetElements()

	// observers := make([]Observer, len(d.observers))
	// copy(observers, d.observers)
	// return observers
}

func (d *BaseDispatcher) isObserverRemove(o Observer) (bool, error) {
	// d.lckObserver.RLock()
	if err := d.lckObserver.TryLock(0); err != nil {
		return false, err
	}

	defer d.lckObserver.Unlock()

	bExist, err := d.observers.Exist(o)
	return !bExist, err

	// for _, obs := range d.observers {
	// 	if o == obs {
	// 		return false
	// 	}
	// }

	// return true
}

//=====================================================
//                    SyncDispatcher
//=====================================================
type SyncDispatcher struct {
	*BaseDispatcher
}

func NewSyncDispatcher(msgName string) *SyncDispatcher {
	return &SyncDispatcher{
		BaseDispatcher: NewBaseDispatcher(msgName),
	}
}

func (d *SyncDispatcher) Notify(params ...interface{}) {
	d.notifyImpl(params...)
}

//=====================================================
//                    AsyncDispatcher
//=====================================================
type AsyncDispatcher struct {
	*BaseDispatcher
	chanMsg chan *NotifyMsg
	evtStop *Event
	evtExit *Event
}

func NewAsyncDispatcher(msgName string, maxMsgBuffSize uint16) *AsyncDispatcher {
	return &AsyncDispatcher{
		BaseDispatcher: NewBaseDispatcher(msgName),
		chanMsg:        make(chan *NotifyMsg, maxMsgBuffSize),
		evtStop:        NewEvent(),
		evtExit:        NewEvent(),
	}
}

func (d *AsyncDispatcher) Notify(params ...interface{}) {
	msgParams := make([]interface{}, 0, len(params))
	msgParams = append(msgParams, params...)

	msg := &NotifyMsg{
		Name:   d.msgName,
		Params: msgParams,
	}

	d.chanMsg <- msg
}

func (d *AsyncDispatcher) Start() {
	ch := d.evtStop.GetChan()
	for {
		select {
		case msg := <-d.chanMsg:
			d.notifyImpl(msg.Params...)

		case <-ch:
			goto Exit0
		}
	}

Exit0:
	d.evtExit.Close()
}

func (d *AsyncDispatcher) Stop() {
	d.evtStop.Close()
	d.evtExit.Wait()
}

//=====================================================
//                    NotifyCenter
//=====================================================
type NotifyCenter struct {
	mapName2Dispatcher map[string]Dispatcher
	lckDispatcher      *FastLock
}

func NewNotifyCenter() *NotifyCenter {
	return &NotifyCenter{
		mapName2Dispatcher: make(map[string]Dispatcher),
		lckDispatcher:      NewFastLock(),
	}
}

// Add an observer to the notify center.
// @param msgName, the name of message which the observer will listen to.
// @param o, the observer to be added.
func (c *NotifyCenter) AddObserver(msgName string, o Observer) {
	// if len(msgName) == 0 {
	// 	return
	// }

	if o == nil {
		return
	}

	// c.lckDispatcher.Lock()
	// defer c.lckDispatcher.Unlock()

	// dispatcher, ok := c.mapName2Dispatcher[msgName]
	// if !ok {
	// 	dispatcher = NewBaseDispatcher(msgName)
	// 	c.mapName2Dispatcher[msgName] = dispatcher
	// }

	dispatcher, err := c.confirmDispatcherExist(msgName)
	if err != nil {
		return
	}

	dispatcher.AddObserver(o)
}

// Remove an observer from the notify center.
// @param msgName, the name of message which the observer listening to.
// @param o, the observer to be removed.
func (c *NotifyCenter) RemoveObserver(msgName string, o Observer) {
	// if len(msgName) == 0 {
	// 	return
	// }

	if o == nil {
		return
	}

	// c.lckDispatcher.Lock()
	// defer c.lckDispatcher.Unlock()

	// dispatcher, ok := c.mapName2Dispatcher[msgName]
	dispatcher, ok := c.GetDispatcher(msgName)
	if !ok {
		return
	}

	dispatcher.RemoveObserver(o)
}

// Notify message.
// @param msgName, the name of message.
// @param params, the params of the message.
func (c *NotifyCenter) Notify(msgName string, params ...interface{}) {
	// if len(msgName) == 0 {
	// 	return
	// }

	// c.lckDispatcher.Lock()
	// defer c.lckDispatcher.Unlock()

	// dispatcher, ok := c.mapName2Dispatcher[msgName]
	dispatcher, ok := c.GetDispatcher(msgName)
	if !ok {
		return
	}

	dispatcher.Notify(params...)
}

// Add an dispatcher to the notify center.
// @param msgName, the name of message which the dispatcher will notify.
// @param dispatcher, the dispatcher to be added.
func (c *NotifyCenter) AddDispatcher(msgName string, dispatcher Dispatcher) {
	if len(msgName) == 0 {
		return
	}

	if dispatcher == nil {
		return
	}

	// c.lckDispatcher.Lock()
	if c.lckDispatcher.TryLock(0) != nil {
		return
	}

	defer c.lckDispatcher.Unlock()

	c.mapName2Dispatcher[msgName] = dispatcher
}

// Remove an dispatcher from the notify center.
// @param msgName, the name of message which the dispatcher is notify.
func (c *NotifyCenter) RemoveDispatcher(msgName string) {
	if len(msgName) == 0 {
		return
	}

	// c.lckDispatcher.Lock()
	if c.lckDispatcher.TryLock(0) != nil {
		return
	}

	defer c.lckDispatcher.Unlock()

	_, ok := c.mapName2Dispatcher[msgName]
	if ok {
		delete(c.mapName2Dispatcher, msgName)
	}
}

// Get an dispatcher from the notify center.
// @param msgName, the name of message which the dispatcher is notify.
// @return Dispatcher, the dispatcher
// @return bool, true mean success, false mean failed
func (c *NotifyCenter) GetDispatcher(msgName string) (Dispatcher, bool) {
	if len(msgName) == 0 {
		return nil, false
	}

	if c.lckDispatcher.TryLock(0) != nil {
		return nil, false
	}

	defer c.lckDispatcher.Unlock()

	d, ok := c.mapName2Dispatcher[msgName]
	return d, ok
}

func (c *NotifyCenter) confirmDispatcherExist(msgName string) (Dispatcher, error) {
	if len(msgName) == 0 {
		return nil, ErrObsMsgEmpty
	}

	// c.lckDispatcher.Lock()
	if err := c.lckDispatcher.TryLock(0); err != nil {
		return nil, err
	}

	defer c.lckDispatcher.Unlock()

	dispatcher, ok := c.mapName2Dispatcher[msgName]
	if !ok {
		dispatcher = NewBaseDispatcher(msgName)
		c.mapName2Dispatcher[msgName] = dispatcher
	}

	return dispatcher, nil
}
