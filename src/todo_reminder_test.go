package zodo

import (
	"testing"
	"time"
)

func Test_isNeedRemind(t *testing.T) {
	type args struct {
		rmdTime   string
		loopType  string
		rmdStatus string
		checkTime time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "isNeedRemind_1",
			args: args{
				rmdTime:   "2022-11-06 19:28",
				loopType:  loopOnce,
				rmdStatus: rmdStatusWaiting,
				checkTime: getCheckTime("2022-11-06 19:28:00"),
			},
			want: true,
		},
		{
			name: "isNeedRemind_2",
			args: args{
				rmdTime:   "2022-11-06 19:28",
				loopType:  loopOnce,
				rmdStatus: rmdStatusFinished,
				checkTime: getCheckTime("2022-11-06 19:28:00"),
			},
			want: false,
		},
		{
			name: "isNeedRemind_3",
			args: args{
				rmdTime:   "19:28",
				loopType:  loopPerDay,
				rmdStatus: rmdStatusWaiting,
				checkTime: getCheckTime("2022-11-06 19:28:00"),
			},
			want: true,
		},
		{
			name: "isNeedRemind_4",
			args: args{
				rmdTime:   "19:28",
				loopType:  loopPerWorkDay,
				rmdStatus: rmdStatusWaiting,
				checkTime: getCheckTime("2022-11-06 19:28:00"),
			},
			want: false,
		},
		{
			name: "isNeedRemind_5",
			args: args{
				rmdTime:   "19:28",
				loopType:  loopPerWorkDay,
				rmdStatus: rmdStatusWaiting,
				checkTime: getCheckTime("2022-11-07 19:28:00"),
			},
			want: true,
		},
		{
			name: "isNeedRemind_6",
			args: args{
				rmdTime:   "19:28",
				loopType:  loopPerWorkDay,
				rmdStatus: rmdStatusWaiting,
				checkTime: getCheckTime("2022-11-07 19:27:00"),
			},
			want: false,
		},
		{
			name: "isNeedRemind_6",
			args: args{
				rmdTime:   "19:28",
				loopType:  loopPerSunday,
				rmdStatus: rmdStatusWaiting,
				checkTime: getCheckTime("2022-11-06 19:28:00"),
			},
			want: true,
		},
		{
			name: "isNeedRemind_7",
			args: args{
				rmdTime:   "19:28",
				loopType:  loopPerSunday,
				rmdStatus: rmdStatusWaiting,
				checkTime: getCheckTime("2022-11-07 19:28:00"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNeedRemind(tt.args.rmdTime, tt.args.loopType, tt.args.rmdStatus, tt.args.checkTime); got != tt.want {
				t.Errorf("isNeedRemind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getCheckTime(t string) time.Time {
	ct, err := time.ParseInLocation(LayoutDateTime, t, time.Local)
	if err != nil {
		panic(err)
	}
	return ct
}
