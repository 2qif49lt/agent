package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	flag "github.com/2qif49lt/pflag"
)

// Command is the struct containing the command name and description
type Command struct {
	Name        string
	Description string
}

type byName []Command

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// TODO(tiborvass): do not show 'daemon' on client-only binaries

func SortCommands(commands []Command) []Command {
	tmp := make([]Command, len(commands))
	copy(tmp, commands)
	sort.Sort(byName(tmp))
	return tmp
}

// Cli represents a command line interface.
type Cli struct {
	Stdout   io.Writer
	handlers []Handler
	Usage    func()
}

// Handler holds the different commands Cli will call
// It should have methods with names starting with `Cmd` like:
// 	func (h yHandler) CmdFoo(args ...string) error
type Handler interface {
	Command(name string) func(...string) error
}

// Initializer can be optionally implemented by a Handler to
// initialize before each call to one of its commands.
type Initializer interface {
	Initialize() error
}

// New instantiates a ready-to-use Cli.
func New(handlers ...Handler) *Cli {
	// make the generic Cli object the first cli handler
	// in order to handle `docker help` appropriately
	cli := new(Cli)
	cli.handlers = append([]Handler{cli}, handlers...)
	return cli
}

var errCommandNotFound = errors.New("command not found")

func (cli *Cli) command(args ...string) (func(...string) error, error) {
	for _, c := range cli.handlers {
		if c == nil {
			continue
		}
		if cmd := c.Command(strings.Join(args, " ")); cmd != nil {
			if ci, ok := c.(Initializer); ok {
				if err := ci.Initialize(); err != nil {
					return nil, err
				}
			}
			return cmd, nil
		}
	}
	return nil, errCommandNotFound
}

// Run executes the specified command.
func (cli *Cli) Run(args ...string) error {
	if len(args) > 1 {
		command, err := cli.command(args[:2]...)
		if err == nil {
			return command(args[2:]...)
		}
		if err != errCommandNotFound {
			return err
		}
	}
	if len(args) > 0 {
		command, err := cli.command(args[0])
		if err != nil {
			if err == errCommandNotFound {
				cli.noSuchCommand(args[0])
				return nil
			}
			return err
		}
		return command(args[1:]...)
	}
	return cli.CmdHelp()
}

func (cli *Cli) noSuchCommand(command string) {
	if cli.Stdout == nil {
		cli.Stdout = os.Stdout
	}
	fmt.Fprintf(cli.Stdout, "agent: '%s' is not a agent command.\nSee 'agent --help'.\n", command)
	os.Exit(1)
}

// Command returns a command handler, or nil if the command does not exist
func (cli *Cli) Command(name string) func(...string) error {
	return map[string]func(...string) error{
		"help": cli.CmdHelp,
	}[name]
}

// CmdHelp displays information on a agent command.
//
// If more than one command is specified, information is only shown for the first command.
//
// Usage: agent help COMMAND or agent COMMAND --help
func (cli *Cli) CmdHelp(args ...string) error {
	if len(args) > 1 {
		command, err := cli.command(args[:2]...)
		if err == nil {
			command("--help")
			return nil
		}
		if err != errCommandNotFound {
			return err
		}
	}
	if len(args) > 0 {
		command, err := cli.command(args[0])
		if err != nil {
			if err == errCommandNotFound {
				cli.noSuchCommand(args[0])
				return nil
			}
			return err
		}
		command("--help")
		return nil
	}

	if cli.Usage == nil {
		cli.Usage()
	}

	return nil
}

// Subcmd is a subcommand of the main "agent" command.
// A subcommand represents an action that can be performed
// from the agent command line client.
//
// To see all available subcommands, run "agent --help".
func Subcmd(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	return flags
}

// StatusError reports an unsuccessful exit by a command.
type StatusError struct {
	Status     string
	StatusCode int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("Status: %s, Code: %d", e.Status, e.StatusCode)
}
