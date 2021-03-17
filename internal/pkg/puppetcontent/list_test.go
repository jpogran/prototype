package puppetcontent

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	type args struct {
		templatePath string
		templateName string
	}
	tests := []struct {
		name    string
		args    args
		want    []ContentTemplateConfig
		wantErr bool
	}{
		{
			name: "lists correct number of templates with correct properties",
			args: args{
				templatePath: "testdata/templates",
				templateName: "",
			},
			want: []ContentTemplateConfig{
				{
					Name:            "examplegood",
					DisplayName:     "Example Good Template",
					Context:         "puppetmodule",
					Tags:            []string{"puppet"},
					TemplateType:    "project",
					TemplateURL:     "https://github.com/puppetlabs/examplegood",
					TemplateVersion: "0.1.0",
				},
			},
		},
		{
			name: "lists the specified template with correct properties",
			args: args{
				templatePath: "testdata/templates",
				templateName: "examplegood",
			},
			want: []ContentTemplateConfig{
				{
					Name:            "examplegood",
					DisplayName:     "Example Good Template",
					Context:         "puppetmodule",
					Tags:            []string{"puppet"},
					TemplateType:    "project",
					TemplateURL:     "https://github.com/puppetlabs/examplegood",
					TemplateVersion: "0.1.0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.templatePath, tt.args.templateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
