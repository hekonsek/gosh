package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	tmpScriptDir, err := ioutil.TempDir("", "gosh_")
	if err != nil {
		panic(err)
	}

	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	scriptPath := filepath.Join(pwd, os.Args[1])

	scriptBytes, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		panic(err)
	}

	targetScriptPath := filepath.Join(tmpScriptDir, filepath.Base(os.Args[1]))
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

	scriptExec := exec.Command("go", "run", filepath.Base(os.Args[1]))
	scriptExec.Dir = tmpScriptDir
	scriptExec.Env = append(os.Environ(), "GO111MODULE=on")
	out, err := scriptExec.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		panic(err)
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
