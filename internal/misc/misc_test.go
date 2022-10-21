package misc

import (
	"testing"
	"time"
	"zodo/internal/cst"
)

func TestCalcBetweenDays(t *testing.T) {
	type args struct {
		t1 time.Time
		t2 time.Time
	}
	tests := []struct {
		name           string
		args           args
		wantNatureDays int
		wantWorkDays   int
	}{
		{
			name: "CalcBetweenDays_1",
			args: args{
				t1: parseTime(cst.LayoutDate, "2022-10-21"),
				t2: parseTime(cst.LayoutDate, "2022-10-21"),
			},
			wantNatureDays: 0,
			wantWorkDays:   0,
		},
		{
			name: "CalcBetweenDays_2",
			args: args{
				t1: parseTime(cst.LayoutDate, "2022-10-21"),
				t2: parseTime(cst.LayoutDate, "2022-10-22"),
			},
			wantNatureDays: 1,
			wantWorkDays:   1,
		},
		{
			name: "CalcBetweenDays_3",
			args: args{
				t1: parseTime(cst.LayoutDate, "2022-10-21"),
				t2: parseTime(cst.LayoutDate, "2022-10-27"),
			},
			wantNatureDays: 6,
			wantWorkDays:   4,
		},
		{
			name: "CalcBetweenDays_4",
			args: args{
				t1: parseTime(cst.LayoutDate, "2022-10-22"),
				t2: parseTime(cst.LayoutDate, "2022-10-21"),
			},
			wantNatureDays: -1,
			wantWorkDays:   -1,
		},
		{
			name: "CalcBetweenDays_5",
			args: args{
				t1: parseTime(cst.LayoutDate, "2022-10-27"),
				t2: parseTime(cst.LayoutDate, "2022-10-21"),
			},
			wantNatureDays: -6,
			wantWorkDays:   -4,
		},
		{
			name: "CalcBetweenDays_6",
			args: args{
				t1: parseTime(cst.LayoutDateTime, "2022-10-21 11:18:00"),
				t2: parseTime(cst.LayoutDate, "2022-10-27"),
			},
			wantNatureDays: 6,
			wantWorkDays:   4,
		},
		{
			name: "CalcBetweenDays_7",
			args: args{
				t1: parseTime(cst.LayoutDate, "2022-10-21"),
				t2: parseTime(cst.LayoutDateTime, "2022-10-27 11:18:00"),
			},
			wantNatureDays: 6,
			wantWorkDays:   4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNatureDays, gotWorkDays := CalcBetweenDays(tt.args.t1, tt.args.t2)
			if gotNatureDays != tt.wantNatureDays {
				t.Errorf("CalcBetweenDays() gotNatureDays = %v, want %v", gotNatureDays, tt.wantNatureDays)
			}
			if gotWorkDays != tt.wantWorkDays {
				t.Errorf("CalcBetweenDays() gotWorkDays = %v, want %v", gotWorkDays, tt.wantWorkDays)
			}
		})
	}
}

func parseTime(layout, val string) time.Time {
	res, err := time.Parse(layout, val)
	if err != nil {
		panic(err)
	}
	return res
}
