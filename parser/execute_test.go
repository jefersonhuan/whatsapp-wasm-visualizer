package parser

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func testFile(t *testing.T) *os.File {
	f, err := os.OpenFile("./tests/conversation.txt", os.O_RDONLY, 0444)
	if err != nil {
		t.Errorf("file for tests couldn't be found: %v", err)
	}
	return f
}

func Test_parseDate(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    uint32
		wantErr bool
	}{
		{
			name:    "recognizes and returns timestamp from date as 1/2/20",
			arg:     "1/2/20",
			want:    1577923200,
			wantErr: false,
		},
		{
			name:    "recognizes and returns timestamp from date as 01/2/20",
			arg:     "01/2/20",
			want:    1577923200,
			wantErr: false,
		},
		{
			name:    "recognizes and returns timestamp from date as 1/02/20",
			arg:     "1/02/20",
			want:    1577923200,
			wantErr: false,
		},
		{
			name:    "recognizes and returns timestamp from date as 01/02/20",
			arg:     "1/2/20",
			want:    1577923200,
			wantErr: false,
		},
		{
			name:    "returns error for non-recognized format",
			arg:     "1/322/20",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.arg)
			if tt.wantErr && err == nil {
				t.Errorf("parseDate() = wanted error, but got none")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
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
			wantMessagesLen: 1000,
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

func TestLoadChatAndParse(t *testing.T) {
	file := testFile(t)
	got := LoadChat(file).Parse()
	if len(got) != 2 {
		t.Errorf("LoadChat().Parse() = %v, want 2", len(got))
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
			want: [][]int64{{1570579200000, 1}, {1573430400000, 2}},
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

func BenchmarkLoad_Chat(b *testing.B) {
	file, err := os.OpenFile("./tests/conversation.txt", os.O_RDONLY, 0444)
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
