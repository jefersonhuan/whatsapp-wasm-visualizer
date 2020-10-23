package main

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestChat_parse(t *testing.T) {
	type fields struct {
		messages []string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int
	}{
		{
			name:   "ignores empty lines and returns empty map",
			fields: fields{messages: []string{"", "", "", ""}},
			want:   map[string]int{},
		},
		{
			name: "only counts lines starting with date",
			fields: fields{
				messages: []string{
					"10/9/19, 07:00 - Jeferson: Lorem Lorem Lorem",
					"Jeferson: Lorem Lorem Lorem\nMais Lorem aqui\nE mais aqui",
					"Jeferson: E mais uma linha",
					"11/11/19, 07:00 - Jeferson: Agora em outro dia",
				},
			},
			want: map[string]int{
				"10/9/19":  1,
				"11/11/19": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat := &Chat{
				messages: tt.fields.messages,
			}
			if got := chat.parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadChat(t *testing.T) {
	t.Skip()
	conversationFile, _ := os.OpenFile("./tests/multi_line.txt", os.O_RDONLY, 0444)

	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantChat *Chat
		wantErr  bool
	}{
		{
			"successfully reads lines from reader and stores them at Chat",
			args{reader: conversationFile},
			&Chat{
				messages: []string{
					"10/17/19, 18:30 - Messages and calls are end-to-end encrypted. No one outside of this chat, not even WhatsApp, can read or listen to them. Tap to learn more.",
					"10/17/19, 18:08 - Leonardo: Lorem Ipsum",
					"10/18/19, 12:31 - Jeferson: Lorem ipsum dolor sit amet,",
					"consectetur adipiscing elit.",
					"",
					"Nunc feugiat scelerisque.",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChat, err := loadChat(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChat, tt.wantChat) {
				t.Errorf("loadChat() gotChat = %v, want %v", gotChat, tt.wantChat)
			}
		})
	}
}
