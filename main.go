package main

import (
        "encoding/json"
        "fmt"
        "time"

        "github.com/gorilla/websocket"
        ping "github.com/t0stbrot/go-ping"
        tracert "github.com/t0stbrot/go-tracert"
)

type Message struct {
        Type    string  `json:"type"`
        Action  string  `json:"action,omitempty"`
        Target  string  `json:"target,omitempty"`
        Content any     `json:"content,omitempty"`
        Options Options `json:"options,omitempty"`
}

type Options struct {
        Target  string `json:"target"`
        Proto   string `json:"proto,omitempty"`
        Timeout int    `json:"timeout,omitempty"`
        Maxhops int    `json:"maxhops,omitempty"`
        TTL     int    `json:"ttl,omitempty"`
}

type Register struct {
        Version string `json:"version"`
}

func main() {
        err := runPrClient()
        if err != nil {
                fmt.Printf("\nWebSocket client error: %v", err)
        }
}

// functions

func traceroute6(conn *websocket.Conn, target string, maxhops int, timeout int) {
        res := tracert.Traceroute6(target, maxhops, timeout)
        jsonOutput, _ := json.Marshal(res)
        returnMsg := Message{Type: "result", Action: "traceroute", Target: target, Content: string(jsonOutput)}
        conn.WriteJSON(returnMsg)
}

func icmp6(conn *websocket.Conn, target string, ttl int, timeout int) {
        res := ping.Ping6(target, ttl, timeout)
        jsonOutput, _ := json.Marshal(res)
        returnMsg := Message{Type: "result", Action: "icmp", Target: target, Content: string(jsonOutput)}
        conn.WriteJSON(returnMsg)
}

func traceroute4(conn *websocket.Conn, target string, maxhops int, timeout int) {
        res := tracert.Traceroute4(target, maxhops, timeout)
        jsonOutput, _ := json.Marshal(res)
        returnMsg := Message{Type: "result", Action: "traceroute", Target: target, Content: string(jsonOutput)}
        conn.WriteJSON(returnMsg)
}

func icmp4(conn *websocket.Conn, target string, ttl int, timeout int) {
        res := ping.Ping4(target, ttl, timeout)
        jsonOutput, _ := json.Marshal(res)
        returnMsg := Message{Type: "result", Action: "icmp", Target: target, Content: string(jsonOutput)}
        conn.WriteJSON(returnMsg)
}

func runPrClient() error {
        // Server to connect to
        serverURL := "wss://winpr.t0stbrot.net/probe"

        for {
                // Dialing the TCP address and trying to connect, then deferring the closing part
                conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
                if err != nil {
                        fmt.Printf("\nfailed connect to master; retry in 1 seconds")
                        time.Sleep(1 * time.Second)
                        continue
                }
                defer conn.Close()

                // Registering at the central service
                registerMsg := Message{Type: "register", Content: Register{Version: "v0.0.2"}}
                conn.WriteJSON(registerMsg)

                fmt.Printf("\nconnected to master")

                // ticker for keepalive
                go func() {
                        ticker := time.NewTicker(5 * time.Second)
                        defer ticker.Stop()

                        for range ticker.C {
                                returnMsg := Message{Type: "keepalive"}
                                conn.WriteJSON(returnMsg)
                        }
                }()

                // for loop for reading of the socket
                for {
                        // listen
                        _, message, err := conn.ReadMessage()
                        if err != nil {
                                fmt.Printf("\nconnection closed; reconnecting")
                                break
                        }

                        // parse message
                        var msg Message
                        if err := json.Unmarshal(message, &msg); err != nil {
                                fmt.Printf("\nFailed to parse message: %v", err)
                                continue
                        }

                        if msg.Target != "" && msg.Action != "" && msg.Type == "request" {
                                if msg.Action == "icmp" {
                                        fmt.Printf("\nrequested %v for %v", msg.Action, msg.Target)
                                        ttl := 64
                                        timeout := 1000
                                        if msg.Options.Target != "" {
                                                ttl = msg.Options.TTL
                                                timeout = msg.Options.Timeout
                                        }
                                        icmp4(conn, msg.Target, ttl, timeout)
                                } else if msg.Action == "traceroute" {
                                        fmt.Printf("\nrequested %v for %v", msg.Action, msg.Target)
                                        maxhops := 30
                                        timeout := 1000
                                        if msg.Options.Target != "" {
                                                maxhops = msg.Options.Maxhops
                                                timeout = msg.Options.Timeout
                                        }
                                        traceroute4(conn, msg.Target, maxhops, timeout)
                                } else if msg.Action == "icmp6" {
                                        fmt.Printf("\nrequested %v for %v", msg.Action, msg.Target)
                                        ttl := 64
                                        timeout := 1000
                                        if msg.Options.Target != "" {
                                                ttl = msg.Options.TTL
                                                timeout = msg.Options.Timeout
                                        }
                                        icmp6(conn, msg.Target, ttl, timeout)
                                } else if msg.Action == "traceroute6" {
                                        fmt.Printf("\nrequested %v for %v", msg.Action, msg.Target)
                                        maxhops := 30
                                        timeout := 1000
                                        if msg.Options.Target != "" {
                                                maxhops = msg.Options.Maxhops
                                                timeout = msg.Options.Timeout
                                        }
                                        traceroute6(conn, msg.Target, maxhops, timeout)
                                }
                                continue
                        }
                }
                // if the for loop breaks, wait for 1 seconds for the top-level for loop to restart the connection
                time.Sleep(1 * time.Second)
        }
}
