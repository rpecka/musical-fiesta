package configfile

import (
	"fmt"
	"os"
	"strings"
)

type configfileWriter struct {
	file *os.File
}

func makeExecCommand(cfgPath string) string {
	return "exec " + cfgPath
}

func makeAliasCommand(alias, command string) string {
	return fmt.Sprintf("alias %s %s", alias, command)
}

func makeBindCommand(bindKey, command string) string {
	return fmt.Sprintf("bind %s %s", bindKey, command)
}

func makeUnbindCommand(bindKey string) string {
	return fmt.Sprintf("unbind %s", bindKey)
}

func makeEcho(s string) string {
	return fmt.Sprintf("echo %s", s)
}

func boolToIntString(b bool) string {
	if b {
		return "1"
	} else {
		return "0"
	}
}

func boolToPlusMinus(b bool) string {
	if b {
		return "+"
	} else {
		return "-"
	}
}

func makeVoiceInputFromFile(enabled bool) string {
	return "voice_inputfromfile " + boolToIntString(enabled)
}

func makeVoiceLoopBack(enabled bool) string {
	return "voice_loopback " + boolToIntString(enabled)
}

func makeVoiceRecord(enabled bool) string {
	return boolToPlusMinus(enabled) + "voicerecord"
}

func makeHostWriteconfig(filename string) string {
	return "host_writeconfig " + filename
}

func chainCommands(commands []string) string {
	return fmt.Sprintf("\"%s\"", strings.Join(commands, "; "))
}

func newWriter(path string) (*configfileWriter, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0755))
	if err != nil {
		return nil, err
	}
	return &configfileWriter{file: file}, nil
}

func (c configfileWriter) close() error {
	return c.file.Close()
}

func (c configfileWriter) writeLine(s string) error {
	_, err := c.file.WriteString(s + "\n")
	return err
}

func (c configfileWriter) writeComment(comment string) error {
	return c.writeLine("// " + comment)
}

func (c configfileWriter) writeHeader(title string) error {
	return c.writeComment(fmt.Sprintf("------------------%s------------------", title))
}

func (c configfileWriter) writeAlias(alias, command string) error {
	return c.writeLine(makeAliasCommand(alias, command))
}

func (c configfileWriter) writeBind(bindKey, command string) error {
	return c.writeLine(makeBindCommand(bindKey, command))
}

func (c configfileWriter) writeEcho(s string) error {
	return c.writeLine(makeEcho(s))
}

func (c configfileWriter) writeEchoHeader(title string) error {
	return c.writeEcho(fmt.Sprintf("------------------%s------------------", title))
}
