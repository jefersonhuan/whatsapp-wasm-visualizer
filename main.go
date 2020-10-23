package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

type Chat struct {
	messages []string
}

func loadChat(reader io.Reader) (chat *Chat, err error) {
	chat = &Chat{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}
		chat.messages = append(chat.messages, scanner.Text())
	}
	return
}

func (chat *Chat) parse() map[string]int {
	result := map[string]int{}
	datePattern := regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{2}`)

	for _, line := range chat.messages {
		if date := datePattern.FindString(line); date != "" {
			result[date] = result[date] + 1
		}
	}
	return result
}

func main() {
	start := time.Now()
	fmt.Println("Initializing WhatsApp Chat parser...")

	file, err := os.OpenFile("./conversations/xlarge.txt", os.O_RDONLY, 044)
	if err != nil {
		return
	}
	chat, err := loadChat(file)
	if err != nil {
		log.Fatalln(err)
	}

	result := chat.parse()
	for k, v := range result {
		fmt.Println(k, ":", v)
	}

	fmt.Println(len(result), "results found")
	fmt.Println("Elapsed", time.Now().Sub(start))
}
