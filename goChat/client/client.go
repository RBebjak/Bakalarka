package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type FetchResponse struct {
	Messages map[string][]Message `json:"messages"`
}

type Rooms struct {
	Rooms []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <username>")
		return
	}

	addr := "localhost:8080"
	username := os.Args[1]

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Send username
	fmt.Fprintf(conn, "%s\n", username)
	fmt.Println(readLine(reader))

	myRooms := new(Rooms)

	console := bufio.NewScanner(os.Stdin)
	for console.Scan() {
		room := strings.TrimSpace(console.Text())

		if room == "exit" {
			fmt.Printf("{Cant use %s as name of room}\n", room)
			continue
		}

		if room == "logout" {
			return
		}

		if room == "/broadcast" {
			fmt.Fprintf(conn, "%s\n", room)
			fmt.Print("Broadcasting message: ")
			console.Scan()
			broadcastingText := strings.TrimSpace(console.Text())
			fmt.Fprintf(conn, "%s\n", broadcastingText)
			fmt.Println(readLine(reader))
			continue
		}

		if room == "/fetchAll" {
			fmt.Fprintf(conn, "%s\n", room)
			fetchAll(reader)
			continue
		}

		// Send room
		fmt.Fprintf(conn, "%s\n", room)
		fmt.Println(readLine(reader))

		myRooms.Rooms = append(myRooms.Rooms, room)

		stopFetch := make(chan string)

		//Goroutine that ask for new messages
		go func() {
			for {
				select {
				case text := <-stopFetch:
					if text == "broadcasting" {
						_ = <-stopFetch
						continue
					}
					return

				case <-time.After(1 * time.Second):
					fmt.Fprintf(conn, "/fetch\n")
					line := readLine(reader)
					var resp struct {
						Messages []Message `json:"messages"`
					}
					if err := json.Unmarshal([]byte(line), &resp); err == nil {
						for _, msg := range resp.Messages {
							fmt.Printf("[%s] %s\n", msg.User, msg.Message)
						}
					}
				}
			}
		}()

		// Read user input from console
		for room != "" && console.Scan() {
			text := strings.TrimSpace(console.Text())
			if text == "" {
				continue
			}

			fmt.Fprintf(conn, "%s\n", text)

			if text == "exit" {
				close(stopFetch)
				room = ""
				fmt.Println(readLine(reader))
				break
			}

			if text == "logout" {
				close(stopFetch)
				return
			}

			if text == "/broadcast" {
				stopFetch <- "broadcasting"
				fmt.Print("Broadcasting message: ")
				console.Scan()
				broadcastingText := strings.TrimSpace(console.Text())
				fmt.Fprintf(conn, "%s\n", broadcastingText)
				stopFetch <- "end"
			}

			if text == "/fetchAll" {
				fetchAll(reader)
				continue
			}

			// Print server response
			fmt.Println(readLine(reader))
		}
	}
}

func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

func fetchAll(reader *bufio.Reader) {
	var resp FetchResponse

	jsonText := readLine(reader)
	fmt.Println(jsonText)

	if err := json.Unmarshal([]byte(jsonText), &resp); err != nil {
		fmt.Println("Failed to parse response:", err)
		return
	}

	for roomName, messages := range resp.Messages {
		if roomName == "/fetchAll" {
			continue
		}

		fmt.Printf("[%s]\n", roomName)

		for _, msg := range messages {
			fmt.Printf("[%s] %s\n", msg.User, msg.Message)
		}

		fmt.Println()
	}
}
