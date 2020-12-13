package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestEmptyRun(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetArgs([]string{""})
	cmd.Execute()
	_, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
}
