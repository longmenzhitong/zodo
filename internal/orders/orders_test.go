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

func Test_parseInput(t *testing.T) {
	type args struct {
		input  string
		orders []string
	}
	orders := []string{"help", "rmd+", "rmd"}
	tests := []struct {
		name      string
		args      args
		wantOrder string
		wantVal   string
	}{
		{
			name: "parseInput_1",
			args: args{
				input:  "help",
				orders: orders,
			},
			wantOrder: "help",
			wantVal:   "",
		},
		{
			name: "parseInput_2",
			args: args{
				input:  "rmd 1 11:02",
				orders: orders,
			},
			wantOrder: "rmd",
			wantVal:   "1 11:02",
		},
		{
			name: "parseInput_3",
			args: args{
				input:  "rmd+ 1 11:02",
				orders: orders,
			},
			wantOrder: "rmd+",
			wantVal:   "1 11:02",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrder, gotVal := parseInput(tt.args.input, tt.args.orders)
			if gotOrder != tt.wantOrder {
				t.Errorf("parseInput() gotOrder = %v, want %v", gotOrder, tt.wantOrder)
			}
			if gotVal != tt.wantVal {
				t.Errorf("parseInput() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
