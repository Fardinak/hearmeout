package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hearmeout"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var config struct {
	host string
}

// Client struct
type Client struct {
	conn net.Conn

	IP       string
	Username string
}

var clients []*Client

func main() {
	if len(os.Args) < 2 {
		config.host = "0.0.0.0:35000"
	} else {
		config.host = os.Args[1]
	}

	fmt.Printf("Listening on %s\n", config.host)
	sock, err := net.Listen("tcp", config.host)
	if err != nil {
		fmt.Printf("[error] Can not start listening on %s\n", config.host)
		fmt.Println("[error] Make sure the port is free")
		os.Exit(10)
	}

	for {
		conn, connerr := sock.Accept()
		if connerr != nil {
			fmt.Println("[error] Can not accept connections")
			os.Exit(11)
		}

		client := Client{conn: conn, IP: conn.RemoteAddr().String()}
		clients = append(clients, &client)

		channel := make(chan string)
		go watchInput(channel, &client)
		go handleDistribution(channel, &client)
	}
}

func watchInput(out chan string, client *Client) {
	defer close(out)

	reader := bufio.NewReader(client.conn)

	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("[error] Client %s connection closed\n", client.IP)
		client.conn.Close()
		removeEntry(client, clients)
		return
	}
	client.Username = strings.TrimSpace(username)

	fmt.Printf("* %s@%s connected\n", client.Username, client.IP)
	sendMessage(
		fmt.Sprintf("%s connected", client.Username),
		client,
		time.Now(),
		clients,
	)

	var usernames []string
	for _, c := range clients {
		usernames = append(usernames, c.Username)
	}
	sendMessage(
		fmt.Sprintf("Connected users: %s", strings.Join(usernames, ", ")),
		&Client{},
		time.Now(),
		[]*Client{client},
	)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[error] Client %s connection closed\n", client.IP)
			fmt.Printf("* %s@%s disconnected\n", client.Username, client.IP)
			sendMessage(
				fmt.Sprintf("%s disconnected", client.Username),
				&Client{},
				time.Now(),
				clients,
			)
			client.conn.Close()
			removeEntry(client, clients)
			return
		}

		out <- string(message)
	}
}

func handleDistribution(in chan string, client *Client) {
	for {
		message := <-in
		if message != "" {
			message = strings.TrimSpace(message)

			sendMessage(message, client, time.Now(), clients)
		}
	}
}

func sendMessage(message string, from *Client, at time.Time, to []*Client) {
	marshalled, err := json.Marshal(hearmeout.Message{
		From: from.Username,
		Body: message,
		Time: strconv.FormatInt(int64(at.Hour()), 10) + ":" + strconv.FormatInt(int64(at.Minute()), 10),
	})
	marshalled = append(marshalled, '\n')

	if err != nil {
		fmt.Printf("[error] Failed to marshal a message from %s@%s\n", from.Username, from.IP)
	}

	for _, client := range to {
		if client.IP == from.IP && client.Username == from.Username {
			continue
		}
		client.conn.Write(marshalled)
	}
}

// remove client entry from stored clients
func removeEntry(client *Client, arr []*Client) []*Client {
	rtn := arr
	index := -1
	for i, value := range arr {
		if value == client {
			index = i
			break
		}
	}

	if index >= 0 {
		// we have a match, create a new array without the match
		rtn = make([]*Client, len(arr)-1)
		copy(rtn, arr[:index])
		copy(rtn[index:], arr[index+1:])
	}

	return rtn
}
