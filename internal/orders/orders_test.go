package orders

import "testing"

func TestParseAdd(t *testing.T) {
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
			name:        "ParseAdd_1",
			args:        args{input: "add "},
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "ParseAdd_2",
			args:        args{input: "add a"},
			wantContent: "a",
			wantErr:     false,
		},
		{
			name:        "ParseAdd_3",
			args:        args{input: "add a b"},
			wantContent: "a b",
			wantErr:     false,
		},
		{
			name:        "ParseAdd_4",
			args:        args{input: "add  a b "},
			wantContent: "a b",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := ParseAdd(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContent != tt.wantContent {
				t.Errorf("ParseAdd() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func TestParseModify(t *testing.T) {
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
			name:        "ParseModify_1",
			args:        args{input: "mod "},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "ParseModify_2",
			args:        args{input: "mod 1"},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "ParseModify_3",
			args:        args{input: "mod a"},
			wantId:      0,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "ParseModify_4",
			args:        args{input: "mod 1 1"},
			wantId:      1,
			wantContent: "1",
			wantErr:     false,
		},
		{
			name:        "ParseModify_5",
			args:        args{input: "mod 1 1 a"},
			wantId:      1,
			wantContent: "1 a",
			wantErr:     false,
		},
		{
			name:        "ParseModify_6",
			args:        args{input: "mod 1  1 a  "},
			wantId:      1,
			wantContent: "1 a",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotContent, err := ParseModify(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseModify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("ParseModify() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotContent != tt.wantContent {
				t.Errorf("ParseModify() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func TestParsePending(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantId  int
		wantErr bool
	}{
		{
			name:    "ParsePending_1",
			args:    args{input: "pending "},
			wantId:  0,
			wantErr: true,
		},
		{
			name:    "ParsePending_2",
			args:    args{input: "pending a"},
			wantId:  0,
			wantErr: true,
		},
		{
			name:    "ParsePending_3",
			args:    args{input: "pending 1"},
			wantId:  1,
			wantErr: false,
		},
		{
			name:    "ParsePending_4",
			args:    args{input: "pending  1 "},
			wantId:  1,
			wantErr: false,
		},
		{
			name:    "ParsePending_5",
			args:    args{input: "pending  1 a"},
			wantId:  0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, err := ParsePending(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePending() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("ParsePending() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestParseDeadline(t *testing.T) {
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
			name:         "ParseDeadline_1",
			args:         args{input: "ddl "},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "ParseDeadline_2",
			args:         args{input: "ddl 1"},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "ParseDeadline_3",
			args:         args{input: "ddl 2022-10-19"},
			wantId:       0,
			wantDeadline: "",
			wantErr:      true,
		},
		{
			name:         "ParseDeadline_4",
			args:         args{input: "ddl 1 2022-10"},
			wantId:       1,
			wantDeadline: "2022-10",
			wantErr:      true,
		},
		{
			name:         "ParseDeadline_5",
			args:         args{input: "ddl 1 2022-10-19"},
			wantId:       1,
			wantDeadline: "2022-10-19",
			wantErr:      false,
		},
		{
			name:         "ParseDeadline_6",
			args:         args{input: "ddl  1  2022-10-19 "},
			wantId:       1,
			wantDeadline: "2022-10-19",
			wantErr:      false,
		},
		{
			name:         "ParseDeadline_6",
			args:         args{input: "ddl  1  2022-10-19 a"},
			wantId:       1,
			wantDeadline: "2022-10-19 a",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotDeadline, err := ParseDeadline(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDeadline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("ParseDeadline() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotDeadline != tt.wantDeadline {
				t.Errorf("ParseDeadline() gotDeadline = %v, want %v", gotDeadline, tt.wantDeadline)
			}
		})
	}
}
