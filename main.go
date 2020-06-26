package main

import (
	"fmt"
	"os"

	"github.com/jenkins-x/jx-git-operator/pkg/cmd"
)

func main() {
	c, _ := cmd.NewMain()
	err := c.Execute()
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		if err != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
