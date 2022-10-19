package todo

import (
	"testing"
	"time"
	"zodo/internal/cst"
)

func Test_calcRemainDays(t *testing.T) {
	type args struct {
		calcTime time.Time
		deadline string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "calcRemainDays_1",
			args: args{
				calcTime: getCalcTime("2020-10-20"),
				deadline: "2020-10-20",
			},
			want: 0,
		},
		{
			name: "calcRemainDays_2",
			args: args{
				calcTime: getCalcTime("2020-10-20"),
				deadline: "2020-10-21",
			},
			want: 1,
		},
		{
			name: "calcRemainDays_3",
			args: args{
				calcTime: getCalcTime("2020-10-20"),
				deadline: "2020-10-19",
			},
			want: -1,
		},
		{
			name: "calcRemainDays_4",
			args: args{
				calcTime: getCalcTime("2020-10-20"),
				deadline: "2020-10-27",
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcRemainDays(tt.args.calcTime, tt.args.deadline); got != tt.want {
				t.Errorf("calcRemainDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getCalcTime(v string) time.Time {
	ct, err := time.Parse(cst.LayoutDate, v)
	if err != nil {
		panic(err)
	}
	return ct
}
