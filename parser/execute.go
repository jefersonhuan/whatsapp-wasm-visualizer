package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"time"
)

// é interessante procurar pelo dia correto, não não se estender por muito tempo
func findOrphan(ctx context.Context, timestamp int64, result *[][]int64) {
	_, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	for i := len(*result) - 1; i >= 0; i-- {
		if (*result)[i][0] == timestamp {
			(*result)[i][1] = (*result)[i][1] + 1
			fmt.Println("Índice encontrado para data orfã", timestamp)
			break
		}
	}
}

func Parse(reader io.Reader) (result [][]int64, nMessages int) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	result = [][]int64{}
	datePattern := regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{2}`)

	currentIndex := -1
	var currentTimestamp time.Time
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}

		if date := datePattern.FindString(scanner.Text()); date != "" {
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
			nMessages++
		}
	}
	return
}

func Convert(result [][]int64) []interface{} {
	data := make([]interface{}, len(result))
	for i, v := range result {
		data[i] = []interface{}{v[0] * 1000, v[1]}
	}
	return data
}
