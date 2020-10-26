package parser

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func testFile(t *testing.T) *os.File {
	f, err := os.OpenFile("./tests/large_conversation.txt", os.O_RDONLY, 0444)
	if err != nil {
		t.Errorf("file for tests couldn't be found: %v", err)
	}
	return f
}

func TestLoadChat(t *testing.T) {
	file := testFile(t)

	tests := []struct {
		name            string
		reader          io.Reader
		wantMessagesLen int
	}{
		{
			name:            "",
			reader:          file,
			wantMessagesLen: 1054,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadChat(tt.reader); len(got.messages) != tt.wantMessagesLen {
				t.Errorf("LoadChat() = %v, want %v", len(got.messages), tt.wantMessagesLen)
			}
		})
	}
}

func TestChat_Parse(t *testing.T) {
	type fields struct {
		messages []string
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]int64
	}{
		{
			name:   "ignores empty lines and returns empty map",
			fields: fields{messages: []string{"", "", "", ""}},
			want:   [][]int64{},
		},
		{
			name: "only counts lines starting with date",
			fields: fields{
				messages: []string{
					"10/9/19, 14:00 - Jeferson: Lorem Lorem Lorem",
					"Jeferson: Lorem Lorem Lorem\nMais Lorem aqui\nE mais aqui",
					"Jeferson: E mais uma linha",
					"11/11/19, 07:00 - Jeferson: Agora em outro dia",
					"11/11/19, 07:01 - Jeferson: Mas duplicado",
				},
			},
			want: [][]int64{{1570579200, 1}, {1573430400, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{
				messages: tt.fields.messages,
			}
			if got := chat.Parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
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

func TestIntegration(t *testing.T) {
	file := testFile(t)
	got := LoadChat(file).Parse()
	if len(got) != 500 {
		t.Errorf("LoadChat().Parse() = %v, want 500", len(got))
	}
}

func BenchmarkLoad_Chat(b *testing.B) {
	file, err := os.OpenFile("./tests/small_conversation.txt", os.O_RDONLY, 0444)
	if err != nil {
		b.Errorf("file for tests couldn't be found: %v", err)
	}
	for n := 0; n < b.N; n++ {
		LoadChat(file)
	}
}

func BenchmarkChat_Parse(b *testing.B) {
	chat := Chat{messages: []string{
		"10/9/19, 07:00 - Jeferson: Lorem Lorem Lorem",
		"Jeferson: Lorem Lorem Lorem\nMais Lorem aqui\nE mais aqui",
		"Jeferson: E mais uma linha",
		"11/11/19, 07:00 - Jeferson: Agora em outro dia",
	}}
	for n := 0; n < b.N; n++ {
		chat.Parse()
	}
}
