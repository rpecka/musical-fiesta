package commands

import "testing"

func TestIsIllegal(t *testing.T) {
	if !IsIllegal("reload") {
		t.Error("the command `reload` should be illegal")
	}
}

func TestNotIllegal(t *testing.T) {
	if IsIllegal("minecraft") {
		t.Error("the command `minecraft` should be legal")
	}
}
