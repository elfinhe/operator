// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	repoRootDir string
	testDataDir string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	repoRootDir = filepath.Join(wd, "../..")
}

func TestReadLayeredYAMLs(t *testing.T) {
	testDataDir = filepath.Join(repoRootDir, "pkg/util/testdata/yaml")
	tests := []struct {
		name     string
		overlays []string
		wantErr  bool
	}{
		{
			name:     "single",
			overlays: []string{"all_off"},
			wantErr:  false,
		},
		{
			name:     "double",
			overlays: []string{"all_off", "all_on"},
			wantErr:  false,
		},
		{
			name:     "triple",
			overlays: []string{"all_off", "all_on", "sds_policy_off"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inDir := filepath.Join(testDataDir, "input")
			outPath := filepath.Join(testDataDir, "output", tt.name+".yaml")
			wantBytes, err := ioutil.ReadFile(outPath)
			want := string(wantBytes)
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v, filename: %v", err, outPath)
			}

			var filenames []string
			for _, ol := range tt.overlays {
				filenames = append(filenames, filepath.Join(inDir, ol+".yaml"))
			}
			got, err := ReadLayeredYAMLs(filenames)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLayeredYAMLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if YAMLDiff(got, want) != "" {
				t.Errorf("ReadLayeredYAMLs() got = %v, want %v", got, want)
			}
		})
	}
}
