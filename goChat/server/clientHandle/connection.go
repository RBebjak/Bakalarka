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

	// First line: username
	if !reader.Scan() {
		return
	}
	user := strings.TrimSpace(reader.Text())
	if user == "" {
		fmt.Fprintf(conn, "{\"error\":\"username required\"}\n")
		return
	}

	var chatLog *chatlog.ChatLog
	lastSeen := 0
	currentRoom := ""

	fmt.Fprintf(conn, "{\"status\":\"hello %s\"}\n", user)

	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())

		switch line {

		case "":
			continue

		case "/fetch":
			if chatLog == nil {
				fmt.Fprintf(conn, "{\"error\":\"not in a room\"}\n")
				continue
			}
			newMsgs, newIndex := chatLog.GetMessagesSince(lastSeen, user)
			resp := struct {
				Messages []message.Message `json:"messages"`
			}{newMsgs}
			data, _ := json.Marshal(resp)
			fmt.Fprintf(conn, "%s\n", data)
			lastSeen = newIndex

		case "exit":
			if chatLog != nil {
				fmt.Fprintf(conn, "{\"status\":\"left room %s\"}\n", currentRoom)
				chatLog = nil
				currentRoom = ""
				lastSeen = 0
			} else {
				fmt.Fprintf(conn, "{\"error\":\"not in a room\"}\n")
			}

		case "logout":
			fmt.Fprintf(conn, "{\"status\":\"bye %s\"}\n", user)
			return

		case "fetchAll":
			if chatLog == nil {
				fmt.Fprintf(conn, "{\"All messages:}\n")
				allmsg := sync.WaitGroup{}
				for reader.Scan() {
					roomLine := strings.TrimSpace(reader.Text())
					if roomLine == "end" {
						allmsg.Wait()
						allmsg.Done()
						break
					}
					allmsg.Add(1)

				}
			}

		default:
			if chatLog == nil {
				room := line
				chatLog = chatlog.GetLog(room)
				currentRoom = room
				lastSeen = 0
				fmt.Fprintf(conn, "{\"status\":\"joined room %s\"}\n", room)
			} else {
				chatLog.AddMessage(user, line)
				fmt.Fprintf(conn, "{\"status\":\"sent\"}\n")
			}
		}
	}
}
