package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"
)

type Chat struct {
	messages []string
}

func parseDate(date string) (int64, error) {
	layouts := []string{"1/2", "01/2", "1/02", "01/02"}
	for _, layout := range layouts {
		t, err := time.Parse(layout+"/06", date)
		if err == nil {
			return t.Unix() * 1000, nil
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

func (chat *Chat) Parse() (result [][]int64) {
	result = [][]int64{}
	datePattern := regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{2}`)

	curIndex := -1
	var curTime int64
	for _, line := range chat.messages {
		if date := datePattern.FindString(line); date != "" {
			timestamp, err := parseDate(date)
			if err != nil || timestamp < curTime {
				continue
			} else if timestamp != curTime {
				curTime = timestamp
				curIndex++
				result = append(result, []int64{timestamp, 1})
			} else {
				result[curIndex][1] = result[curIndex][1] + 1
			}
		}
	}
	return
}

func Convert(result [][]int64) (data []interface{}) {
	var wg sync.WaitGroup
	data = make([]interface{}, len(result))

	wg.Add(len(result))
	go func() {
		for i, v := range result {
			data[i] = []interface{}{v[0], v[1]}
			wg.Done()
		}
	}()
	wg.Wait()
	return
}
