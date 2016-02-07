package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hearmeout"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"
)

var config struct {
	host     string
	username string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("[error] Specify the server")
		os.Exit(1)
	}
	config.host = os.Args[1]

	if len(os.Args) < 3 {
		user, _ := user.Current()
		config.username = user.Username
	} else {
		config.username = os.Args[2]
	}

	fmt.Printf("Connecting to %s...\n", config.host)
	conn, err := net.Dial("tcp", config.host)

	if err != nil {
		fmt.Println("[error] Connection refused")
		os.Exit(2)
	}

	fmt.Fprintln(conn, config.username)

	go watchIncoming(conn)
	watchStdin(conn)
}

func watchIncoming(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, msgerr := reader.ReadBytes('\n')
		if msgerr != nil {
			fmt.Println("[error] Shit happend!")
			sayerr := exec.Command("say", "Shit. Server died!").Start()
			if sayerr != nil {
				fmt.Println("[error] Something went wrong with the say command")
			}
			panic(msgerr)
		}

		var parsed = hearmeout.Message{}
		json.Unmarshal(message, &parsed)

		if parsed.From == "" {
			fmt.Printf("\r%s [%s]: %s\n", parsed.Time, "*", parsed.Body)
		} else {
			fmt.Printf("\r%s [%s]: %s\n", parsed.Time, parsed.From, parsed.Body)
		}
		fmt.Print("> ")

		sayerr := exec.Command("say", parsed.Body).Start()
		if sayerr != nil {
			fmt.Println("[error] Something went wrong with the say command")
		}
	}
}

func watchStdin(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[error] Shit happend!")
			panic(err)
		}

		conn.Write([]byte(message))

		var now = time.Now()
		var time = strconv.FormatInt(int64(now.Hour()), 10) + ":" + strconv.FormatInt(int64(now.Minute()), 10)

		fmt.Printf("\033[F%s [%s]: %s", time, config.username, message)
		fmt.Print("> ")
	}
}
