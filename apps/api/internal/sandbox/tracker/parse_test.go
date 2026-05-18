package tracker

import "testing"

func TestParseInput_SubmitOnEnter(t *testing.T) {
	var buf []byte
	cmds := ParseInput(&buf, []byte("ls -la\r"))
	if len(cmds) != 1 || cmds[0] != "ls -la" {
		t.Fatalf("cmds = %v", cmds)
	}
}

func TestParseInput_Backspace(t *testing.T) {
	var buf []byte
	_ = ParseInput(&buf, []byte("lss"))
	_ = ParseInput(&buf, []byte{8})
	cmds := ParseInput(&buf, []byte("\r"))
	if len(cmds) != 1 || cmds[0] != "ls" {
		t.Fatalf("cmds = %v", cmds)
	}
}

func TestParseInput_CtrlC(t *testing.T) {
	var buf []byte
	_ = ParseInput(&buf, []byte("partial"))
	_ = ParseInput(&buf, []byte{3})
	cmds := ParseInput(&buf, []byte("ok\r"))
	if len(cmds) != 1 || cmds[0] != "ok" {
		t.Fatalf("cmds = %v", cmds)
	}
}
