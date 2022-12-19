// Copyright 2022 Giuseppe De Palma, Matteo Trentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test for the copy function
func Test_copy(t *testing.T) {
	createSource := func(t *testing.T, numberOfFiles, mode int) (string, error) {
		t.Helper()
		data := []byte("a funless function")

		// create a folder for source files
		srcDir, err := os.MkdirTemp("", "funless-test-src-")
		require.NoError(t, err)

		// create n files inside the created folder
		for i := 1; i <= numberOfFiles; i++ {
			srcFile := fmt.Sprintf("%s/test-file-%d", srcDir, i)
			err = os.WriteFile(srcFile, data, os.FileMode(mode))
			require.NoError(t, err)
		}

		return srcDir, nil
	}

	fileModes := []int{0600, 0640, 0644, 0700, 0755}
	numberOfFiles := 2

	for _, mode := range fileModes {
		// set up source with 2 files
		srcDir, _ := createSource(t, numberOfFiles, mode)
		defer os.RemoveAll(srcDir)

		// create a destination folder to copy the files to
		destDir, destDirErr := os.MkdirTemp("", "funless-test-dest-")
		if destDirErr != nil {
			t.Fatalf("Error creating destination folder\n%v", destDirErr)
		}
		defer os.RemoveAll(destDir)

		err := Copy(srcDir, destDir)
		require.NoError(t, err)

		for i := 1; i <= numberOfFiles; i++ {
			fileInfo, err := os.Stat(fmt.Sprintf("%s/test-file-%d", destDir, i))
			require.NoError(t, err)
			require.Equal(t, os.FileMode(mode), fileInfo.Mode())
		}
	}
}
