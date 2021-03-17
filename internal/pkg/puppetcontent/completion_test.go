package puppetcontent

import (
	"reflect"
	"testing"
)

func TestCompleteName(t *testing.T) {
	type args struct {
		templatePath string
		match        string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "accurate completion from single letter",
			args: args{
				templatePath: "testdata/templates",
				match:        "e",
			},
			want: []string{
				"examplegood" + "\t" + "Example Good Template",
			},
		},
		{
			name: "accurate completion from multiple letters",
			args: args{
				templatePath: "testdata/templates",
				match:        "example",
			},
			want: []string{
				"examplegood" + "\t" + "Example Good Template",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompleteName(tt.args.templatePath, tt.args.match); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompleteName() = %v, want %v", got, tt.want)
			}
		})
	}
}
