package broadcast

import (
	chatlog "goChat/server/chatLog"
	"sync"
)

func Broadcast(user string, line string, rooms []string) {
	if line == "" {
		return
	}
	var wg sync.WaitGroup
	for _, room := range rooms {
		wg.Add(1)
		go func(room string) {
			defer wg.Done()
			chatLog := chatlog.GetLog(room)
			chatLog.AddMessage(user, line)
		}(room)
	}
	wg.Wait()
}
