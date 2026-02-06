package chatlog

import (
	message "goChat/server/message"
)

type ChatLog struct {
	Messages []message.Message
	writeCh  chan message.Message
}

func NewChatLog() *ChatLog {
	cl := &ChatLog{
		writeCh: make(chan message.Message, 100),
	}
	go cl.run()
	return cl
}

func (cl *ChatLog) run() {
	for msg := range cl.writeCh {
		cl.Messages = append(cl.Messages, msg)
	}
}

var logs = make(map[string]*ChatLog)

func GetLog(name string) *ChatLog {
	if logs[name] == nil {
		logs[name] = NewChatLog()
	}
	return logs[name]
}

func (cl *ChatLog) AddMessage(user, text string) {
	cl.writeCh <- message.Message{User: user, Message: text}
}

func (cl *ChatLog) GetMessagesSince(lastSeen int, excludeUser string) ([]message.Message, int) {
	if lastSeen < len(cl.Messages) {
		var result []message.Message
		for _, msg := range cl.Messages[lastSeen:] {
			if msg.User != excludeUser {
				result = append(result, msg)
			}
		}
		return result, len(cl.Messages)
	}
	return []message.Message{}, lastSeen
}
