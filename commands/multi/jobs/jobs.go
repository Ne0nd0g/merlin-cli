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

package jobs

import (
	// Standard
	"fmt"
	"log/slog"
	"strings"

	// 3rd Party
	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"

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
	cmd.name = "jobs"
	cmd.menus = []menu.Menu{menu.AGENT, menu.MAIN}
	cmd.os = os.LOCAL
	description := "Display all unfinished jobs"
	usage := "jobs"
	example := "Merlin» jobs\n\n" +
		"\t\t\t AGENT                 |     ID     |  COMMAND   | STATUS  |       CREATED        |         SENT\n" +
		"\t+--------------------------------------+------------+------------+---------+----------------------+----------------------+\n" +
		"\t  d07edfda-e119-4be2-a20f-918ab701fa3c | UjNoTALgcn | pwd        | Created | 2021-08-03T01:39:57Z |\n" +
		"\t  99dbe632-984c-4c98-8f38-11535cb5d937 | UHOddpFQTm | run whoami | Sent    | 2021-08-03T01:40:11Z | 2021-08-03T01:40:17Z"
	notes := "Only the first 30 characters of the COMMAND are displayed"
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

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)

	var data [][]string
	switch m {
	case menu.AGENT:
		jobs, err := rpc.GetAgentActiveJobs(id)
		if err != nil {
			response.Message = message.NewErrorMessage(err)
			return
		}
		table.SetHeader([]string{"ID", "Command", "Status", "Created", "Sent"})
		for _, job := range jobs {
			var row []string
			if len(job.Command) < 30 {
				row = []string{job.ID, job.Command, job.Status, job.Created, job.Sent}
			} else {
				row = []string{job.ID, job.Command[:30], job.Status, job.Created, job.Sent}
			}
			data = append(data, row)
		}
	default:
		jobs, err := rpc.GetAllActiveJobs()
		if err != nil {
			response.Message = message.NewErrorMessage(err)
			return
		}
		table.SetHeader([]string{"Agent", "ID", "Command", "Status", "Created", "Sent"})
		for _, job := range jobs {
			var row []string
			if len(job.Command) < 30 {
				row = []string{job.AgentID, job.ID, job.Command, job.Status, job.Created, job.Sent}
			} else {
				row = []string{job.AgentID, job.ID, job.Command[:30], job.Status, job.Created, job.Sent}
			}
			data = append(data, row)
		}
	}

	table.AppendBulk(data)
	table.Render()

	response.Message = message.NewUserMessage(message.Plain, fmt.Sprintf("\n%s", tableString.String()))

	return
}

// Help returns a help.Help structure that can be used to view a command's Description, Notes, Usage, and an example
func (c *Command) Help(menu.Menu) help.Help {
	return c.help
}

// Menu checks to see if the command is supported for the provided menu
func (c *Command) Menu(m menu.Menu) bool {
	for _, v := range c.menus {
		if v == m {
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
