package connection

import (
	"bufio"
	"encoding/json"
	"fmt"
	"goChat/server/broadcast"
	chatlog "goChat/server/chatLog"
	fetchall "goChat/server/fetchAll"
	"goChat/server/message"
	"net"
	"strings"
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

		if line == "/broadcast" {
			reader.Scan()
			lineForAll := strings.TrimSpace(reader.Text())
			broadcast.Broadcast(user, lineForAll, rooms)
			fmt.Fprint(conn, "{\"status\":\"sent broadcast\"}\n")
			continue
		}

		if line == "/fetchAll" {
			allMessages := fetchall.FetchAll(user, rooms)
			resp := struct {
				Messages map[string][]message.Message `json:"messages"`
			}{allMessages}
			data, _ := json.Marshal(resp)
			fmt.Fprintf(conn, "%s\n", data)
			continue
		}

		if chatLog == nil {
			if line == "/fetch" || line == "exit" || line == "/broadcast" {
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
