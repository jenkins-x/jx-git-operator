package cmd

import (
	"fmt"
	"os"

	"github.com/jenkins-x/jx-git-operator/pkg/poller"

	"github.com/spf13/cobra"
)

// NewMain creates a command object
func NewMain() (*cobra.Command, *poller.Options) {
	o := &poller.Options{}

	cmd := &cobra.Command{
		Short: "Starts a pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&o.Dir, "dir", "d", "", "the work directory where git clones take place. If none is specified a temporary dir is used")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "", "", "the namespace to use")
	return cmd, o
}
