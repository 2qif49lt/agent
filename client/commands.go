package client

// Command returns a cli command handler if one exists
func (cli *AgentCli) Command(name string) func(...string) error {
	return map[string]func(...string) error{
		"info": cli.CmdInfo,
		"ping": cli.CmdPing,
	}[name]
}
