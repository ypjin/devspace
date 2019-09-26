// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devspace-cloud/devspace/cmd"
	"github.com/spf13/cobra/doc"
)

const cliDocsDir = "./docs/pages/cli/commands"
const headerTemplate = `---
title: "%s"
sidebar_label: %s
---

`

var fixSynopsisRegexp = regexp.MustCompile("(?smi)(## devspace.*?\n)(.*?)#(## Synopsis\n*\\s*)(.*?)(\\s*\n\n)((```)(.*?))?#(## Options)(.*?)#(## SEE ALSO)(\\s*\\* \\[devspace\\][^\n]*)?(.*)\n###### Auto generated by spf13/cobra on .*$")

// Run executes the command logic
func main() {
	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		command := strings.Split(base, "_")
		title := strings.Join(command, " ")
		sidebarLabel := title
		l := len(command)

		if l > 2 {
			sidebarLabel = command[l-1]
		}

		return fmt.Sprintf(headerTemplate, "Command - "+title, sidebarLabel)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "../../cli/commands/" + strings.ToLower(base)
	}

	rootCmd := cmd.GetRoot()

	err := doc.GenMarkdownTreeCustom(rootCmd, cliDocsDir, filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(cliDocsDir, func(path string, info os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if stat.IsDir() {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := fixSynopsisRegexp.ReplaceAllString(string(content), "$2$3$7$8```\n$4\n```\n$9$10## See Also$13")

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
