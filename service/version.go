// Copyright © 2018 Job King'ori Maina <j@kingori.co>

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

var (
	majorVersion = "0"
	minorVersion = "0"
	patchVersion = "0"
)

// Version represents version information.
type Version struct {
	Major string
	Minor string
	Patch string
}

// GetVersion returns version information.
func GetVersion() Version {
	return Version{
		Major: majorVersion,
		Minor: minorVersion,
		Patch: patchVersion,
	}
}

// Str returns version information as a string
func (v *Version) Str() string {
	return v.Major + "." + v.Minor + "." + v.Patch
}
