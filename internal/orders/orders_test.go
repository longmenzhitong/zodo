package orders

import "testing"

func Test_parseStr(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name        string
		args        args
		wantContent string
		wantErr     bool
	}{
		{
			name:        "parseStr_1",
			args:        args{input: "add "},
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "parseStr_2",
			args:        args{input: "add a"},
			wantContent: "a",
			wantErr:     false,
		},
		{
			name:        "parseStr_3",
			args:        args{input: "add a b"},
			wantContent: "a b",
			wantErr:     false,
		},
		{
			name:        "parseStr_4",
			args:        args{input: "add  a b "},
			wantContent: "a b",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := parseStr(tt.args.input, prefixAdd)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContent != tt.wantContent {
				t.Errorf("parseStr() gotContent = %v, want %v", gotContent, tt.wantContent)
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
			args:        args{input: "mod "},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "parseIdAndStr_2",
			args:        args{input: "mod 1"},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "parseIdAndStr_3",
			args:        args{input: "mod a"},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "parseIdAndStr_4",
			args:        args{input: "mod 1 1"},
			wantId:      1,
			wantContent: "1",
			wantErr:     false,
		},
		{
			name:        "parseIdAndStr_5",
			args:        args{input: "mod 1 1 a"},
			wantId:      1,
			wantContent: "1 a",
			wantErr:     false,
		},
		{
			name:        "parseIdAndStr_6",
			args:        args{input: "mod 1  1 a  "},
			wantId:      1,
			wantContent: "1 a",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotContent, err := parseIdAndStr(tt.args.input, prefixModify)
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
			args:         args{input: "ddl "},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_2",
			args:         args{input: "ddl 1"},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_3",
			args:         args{input: "ddl 2022-10-19"},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_4",
			args:         args{input: "ddl 1 2022-10"},
			wantId:       1,
			wantDeadline: "2022-10",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_5",
			args:         args{input: "ddl 1 2022-10-19"},
			wantId:       1,
			wantDeadline: "2022-10-19",
			wantErr:      false,
		},
		{
			name:         "parseDeadline_6",
			args:         args{input: "ddl  1  2022-10-19 "},
			wantId:       1,
			wantDeadline: "2022-10-19",
			wantErr:      false,
		},
		{
			name:         "parseDeadline_7",
			args:         args{input: "ddl  1  2022-10-19 a"},
			wantId:       1,
			wantDeadline: "2022-10-19 a",
			wantErr:      true,
		},
		{
			name:         "parseDeadline_8",
			args:         args{input: "ddl 1 10-19"},
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
