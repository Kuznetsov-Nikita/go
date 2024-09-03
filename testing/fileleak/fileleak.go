//go:build !solution

package fileleak

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func VerifyNone(t testingT) {
	startFiles := make(map[string]int)
	files, _ := os.ReadDir("/proc/self/fd")
	for _, file := range files {
		fd := file.Name()
		link, _ := os.Readlink(filepath.Join("/proc/self/fd", fd))
		if _, ok := startFiles[link]; !ok {
			startFiles[link] = 1
		}
		startFiles[link]++
	}

	t.Cleanup(func() {
		endFiles := make(map[string]int)
		files, _ := os.ReadDir("/proc/self/fd")
		for _, file := range files {
			fd := file.Name()
			link, _ := os.Readlink(filepath.Join("/proc/self/fd", fd))
			if _, ok := endFiles[link]; !ok {
				endFiles[link] = 1
			}
			endFiles[link]++
		}

		var leakedFiles []string
		for file, cnt := range endFiles {
			if _, ok := startFiles[file]; !ok {
				leakedFiles = append(leakedFiles, file)
			} else {
				if cnt > startFiles[file] {
					leakedFiles = append(leakedFiles, file)
				}
			}
		}

		if len(leakedFiles) > 0 {
			sort.Strings(leakedFiles)
			t.Errorf("File leaks detected: %s", strings.Join(leakedFiles, ", "))
		}
	})
}
