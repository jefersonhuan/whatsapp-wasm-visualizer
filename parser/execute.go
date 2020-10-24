package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
)

type Chat struct {
	messages []string
}

func parseDate(date string) (uint32, error) {
	layouts := []string{"1/2", "01/2", "1/02", "01/02"}
	for _, layout := range layouts {
		t, err := time.Parse(layout+"/06", date)
		if err == nil {
			return uint32(t.Unix()), nil
		}
	}
	return 0, fmt.Errorf("date %s is not valid", date)
}

func LoadChat(reader io.Reader) (chat *Chat) {
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

func (chat *Chat) Parse() (result []uint32) {
	result = []uint32{}
	datePattern := regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{2}`)

	curIndex := -1
	var curTime uint32
	for _, line := range chat.messages {
		if date := datePattern.FindString(line); date != "" {
			timestamp, err := parseDate(date)
			if err != nil || timestamp < curTime {
				continue
			} else if timestamp != curTime {
				curTime = timestamp
				curIndex += 2
				result = append(result, []uint32{timestamp, 1}...)
			} else {
				result[curIndex] = result[curIndex] + 1
			}
		}
	}
	return
}
