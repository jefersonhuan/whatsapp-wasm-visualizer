package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"
)

type Chat struct {
	messages []string
}

// é interessante procurar pelo dia correto, não não se estender ou permitir vazamentos
func findOrphan(ctx context.Context, timestamp int64, result *[][]int64) {
	interrupt := false
	_, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer (func() {
		interrupt = true
		cancel()
	})()

	done := make(chan bool)
	go func() {
		for i := len(*result) - 1; i >= 0 && !interrupt; i-- {
			if (*result)[i][0] == timestamp {
				(*result)[i][1] = (*result)[i][1] + 1
				fmt.Println("Índice encontrado para data orfã", timestamp)
				break
			}
		}
		done <- true
	}()
	<-done
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
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	result = [][]int64{}
	datePattern := regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{2}`)

	currentIndex := -1
	var currentTimestamp time.Time
	for _, line := range chat.messages {
		if date := datePattern.FindString(line); date != "" {
			timestamp, err := time.Parse("1/2/06", date)
			if err != nil {
				continue
			} else if timestamp.Before(currentTimestamp) {
				// algumas mensagens ficam perdidas (muitas vezes mídia)
				// então tenta-se, por um acaso, encontrar a verdadeira data
				findOrphan(ctx, timestamp.Unix(), &result)
			} else if timestamp != currentTimestamp {
				currentTimestamp = timestamp
				currentIndex++
				result = append(result, []int64{timestamp.Unix(), 1})
			} else {
				result[currentIndex][1] = result[currentIndex][1] + 1
			}
		}
	}
	return
}

func Convert(result [][]int64) (data []interface{}) {
	var wg sync.WaitGroup
	data = make([]interface{}, len(result))

	// por se tratar de uma operação com interface, a goroutine acaba sendo vantajosa na maioria dos casos
	// mas sendo mais vantajoso em conversas de maior volume
	wg.Add(len(result))
	go func() {
		for i, v := range result {
			data[i] = []interface{}{v[0] * 1000, v[1]}
			wg.Done()
		}
	}()
	wg.Wait()
	return
}
