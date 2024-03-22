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

package reload

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
	moduleMemory "github.com/Ne0nd0g/merlin-cli/module/memory"
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
	// name MUST be unique across all commands and is used as the primary key in the database
	cmd.name = "reload"
	cmd.menus = []menu.Menu{menu.MODULE}
	cmd.os = os.LOCAL
	description := "reload the module to reset the module's state"
	// Style guide for usage https://developers.google.com/style/code-syntax
	usage := "reload"
	example := "Merlin[module][Invoke-Mimikatz]» reload\n\tMerlin[module][Invoke-Mimikatz]»"
	notes := ""
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
	if len(args) > 1 {
		switch strings.ToLower(args[1]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, c.help.Description(), c.help.Usage(), c.help.Example(), c.help.Notes()))
			return
		}
	}

	// Get options from the local repository
	repo := moduleMemory.NewRepository()
	module, err := repo.Get(id)
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("pkg/cli/commands/reload.DoModule(): there was an error getting module ID %s from the repository: %s", id, err))
		return
	}

	err = repo.UpdateOption(id, "agent", "")
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("pkg/cli/commands/reload.DoModule(): there was an error resetting the agent ID: %s", err))
		return
	}

	err = repo.Reload(id)
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("pkg/cli/commands/reload.DoModule(): there was an error reloading the module: %s", err))
		return
	}

	// Store the updated module
	/*
		err = repo.Update(id, module)
		if err != nil {
			response.Message = message.NewErrorMessage(fmt.Errorf("pkg/cli/commands/reload.DoModule(): there was an error storing the updated module: %s", err))
			return
		}

	*/
	response.Message = message.NewUserMessage(message.Success, fmt.Sprintf("The '%s' module was reloaded", module.String()))
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
