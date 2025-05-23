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

package agent

// Process is a structure that holds information about the Process the Agent is running in/as
type Process struct {
	ID        int32  // The process ID that the agent is running in
	Integrity int32  // The integrity level of the process the agent is running in
	Name      string // The process name that the agent is running in
	UserGUID  string // The GUID of the user that the agent is running as
	UserName  string // The username that the agent is running as
	Domain    string // The domain the user running the process belongs to
}
