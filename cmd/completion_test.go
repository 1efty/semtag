package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestZshCompletion(t *testing.T) {
	cmd := completionCmd
	b := bytes.NewBufferString("")
	cmd.SetArgs([]string{""})
	cmd.Execute()
	_, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
}
