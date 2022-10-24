package orders

import (
	"reflect"
	"testing"
)

func Test_parseIds(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantIds []int
		wantErr bool
	}{
		{
			name: "parseIds_1",
			args: args{
				input: "1",
			},
			wantIds: []int{1},
			wantErr: false,
		},
		{
			name: "parseIds_2",
			args: args{
				input: "1 2",
			},
			wantIds: []int{1, 2},
			wantErr: false,
		},
		{
			name: "parseIds_3",
			args: args{
				input: "1  2",
			},
			wantIds: []int{1, 2},
			wantErr: false,
		},
		{
			name: "parseIds_4",
			args: args{
				input: "1  a",
			},
			wantIds: []int{1},
			wantErr: true,
		},
		{
			name: "parseIds_5",
			args: args{
				input: "28 29",
			},
			wantIds: []int{28, 29},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIds, err := parseIds(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIds, tt.wantIds) {
				t.Errorf("parseIds() gotIds = %v, want %v", gotIds, tt.wantIds)
			}
		})
	}
}

func Test_parseIdAndStr(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name        string
		args        args
		wantId      int
		wantContent string
		wantErr     bool
	}{
		{
			name:        "parseIdAndStr_1",
			args:        args{input: ""},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "parseIdAndStr_2",
			args:        args{input: "1"},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "parseIdAndStr_3",
			args:        args{input: "a"},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "parseIdAndStr_4",
			args:        args{input: "1 1"},
			wantId:      1,
			wantContent: "1",
			wantErr:     false,
		},
		{
			name:        "parseIdAndStr_5",
			args:        args{input: "1 1 a"},
			wantId:      1,
			wantContent: "1 a",
			wantErr:     false,
		},
		{
			name:        "parseIdAndStr_6",
			args:        args{input: "1  1 a"},
			wantId:      1,
			wantContent: "1 a",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotContent, err := parseIdAndStr(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIdAndStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("parseIdAndStr() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotContent != tt.wantContent {
				t.Errorf("parseIdAndStr() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func Test_parseDeadline(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name         string
		args         args
		wantId       int
		wantDeadline string
		wantErr      bool
	}{
		{
			name:         "parseDeadline_1",
			args:         args{input: ""},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_2",
			args:         args{input: "1"},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_3",
			args:         args{input: "2022-10-19"},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_4",
			args:         args{input: "1 2022-10"},
			wantId:       1,
			wantDeadline: "2022-10",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_5",
			args:         args{input: "1 2022-10-19"},
			wantId:       1,
			wantDeadline: "2022-10-19",
			wantErr:      false,
		},
		{
			name:         "parseDeadline_6",
			args:         args{input: "1  2022-10-19 a"},
			wantId:       1,
			wantDeadline: "2022-10-19 a",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_7",
			args:         args{input: "1 10-19"},
			wantId:       1,
			wantDeadline: "2022-10-19",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotDeadline, err := parseDeadline(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDeadline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("parseDeadline() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotDeadline != tt.wantDeadline {
				t.Errorf("parseDeadline() gotDeadline = %v, want %v", gotDeadline, tt.wantDeadline)
			}
		})
	}
}
