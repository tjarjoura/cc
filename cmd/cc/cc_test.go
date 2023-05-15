package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
)

func dumpFiles(t *testing.T, files ...string) {
	for _, f := range files {
		contents, err := os.ReadFile(f)
		if err != nil {
			t.Errorf("failed to read %s: %s", f, err)
		} else {
			t.Logf("%s:\n%s", f, contents)
		}
	}
}

// Compile and check the status codes of all test programs
func TestCC(t *testing.T) {
	log.SetFlags(0)

	_, thisFile, _, _ := runtime.Caller(0)
	testDir := path.Dir(thisFile)

	testPrograms := fmt.Sprintf("%s/test_programs/", testDir)

	t.Logf("test programs are in %s", testPrograms)
	sourceFiles, err := ioutil.ReadDir(testPrograms)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range sourceFiles {
		t.Run(f.Name(), func(t *testing.T) {
			fullPath := path.Join(testPrograms, f.Name())
			var asmFiles, objFiles []string
			outFile := strings.ReplaceAll(fullPath, ".c", "")
			defer func() {
				allFiles := append(asmFiles, outFile)
				allFiles = append(allFiles, objFiles...)
				for _, f := range allFiles {
					os.Remove(f)
				}
			}()

			asmFiles, err := compile("",
				path.Join(testPrograms, f.Name()))
			if err != nil {
				t.Fatalf("error compiling %s: %s", f.Name(), err)
			}

			objFiles, err = assemble("", asmFiles...)
			if err != nil {
				t.Errorf("error assembling %s: %s", f.Name(), err)
				dumpFiles(t, asmFiles...)
				t.FailNow()
			}

			if err := link(outFile, objFiles...); err != nil {
				t.Fatalf("error linking %s: %s", outFile, err)
			}

			cmd := exec.Command(outFile)
			if err := cmd.Run(); err != nil {
				t.Errorf("error running %s: %s", outFile, err)
				dumpFiles(t, asmFiles...)
				t.FailNow()
			}
		})
	}
}
