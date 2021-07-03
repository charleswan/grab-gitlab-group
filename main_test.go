package main

import "testing"

func TestGitCloneProject(t *testing.T) {
	l := NewLimiter(5)
	if !gitCloneProject(l, "coin") {
		t.Fail()
	}
}
