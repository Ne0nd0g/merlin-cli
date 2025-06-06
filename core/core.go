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

package core

import (
	"os"
	"sync"
)

// STDOUT is a global mutex to prevent concurrent writes to STDOUT
var STDOUT sync.Mutex

// Debug puts Merlin into debug mode and displays debug messages
var Debug = false

// Verbose puts Merlin into verbose mode and displays verbose messages
var Verbose = false

// CurrentDir is the current directory where Merlin was executed from
var CurrentDir, _ = os.Getwd()
