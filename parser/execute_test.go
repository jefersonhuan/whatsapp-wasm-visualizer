package parser

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func testFile(t *testing.T) *os.File {
	f, err := os.OpenFile("./tests/large_conversation.txt", os.O_RDONLY, 0444)
	if err != nil {
		t.Errorf("file for tests couldn't be found: %v", err)
	}
	return f
}

func TestChat_Parse(t *testing.T) {
	type fields struct {
		messages string
	}
	tests := []struct {
		name                 string
		fields               fields
		want                 [][]int64
		wantNumberOfMessages int
	}{
		{
			name:                 "ignores empty lines and returns empty map",
			fields:               fields{messages: "\n\n\n\n"},
			want:                 [][]int64{},
			wantNumberOfMessages: 0,
		},
		{
			name: "only counts lines starting with date",
			fields: fields{
				messages: "10/9/19, 14:00 - Jeferson: Lorem Lorem Lorem\n" +
					"Jeferson: Lorem Lorem Lorem\nMais Lorem aqui\nE mais aqui\n" +
					"Jeferson: E mais uma linha\n" +
					"11/11/19, 07:00 - Jeferson: Agora em outro dia\n" +
					"11/11/19, 07:01 - Jeferson: Mas duplicado\n",
			},
			want:                 [][]int64{{1570579200, 1}, {1573430400, 2}},
			wantNumberOfMessages: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, nMes := Parse(strings.NewReader(tt.fields.messages))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
			if nMes != tt.wantNumberOfMessages {
				t.Errorf("Parse(), nMessages = %v, want %v", nMes, tt.wantNumberOfMessages)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	input := [][]int64{{1570579200, 1}, {1573430400, 2}, {1573450400, 50}}
	got := Convert(input)
	if len(got) != len(input) {
		t.Errorf("Convert() = %v, want %v", len(got), len(input))
	}

	for i, d := range got {
		intermediary := d.([]interface{})
		if intermediary[0] != input[i][0]*1000 || intermediary[1] != input[i][1] {
			t.Errorf("Convert() = got %v, want [%v, %v]", intermediary, input[i][0]*1000, input[i][1])
		}
	}
}

func BenchmarkChat_Parse(b *testing.B) {
	input := "10/9/19, 07:00 - Jeferson: Lorem Lorem Lorem\n" +
		"Jeferson: Lorem Lorem Lorem\nMais Lorem aqui\nE mais aqui\n" +
		"Jeferson: E mais uma linha\n" +
		"11/11/19, 07:00 - Jeferson: Agora em outro dia\n"
	reader := strings.NewReader(input)
	for n := 0; n < b.N; n++ {
		Parse(reader)
	}
}
