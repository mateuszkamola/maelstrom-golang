package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

var messageCounter int = 1

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	for {
		handle(line, os.Stdout)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		line, err = reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
	}
}

func handle(line string, output *os.File) {
	msg := parseMsg(line)
	var responseBody map[string]interface{}
	fmt.Println("Hello")
	if msg.Body["type"].(string) == "echo" {
		responseBody = createEchoOkResponse(msg)
	} else {
		responseBody = createInitOkResponse(msg)
	}
	response := Message{
		Src: msg.Dst,
		Dst: msg.Src,
		Body: responseBody,
	}
	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, err = output.Write(data)
	if err != nil {
		panic(err)
	}
	_, err = output.Write([]byte{'\n'})
}

func createEchoOkResponse(msg Message) map[string]interface{} {
	messageCounter++
	response := make(map[string]interface{})
	response["type"] = "echo_ok"
	response["msg_id"] = messageCounter
	msgId, ok := msg.Body["msg_id"]
	if ok {
		response["in_reply_to"] = int(msgId.(float64))
	}
	response["echo"] = msg.Body["echo"].(string)
	return response
}

func createInitOkResponse(msg Message) map[string]interface{} {
	messageCounter++
	response := make(map[string]interface{})
	response["type"] = "init_ok"
	response["msg_id"] = messageCounter
	msgId, ok := msg.Body["msg_id"]
	if ok {
		response["in_reply_to"] = int(msgId.(float64))
	}
	return response
}

func parseMsg(line string) Message {
	msg := Message{}
	json.Unmarshal([]byte(line), &msg) 
	return msg
}

type Message struct {
	Src string `json:"src"`
	Dst string `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

type MessageBody struct {
	Type string `json:"type"`
	MsgId int `json:"msg_id"`
	InReplyTo int `json:"in_reply_to"`
}

type EchoBodyOk struct {
	Type string `json:"type"`
	MsgId int `json:"msg_id"`
	InReplyTo int `json:"in_reply_to"`
	Echo string `json:"echo"`
}

type InitBody struct {
	Type string `json:"type"`
	MsgId int `json:"msg_id"`
	InReplyTo int `json:"in_reply_to"`
	NodeId string `json:"node_id"`
	NodeIds string `json:"node_ids"`
}
