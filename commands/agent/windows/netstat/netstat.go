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

package netstat

import (
	// Standard
	"fmt"
	"log/slog"
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
	cmd.name = "netstat"
	cmd.menus = []menu.Menu{menu.AGENT}
	cmd.os = os.WINDOWS
	description := "Get a list of network connections"
	// Style guide for usage https://developers.google.com/style/code-syntax
	usage := "netstat [-p tcp|udp]"
	example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» netstat\n" +
		"\t[-] Created job JEFMANkdaU for agent c1090dbc-f2f7-4d90-a241-86e0c0217786\n" +
		"\t[-] Results job JEFMANkdaU for agent c1090dbc-f2f7-4d90-a241-86e0c0217786\n" +
		"\t[+]\n" +
		"\tProto Local Addr              Foreign Addr            State        PID/Program name\n" +
		"\tudp   0.0.0.0:123             0.0.0.0:0                            3272/svchost.exe\n" +
		"\tudp   0.0.0.0:500             0.0.0.0:0                            3104/svchost.exe\n" +
		"\tudp   0.0.0.0:3389            0.0.0.0:0                            984/svchost.exe\n" +
		"\tudp6  :::123                  0.0.0.0:0                            3272/svchost.exe\n" +
		"\tudp6  :::500                  0.0.0.0:0                            3104/svchost.exe\n" +
		"\tudp6  :::3389                 0.0.0.0:0                            984/svchost.exe\n" +
		"\ttcp   0.0.0.0:135             0.0.0.0:0               LISTEN       964/svchost.exe\n" +
		"\ttcp   0.0.0.0:445             0.0.0.0:0               LISTEN       4/System\n" +
		"\ttcp   0.0.0.0:3389            0.0.0.0:0               LISTEN       984/svchost.exe\n" +
		"\ttcp   127.0.0.1:52945         127.0.0.1:5357          TIME_WAIT\n" +
		"\ttcp   127.0.0.1:54441         127.0.0.1:5357          TIME_WAIT\n" +
		"\ttcp   192.168.1.11:59757      72.21.91.29:80          CLOSE_WAIT   6496/SearchApp.exe\n" +
		"\ttcp   192.168.1.11:59763      72.21.91.29:80          CLOSE_WAIT   12076/YourPhone.exe\n" +
		"\ttcp6  :::135                  :::0                    LISTEN       964/svchost.exe\n" +
		"\ttcp6  :::445                  :::0                    LISTEN       4/System\n" +
		"\ttcp6  :::3389                 :::0                    LISTEN       984/svchost.exe"
	notes := "This command is only available on Windows. It uses the Windows API to enumerate network " +
		"connections and listening ports. Without any arguments, the netstat command returns all TCP and UDP network " +
		"connections.\n" +
		"\tUse 'netstat -p tcp' to only return TCP connections and 'netstat -p udp' to only return UDP connections."
	cmd.help = help.NewHelp(description, example, notes, usage)
	return &cmd
}

// Completer returns the data that is displayed in the CLI for tab completion depending on the menu the command is for
// Errors are not returned to ensure the CLI is not interrupted.
// Errors are logged and can be viewed by enabling debug output in the CLI
func (c *Command) Completer(menu.Menu, uuid.UUID) readline.PrefixCompleterInterface {
	return readline.PcItem(c.name)
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

	// Check for help first
	// 0. netstat, 1. -p|-h, 2. tcp|udp
	if len(args) > 1 {
		switch strings.ToLower(args[1]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, c.help.Description(), c.help.Usage(), c.help.Example(), c.help.Notes()))
			return
		case "-p":
			if len(args) < 2 {
				response.Message = message.NewUserMessage(message.Warn, "Invalid argument for -p. Valid arguments are 'tcp' or 'udp'.")
				return
			}
		default:
			response.Message = message.NewUserMessage(message.Info, c.help.Usage())
			return
		}
	}
	response.Message = rpc.Netstat(id, args[1:])
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
