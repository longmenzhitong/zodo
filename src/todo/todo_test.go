package todo

import (
	"reflect"
	"testing"
)

func Test_sortTodo(t *testing.T) {
	type args struct {
		tds []*todo
	}
	tests := []struct {
		name string
		args args
		want []*todo
	}{
		{
			name: "_sort_1",
			args: args{tds: []*todo{
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-09"},
				{Id: 2, Content: "todo_2", Status: statusPending, Deadline: "2022-11-08"},
			}},
			want: []*todo{
				{Id: 2, Content: "todo_2", Status: statusPending, Deadline: "2022-11-08"},
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-09"},
			},
		},
		{
			name: "_sort_2",
			args: args{tds: []*todo{
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-09"},
				{Id: 2, Content: "todo_2", Status: statusProcessing, Deadline: "2022-11-09"},
			}},
			want: []*todo{
				{Id: 2, Content: "todo_2", Status: statusProcessing, Deadline: "2022-11-09"},
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-09"},
			},
		},
		{
			name: "_sort_3",
			args: args{tds: []*todo{
				{Id: 2, Content: "todo_2", Status: statusPending, Deadline: "2022-11-09"},
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-09"},
			}},
			want: []*todo{
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-09"},
				{Id: 2, Content: "todo_2", Status: statusPending, Deadline: "2022-11-09"},
			},
		},
		{
			name: "_sort_4",
			args: args{tds: []*todo{
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-10"},
				{Id: 2, Content: "todo_2", Status: statusPending, Deadline: "2022-11-09"},
				{Id: 3, Content: "todo_3", Status: statusProcessing, Deadline: "2022-11-09"},
				{Id: 5, Content: "todo_5", Status: statusProcessing, Deadline: "2022-11-09"},
				{Id: 4, Content: "todo_4", Status: statusProcessing, Deadline: "2022-11-09"},
			}},
			want: []*todo{
				{Id: 3, Content: "todo_3", Status: statusProcessing, Deadline: "2022-11-09"},
				{Id: 4, Content: "todo_4", Status: statusProcessing, Deadline: "2022-11-09"},
				{Id: 5, Content: "todo_5", Status: statusProcessing, Deadline: "2022-11-09"},
				{Id: 2, Content: "todo_2", Status: statusPending, Deadline: "2022-11-09"},
				{Id: 1, Content: "todo_1", Status: statusPending, Deadline: "2022-11-10"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortTodo(tt.args.tds); !reflect.DeepEqual(got, tt.want) {
				argIds := make([]int, 0)
				for _, td := range got {
					argIds = append(argIds, td.Id)
				}
				wantIds := make([]int, 0)
				for _, td := range tt.want {
					wantIds = append(wantIds, td.Id)
				}
				t.Errorf("got = %v, want %v", argIds, wantIds)
			}
		})
	}
}
