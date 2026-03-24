package clientHandle

import (
	"bufio"
	"encoding/json"
	"fmt"
	chatlog "goChat/server/chatLog"
	"goChat/server/message"
	"net"
	"strings"
	"sync"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewScanner(conn)

	if !reader.Scan() {
		return
	}
	user := strings.TrimSpace(reader.Text())
	if user == "" {
		fmt.Fprintf(conn, "{\"error\":\"username required\"}\n")
		return
	}

	var chatLog *chatlog.ChatLog
	var rooms []string
	lastSeen := 0
	currentRoom := ""

	fmt.Fprintf(conn, "{\"status\":\"hello %s\"}\n", user)

	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())

		if chatLog == nil {
			if line == "/fetch" || line == "exit" {
				fmt.Fprintf(conn, "{\"error\":\"not in a room\"}\n")
				continue
			}
			room := line
			rooms = append(rooms, room)
			chatLog = chatlog.GetLog(room)
			currentRoom = room
			lastSeen = 0
			fmt.Fprintf(conn, "{\"status\":\"joined room %s\"}\n", room)
			continue
		}

		switch line {
		case "logout":
			fmt.Fprintf(conn, "{\"status\":\"logged out\"}\n")
			return

		case "":
			continue

		case "/fetch":
			newMsgs, newIndex := chatLog.GetMessagesSince(lastSeen, user)
			resp := struct {
				Messages []message.Message `json:"messages"`
			}{newMsgs}
			data, _ := json.Marshal(resp)
			fmt.Fprintf(conn, "%s\n", data)
			lastSeen = newIndex

		case "/broadcast":
			reader.Scan()
			line = strings.TrimSpace(reader.Text())
			Broadcast(user, line, rooms)

		case "/fetchAll":
			allMessages := FetchAll(user, rooms)
			resp := struct {
				Messages map[string][]message.Message `json:"messages"`
			}{allMessages}
			data, _ := json.Marshal(resp)
			fmt.Fprintf(conn, "%s\n", data)

		case "exit":
			fmt.Fprintf(conn, "{\"status\":\"left room %s\"}\n", currentRoom)
			chatLog = nil
			currentRoom = ""
			lastSeen = 0

		default:
			chatLog.AddMessage(user, line)
			fmt.Fprintf(conn, "{\"status\":\"sent\"}\n")
		}
	}
}

func Broadcast(user string, line string, rooms []string) {
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

func FetchAll(user string, rooms []string) map[string][]message.Message {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	allMessages := make(map[string][]message.Message)
	for _, room := range rooms {
		wg.Add(1)
		go func(room string) {
			defer wg.Done()
			defer mutex.Unlock()
			chatLog := chatlog.GetLog(room)
			messages, _ := chatLog.GetMessagesSince(0, user)
			mutex.Lock()
			allMessages[room] = messages
		}(room)
	}
	wg.Wait()
	return allMessages
}
