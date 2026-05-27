package main

import "testing"

func TestIsVersionArg(t *testing.T) {
	for _, arg := range []string{"-v", "--version", "-version"} {
		if !isVersionArg(arg) {
			t.Fatalf("expected %q to be a version arg", arg)
		}
	}
	for _, arg := range []string{"version", "-V", "--help", "convert"} {
		if isVersionArg(arg) {
			t.Fatalf("did not expect %q to be a version arg", arg)
		}
	}
}
