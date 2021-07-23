package format

import "testing"

func TestReplaceArgs(t *testing.T) {
	type args struct {
		target string
		source map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{target: "a${b}c$d f$,${},$g a=${b}", source: map[string]string{"b": "1", "d": "2"}},
			want: "a1c2 f$,${}, a=1",
		},
		{
			name: "head",
			args: args{target: "${b}c$d f$,${},$g a=${b}", source: map[string]string{"b": "1", "d": "2"}},
			want: "1c2 f$,${}, a=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceArgs(tt.args.target, tt.args.source); got != tt.want {
				t.Errorf("ReplaceArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
