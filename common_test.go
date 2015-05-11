package grim

import (
	"io/ioutil"
	"os"
	"testing"
)

const badPath = "/\\/\\/\\"
const badContents = "}{"

func getEnvOrSkip(t *testing.T, name string) string {
	value := os.Getenv(name)

	if value == "" {
		t.Skipf("this test requires the environment variable %q to be set", name)
	}

	return value
}

func withTempFile(t *testing.T, contents string, f func(string)) {
	withTempFilePerms(t, contents, 0660, f)
}

func withTempScript(t *testing.T, contents string, f func(string)) {
	withTempFilePerms(t, contents, 0770, f)
}

func withTempFilePerms(t *testing.T, contents string, mode os.FileMode, f func(string)) {
	file, err := ioutil.TempFile("", "grim_testing_")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	file.Close()

	fn := file.Name()

	defer func() {
		if err := os.Remove(fn); err != nil {
			t.Errorf("failed to remove temp file at %v: %v", fn, err)
		}
	}()

	if err := ioutil.WriteFile(fn, []byte(contents), mode); err != nil {
		t.Fatalf("failed to write temp file contents at %q: %v", fn, err)
	}

	os.Chmod(fn, mode)
	f(fn)
}

func withTempDir(t *testing.T, f func(string)) {
	dir, err := ioutil.TempDir("", "grim_testing_")

	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("failed to remove temp dir at %v: %v", dir, err)
		}
	}()

	f(dir)
}
