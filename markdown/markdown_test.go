package markdown

import (
	"testing"
)

func TestText2tag(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				text: "test",
			},
			want: "test",
		},
		{
			name: "test with space",
			args: args{
				text: "test test",
			},
			want: "test-test",
		},
		{
			name: "test with fullwidth",
			args: args{
				text: "test　test",
			},
			want: "testtest",
		},
		{
			name: "test with fullwidth chars",
			args: args{
				text: "test　！＠＃＄％＾＆＊（）＋｜〜＝￥｀「」｛｝；’：”、。・＜＞？【】『』《》〔〕［］‹›«»〘〙〚〛test",
			},
			want: "testtest",
		},
		{
			name: "test with #",
			args: args{
				text: "test#test",
			},
			want: "testtest",
		},
		{
			name: "test with .",
			args: args{
				text: "test.test",
			},
			want: "testtest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Text2tag(tt.args.text); got != tt.want {
				t.Errorf("Text2tag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildHeading(t *testing.T) {
	type args struct {
		level int
		text  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "level 1",
			args: args{
				level: 1,
				text:  "test",
			},
			want: "# test",
		},
		{
			name: "level 6",
			args: args{
				level: 6,
				text:  "test",
			},
			want: "###### test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildHeading(tt.args.level, tt.args.text); got != tt.want {
				t.Errorf("BuildHeading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildLink(t *testing.T) {
	type args struct {
		text        string
		destination string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				text:        "test",
				destination: "https://example.com",
			},
			want: "[test](https://example.com)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildLink(tt.args.text, tt.args.destination); got != tt.want {
				t.Errorf("BuildLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildList(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				text: "test",
			},
			want: "- test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildList(tt.args.text); got != tt.want {
				t.Errorf("BuildList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildOrderedList(t *testing.T) {
	type args struct {
		order int
		text  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "order 1",
			args: args{
				order: 1,
				text:  "test",
			},
			want: "1. test",
		},
		{
			name: "order 6",
			args: args{
				order: 6,
				text:  "test",
			},
			want: "6. test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildOrderedList(tt.args.order, tt.args.text); got != tt.want {
				t.Errorf("BuildOrderedList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildCheckbox(t *testing.T) {
	type args struct {
		text    string
		checked bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "checked",
			args: args{
				text:    "test",
				checked: true,
			},
			want: "- [x] test",
		},
		{
			name: "unchecked",
			args: args{
				text:    "test",
				checked: false,
			},
			want: "- [ ] test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildCheckbox(tt.args.text, tt.args.checked); got != tt.want {
				t.Errorf("BuildCheckbox() = %v, want %v", got, tt.want)
			}
		})
	}
}
