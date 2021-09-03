package filter

import (
	"reflect"
	"strconv"
	"testing"
)

func TestFilterStringSlice(t *testing.T) {
	type args struct {
		source    []string
		condition func(val string) bool
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []string
	}{
		{
			name: "success",
			args: args{
				source: []string{"1", "2", "3", "a", "b", "c"},
				condition: func(val string) bool {
					_, err := strconv.Atoi(val) //数字过掉

					return err == nil
				},
			},
			want:  []string{"1", "2", "3"},
			want1: []string{"a", "b", "c"}},
		{
			name: "no_condition",
			args: args{
				source: []string{"1", "2", "3", "a", "b", "c"},
			},
			want:  []string{},
			want1: []string{"1", "2", "3", "a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StringSlice(tt.args.source, tt.args.condition)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSlice() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("StringSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
