# gosh - simple Go scripting

**gosh** is a simple tool for executing Go scripts. Its primary feature is an ability read Go Modules metadata from the
script itself, so you can be sure that versions of your script dependencies are fixed.

## Usage

Define your script as simple application with `main()` function. If you would like to be sure that version of dependency 
you are relying on is fixed, just add it as as comment on the beginning of the script using standard Go Modules syntax.

For example this is how you can ensure that your script uses `github.com/fatih/color` dependency in version `1.7.0`: 

```
// require github.com/fatih/color v1.7.0
package main

import "github.com/fatih/color"

func main() {
	color.Green("Hello world!\n")
}
```

Now you can execute the script:

```
$ gosh script.go
Hello world!
```
