package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

func startConsole() {
	log.Println("You can now type commands into console.")
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) == 0 {
				continue
			}
			verb, rest, _ := strings.Cut(scanner.Text(), " ")
			switch verb {
			case "list":
				lobbiesMu.Lock()
				for id := range lobbies {
					log.Println(id)
				}
				lobbiesMu.Unlock()
			case "kill":
				if len(fields) > 1 {
					if fields[1] == "all" {
						lobbiesMu.Lock()
						ids := make([]string, 0, len(lobbies))
						for id := range lobbies {
							ids = append(ids, id)
						}
						lobbiesMu.Unlock()
						for _, id := range ids {
							killLobby(id)
						}
					} else {
						for i := range len(fields) - 1 {
							killLobby(fields[i+1])
						}
					}
				}
			case "end":
				if len(fields) > 1 {
					if fields[1] == "all" {
						lobbiesMu.Lock()
						ids := make([]string, 0, len(lobbies))
						for id := range lobbies {
							ids = append(ids, id)
						}
						lobbiesMu.Unlock()
						for _, id := range ids {
							endLobby(id)
						}
					} else {
						for i := range len(fields) - 1 {
							endLobby(fields[i+1])
						}
					}
				}
			case "unregister":
				if len(fields) > 1 {
					if fields[1] == "all" {
						unregisterAll()
					} else {
						for i := range len(fields) - 1 {
							unregisterLobby(fields[i+1])
						}
					}
				}
			case "create":
				if len(fields) > 1 {
					var config LobbyConfig
					if err := json.Unmarshal([]byte(rest), &config); err != nil {
						log.Println("Error parsing json:", err)
						continue
					}
					createLobby(config)
				}
			case "status":
				if len(fields) > 1 {
					rest = strings.TrimSpace(rest)
					if rest == "clear" {
						setStatus("", "")
					} else {
						i := strings.LastIndexByte(rest, ' ')
						text := rest
						color := "red"
						if i > 0 && strings.HasPrefix(rest[i+1:], "#") {
							text = strings.TrimSpace(rest[:i])
							color = strings.TrimPrefix(rest[i+1:], "#")
						}
						setStatus(text, color)
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println("Console read error:", err)
		}
	}()
}
