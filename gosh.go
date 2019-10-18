package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var RootCommand = &cobra.Command{
	Use:   "gosh script",
	Short: "Gosh - simple Go scripting",

	Run: func(cmd *cobra.Command, args []string) {
		showHelpIfNeeded(cmd)

		scriptArg := os.Args[1]

		tmpScriptDir, err := ioutil.TempDir("", "gosh_")
		if err != nil {
			panic(err)
		}

		pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		scriptPath := filepath.Join(pwd, scriptArg)

		scriptBytes, err := ioutil.ReadFile(scriptPath)
		if err != nil {
			panic(err)
		}

		targetScriptPath := filepath.Join(tmpScriptDir, filepath.Base(scriptArg))
		err = ioutil.WriteFile(targetScriptPath, scriptBytes, 0644)
		if err != nil {
			panic(err)
		}

		dependencies := parseDependencies(string(scriptBytes))

		targetModPath := filepath.Join(tmpScriptDir, "go.mod")
		targetModBytes := []byte(`module gosh

go 1.12

` + strings.Join(dependencies, "\n"))
		err = ioutil.WriteFile(targetModPath, targetModBytes, 0644)
		if err != nil {
			panic(err)
		}

		scriptExec := exec.Command("go", "run", filepath.Base(scriptArg))
		scriptExec.Dir = tmpScriptDir
		scriptExec.Env = append(os.Environ(), "GO111MODULE=on")
		out, err := scriptExec.CombinedOutput()
		fmt.Print(string(out))
		if err != nil {
			panic(err)
		}
	},
}

func main() {
	exitOnCliError(RootCommand.Execute())
}

func showHelpIfNeeded(cmd *cobra.Command) {
	if len(os.Args) != 2 {
		exitOnCliError(cmd.Help())
		os.Exit(UnixExitCodeGeneralError)
	}
}

func parseDependencies(script string) []string {
	var dependencies []string
	for _, line := range strings.Split(script, "\n") {
		if strings.HasPrefix(line, "//") {
			dependencies = append(dependencies, strings.Replace(line, "//", "", 1))
		}
		if strings.HasPrefix(line, "package") {
			break
		}
	}
	return dependencies
}

// Util

const UnixExitCodeGeneralError = 1

func cliError(err error) bool {
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return true
	}
	return false
}

func exitOnCliError(err error) {
	if cliError(err) {
		os.Exit(UnixExitCodeGeneralError)
	}
}
