package puppetcontent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDeploy(t *testing.T) {
	// tempDir := t.TempDir()
	tempDir, _ := ioutil.TempDir("", "testdeploy")
	type args struct {
		SelectedTemplate   string
		LocalTemplateCache string
		TargetOutput       string
		TargetName         string
	}
	tests := []struct {
		name    string
		tempdir string
		args    args
		want    []string
		want1   []error
	}{
		{
			name: "deploys a project with correct number of files",
			args: args{
				SelectedTemplate:   "puppet-module",
				LocalTemplateCache: "testdata/templates/exampledeploytest",
				TargetOutput:       filepath.Join(tempDir, "foo"),
				TargetName:         "foo",
			},
			want: []string{
				filepath.Join(tempDir, "foo"),
				filepath.Join(tempDir, "foo/.vscode"),
				filepath.Join(tempDir, "foo/.vscode/extensions.json"),
				filepath.Join(tempDir, "foo/Gemfile"),
				filepath.Join(tempDir, "foo/README.md"),
				filepath.Join(tempDir, "foo/Rakefile"),
				filepath.Join(tempDir, "foo/data"),
				filepath.Join(tempDir, "foo/data/common.yaml"),
				filepath.Join(tempDir, "foo/examples"),
				filepath.Join(tempDir, "foo/examples/.gitkeep"),
				filepath.Join(tempDir, "foo/files"),
				filepath.Join(tempDir, "foo/files/.gitkeep"),
				filepath.Join(tempDir, "foo/hiera.yaml"),
				filepath.Join(tempDir, "foo/manifests"),
				filepath.Join(tempDir, "foo/manifests/.gitkeep"),
				filepath.Join(tempDir, "foo/metadata.json"),
				filepath.Join(tempDir, "foo/spec"),
				filepath.Join(tempDir, "foo/spec/default_facts.yml"),
				filepath.Join(tempDir, "foo/spec/spec_helper.rb"),
				filepath.Join(tempDir, "foo/tasks"),
				filepath.Join(tempDir, "foo/tasks/.gitkeep"),
				filepath.Join(tempDir, "foo/templates"),
				filepath.Join(tempDir, "foo/templates/.gitkeep"),
			},
			want1: nil,
		},
		{
			name: "deploys an item with correct number of files",
			args: args{
				SelectedTemplate:   "puppet-class",
				LocalTemplateCache: "testdata/templates/exampledeploytest",
				TargetOutput:       filepath.Join(tempDir, "foo"),
				TargetName:         "foo",
			},
			want: []string{
				tempDir,
				filepath.Join(tempDir, "manifests"),
				filepath.Join(tempDir, "manifests/foo.pp"),
				filepath.Join(tempDir, "spec"),
				filepath.Join(tempDir, "spec/classes"),
				filepath.Join(tempDir, "spec/classes/foospec.rb"),
			},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := Deploy(tt.args.SelectedTemplate, tt.args.LocalTemplateCache, tt.args.TargetOutput, tt.args.TargetName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deploy() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Deploy() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	defer os.RemoveAll(tempDir)
}

func TestDeployDual(t *testing.T) {
	// tempDir := t.TempDir()
	tempDir, _ := ioutil.TempDir("", "testdeploy")

	t.Run("test project and item deploy with name and output", func(t *testing.T) {
		_, _ = Deploy(
			"puppet-module",
			"testdata/templates/exampledeploytest",
			filepath.Join(tempDir, "foo"),
			"foo",
		)
		items, _ := Deploy(
			"puppet-class",
			"testdata/templates/exampledeploytest",
			filepath.Join(tempDir, "foo"),
			"foo",
		)
		itemWant := []string{
			tempDir,
			filepath.Join(tempDir, "manifests"),
			filepath.Join(tempDir, "manifests/foo.pp"),
			filepath.Join(tempDir, "spec"),
			filepath.Join(tempDir, "spec/classes"),
			filepath.Join(tempDir, "spec/classes/foospec.rb"),
		}
		if !reflect.DeepEqual(items, itemWant) {
			t.Errorf("Deploy() got = %v, want %v", items, itemWant)
		}
		defer os.RemoveAll(tempDir)
	})
}
