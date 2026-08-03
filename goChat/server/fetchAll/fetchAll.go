package fetchall

import (
	chatlog "goChat/server/chatLog"
	"goChat/server/message"
	"sync"
)

func FetchAll(user string, rooms []string) map[string][]message.Message {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var chatLog *chatlog.ChatLog
	allMessages := make(map[string][]message.Message)
	for _, room := range rooms {
		wg.Add(1)
		go func(room string, chatLog *chatlog.ChatLog) {
			defer wg.Done()
			mutex.Lock()
			chatLog = chatlog.GetLog(room)
			messages, _ := chatLog.GetMessagesSince(0, user)
			allMessages[room] = messages
			mutex.Unlock()
		}(room, chatLog)
	}
	wg.Wait()
	return allMessages
}
