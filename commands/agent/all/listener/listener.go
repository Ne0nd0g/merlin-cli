/*
Merlin is a post-exploitation command and control framework.

This file is part of Merlin.
Copyright (C) 2024 Russel Van Tuyl

Merlin is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
any later version.

Merlin is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Merlin.  If not, see <http://www.gnu.org/licenses/>.
*/

package listener

import (
	// Standard
	"fmt"
	"log/slog"
	"net"
	"strings"

	// 3rd Party
	"github.com/chzyer/readline"
	"github.com/google/uuid"

	// Internal
	"github.com/Ne0nd0g/merlin-cli/commands"
	"github.com/Ne0nd0g/merlin-cli/entity/help"
	"github.com/Ne0nd0g/merlin-cli/entity/menu"
	"github.com/Ne0nd0g/merlin-cli/entity/os"
	"github.com/Ne0nd0g/merlin-cli/message"
	"github.com/Ne0nd0g/merlin-cli/services/rpc"
)

// Command is an aggregate structure for a command executed on the command line interface
type Command struct {
	name   string      // name is the name of the command
	help   help.Help   // help is the Help structure for the command
	menus  []menu.Menu // menu is the Menu the command can be used in
	native bool        // native is true if the command is executed by an Agent using only Golang native code
	os     os.OS       // os is the supported operating system the Agent command can be executed on
}

// NewCommand is a factory that builds and returns a Command structure that implements the Command interface
func NewCommand() *Command {
	var cmd Command
	cmd.name = "listener"
	cmd.menus = []menu.Menu{menu.AGENT}
	cmd.os = os.ALL
	description := "Start, stop, or list peer-to-peer listeners on the Agent"
	// Style guide for usage https://developers.google.com/style/code-syntax
	usage := "listener {list|start|stop} [protocol] [address]"
	example := ""
	notes := "Use '-h' after the subcommand to get more information"
	cmd.help = help.NewHelp(description, example, notes, usage)
	return &cmd
}

// Completer returns the data that is displayed in the CLI for tab completion depending on the menu the command is for
// Errors are not returned to ensure the CLI is not interrupted.
// Errors are logged and can be viewed by enabling debug output in the CLI
func (c *Command) Completer(menu.Menu, uuid.UUID) readline.PrefixCompleterInterface {
	comp := readline.PcItem(c.name,
		readline.PcItem("list"),
		readline.PcItem("start",
			readline.PcItem("tcp",
				readline.PcItem("127.0.0.1:7777"),
			),
			readline.PcItem("udp",
				readline.PcItem("127.0.0.1:7777"),
			),
			readline.PcItem("smb",
				readline.PcItem("merlinpipe"),
			),
		),
		readline.PcItem("stop",
			readline.PcItem("tcp",
				readline.PcItem("127.0.0.1:7777"),
			),
			readline.PcItem("udp",
				readline.PcItem("127.0.0.1:7777"),
			),
			readline.PcItem("smb",
				readline.PcItem("merlinpipe"),
			),
		),
	)
	return comp
}

// Do executes the command and returns a Response to the caller to facilitate changes in the CLI service
// m, an optional parameter, is the Menu the command was executed from
// id, an optional parameter, used to identify a specific Agent or Listener
// arguments, and optional, parameter, is the full unparsed string entered on the command line to include the
// command itself passed into command for processing
func (c *Command) Do(m menu.Menu, id uuid.UUID, arguments string) (response commands.Response) {
	slog.Debug("entering into function", "menu", m, "id", id, "arguments", arguments)
	// Parse the arguments
	args := strings.Split(arguments, " ")

	// Validate at least one argument, in addition to the command, was provided
	if len(args) < 2 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command requires at least one argument\n%s", c, c.help.Usage()))
		return
	}

	switch strings.ToLower(args[1]) {
	case "list":
		return c.List(id, arguments)
	case "start", "stop":
		return c.Start(id, arguments)
	case "help", "-h", "--help", "?", "/?":
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, c.help.Description(), c.help.Usage(), c.help.Example(), c.help.Notes()))
		return
	default:
		response.Message = message.NewUserMessage(message.Info, c.help.Usage())
		return
	}
}

func (c *Command) List(id uuid.UUID, arguments string) (response commands.Response) {
	sub := "list"

	// Create Help for this sub command
	description := "Instruct the Agent to return a list of its peer-to-peer listeners"
	example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» listener list\n" +
		"\t[-] Created job OebvmmBQPr for agent c1090dbc-f2f7-4d90-a241-86e0c0217786 at 2023-07-23T16:44:33Z\n" +
		"\t[-] Results of job OebvmmBQPr for agent c1090dbc-f2f7-4d90-a241-86e0c0217786 at 2023-07-23T16:44:53Z\n\n" +
		"\t[+] Peer-to-Peer Listeners (3):\n" +
		"\t0. TCP listener on 127.0.0.1:7777\n" +
		"\t1. UDP listener on [::]:8888\n" +
		"\t2. SMB listener on \\\\.\\pipe\\merlinpipe\n"
	notes := "The string `[::]` signifies all IP interfaces"
	usage := "listener list"
	h := help.NewHelp(description, example, notes, usage)

	// Parse the arguments
	args := strings.Split(arguments, " ")

	// 0. listener, 1. list, 2. -h
	if len(args) > 2 {
		switch strings.ToLower(args[2]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s %s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, sub, h.Description(), h.Usage(), h.Example(), h.Notes()))
			return
		}
	}
	response.Message = rpc.Listener(id, args[1:])
	return
}

func (c *Command) Start(id uuid.UUID, arguments string) (response commands.Response) {
	// Parse the arguments
	args := strings.Split(arguments, " ")

	var h help.Help
	var sub string

	if len(args) < 2 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command requires at least one argument\n%s", c, c.help.Usage()))
		return
	}

	switch strings.ToLower(args[1]) {
	case "start":
		sub = "start"
		description := "Instruct the Agent to start a peer-to-peer listener"
		example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» listener start tcp 127.0.0.1:7777 \n" +
			"\t[-] Created job LkIuWumcOt for agent c1090dbc-f2f7-4d90-a241-86e0c0217786 at 2023-07-23T16:40:24Z\n" +
			"\t[-] Results of job LkIuWumcOt for agent c1090dbc-f2f7-4d90-a241-86e0c0217786 at 2023-07-23T16:40:37Z\n" +
			"\t[+] Successfully started TCP listener on 127.0.0.1:7777\n\n" +
			"\tMerlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» listener start udp 0.0.0.0:8888\n" +
			"\t[-] Created job suVecDPJhC for agent d942a9a5-a68e-42e7-8d26-71ac45e8345a at 2023-07-23T16:41:43Z\n" +
			"\t[-] Results of job suVecDPJhC for agent d942a9a5-a68e-42e7-8d26-71ac45e8345a at 2023-07-23T16:41:56Z\n" +
			"\t[+] Successfully started UDP listener on 0.0.0.0:8888\n"
		notes := "Use '0.0.0.0' for all IPv4 interfaces. Only provide the name of the pipe for the SMB listener (e.g., merlinPipe)"
		usage := "listener start {smb|tcp|udp} {namedPipe|<interface:port>}"
		h = help.NewHelp(description, example, notes, usage)
	case "stop":
		sub = "stop"
		description := "Instruct the Agent to stop a peer-to-peer listener"
		example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» listener stop tcp 127.0.0.1:7777\n" +
			"\t[-] Created job zlVVVBDCVS for agent c1090dbc-f2f7-4d90-a241-86e0c0217786 at 2023-07-23T16:53:58Z\n" +
			"\t[-] Results of job zlVVVBDCVS for agent c1090dbc-f2f7-4d90-a241-86e0c0217786 at 2023-07-23T16:54:18Z\n" +
			"\t[+] Successfully closed TCP listener on 127.0.0.1:7777"
		notes := ""
		usage := "listener stop {smb|tcp|udp} {namedPipe|<interface:port>}"
		h = help.NewHelp(description, example, notes, usage)
	default:
		response.Message = message.NewErrorMessage(fmt.Errorf("unknown listener command '%s'\n%s", args[1], c.help.Usage()))
		return
	}

	// 0. listener, 1. start, 2. -h
	if len(args) > 2 {
		switch strings.ToLower(args[2]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s %s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, sub, h.Description(), h.Usage(), h.Example(), h.Notes()))
			return
		}
	}

	// 0. listener, 1. start, 2. protocol, 3. interface:port/named pipe
	if len(args) < 4 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s %s' command requires two arguments\n%s", c, sub, h.Usage()))
		return
	}

	switch strings.ToLower(args[2]) {
	case "smb", "tcp", "udp":
		// Pass
	default:
		response.Message = message.NewErrorMessage(fmt.Errorf("'%s' is not a valid protocol", args[2]))
		return
	}

	if strings.ToLower(args[2]) == "tcp" || strings.ToLower(args[2]) == "udp" {
		// Client side validate interface and port
		addr := strings.Split(args[3], ":")
		if len(addr) != 2 {
			response.Message = message.NewErrorMessage(fmt.Errorf("'%s' is not a valid IP address and port:\n%s", args[3], h.Usage()))
			return
		}
		if net.ParseIP(addr[0]) == nil {
			response.Message = message.NewErrorMessage(fmt.Errorf("'%s' is not a valid IP address", addr[0]))
			return
		}
	}
	response.Message = rpc.Listener(id, args[1:])
	return
}

// Help returns a help.Help structure that can be used to view a command's Description, Notes, Usage, and an example
func (c *Command) Help(menu.Menu) help.Help {
	return c.help
}

// Menu checks to see if the command is supported for the provided menu
func (c *Command) Menu(m menu.Menu) bool {
	for _, v := range c.menus {
		if v == m || v == menu.ALLMENUS {
			return true
		}
	}
	return false
}

// OS returns the supported operating system the Agent command can be executed on
func (c *Command) OS() os.OS {
	return c.os
}

// String returns the unique name of the command as a string
func (c *Command) String() string {
	return c.name
}
