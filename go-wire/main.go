package main

import (
	"errors"
	"fmt"
	"time"
)

type Message string
type Greeter struct {
	Message Message
}
type Event struct {
	Greeter Greeter
}

// provider function
func NewMessage() (Message, error) {
	if time.Now().Unix()%2 == 0 {
		return "", errors.New("failed to create message")
	}
	return Message("Hello, Dependency Injection in Go!"), nil
}

// NewGreeter: 它的入参不需要变。
// 注意：Wire 极其聪明，它知道 Greeter 需要的是 Message，而不是 error。
// 它会自动处理完 NewMessage 的 error 后，才把 Message 传给这里。
func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}
func (e Event) Start() {
	msg := e.Greeter.Message
	fmt.Println(msg)
}

type IMessage interface {
	Getcontent() string
}
type SimpleMessage struct {
	Content string
}

func (sm *SimpleMessage) Getcontent() string {
	return sm.Content
}

func NewSimpleMessage() (*SimpleMessage, func(), error) {
	msg := &SimpleMessage{Content: "This is a clean message"}
	cleanup := func() {
		fmt.Println("Cleaning up SimpleMessage resources")
	}
	return msg, cleanup, nil
}
func NewGreeter(im IMessage) Greeter {
	return Greeter{Message: Message(im.Getcontent())}
}

func main() {
	e, _, err := InitializeEvent()
	if err != nil {
		fmt.Println("Error initializing event:", err)
		return
	}
	e.Start()
}
