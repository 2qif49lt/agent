package daemoncmd

import (
	"fmt"
	"github.com/2qif49lt/agent/daemon"
)

// NewDaemonCommand creats a new cobra.Command for `agent daemon`
func NewDaemonCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "daemon ",
		Short: "Manage agent daemon",
		Run: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.UsageString())
		},
	}

	cmd.AddCommand(
		newInstallCommand(),
		newStartCommand(),
		newStopCommand(),
		newReStartCommand(),
		newStatusCommand(),
		newUnInstallCommand(),
	)
	return cmd
}
