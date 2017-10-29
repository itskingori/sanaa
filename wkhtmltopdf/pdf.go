// Copyright © 2017 Job King'ori Maina <j@kingori.co>
//
// This file is part of sanaa.
//
// sanaa is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// sanaa is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with sanaa. If not, see <http://www.gnu.org/licenses/>.

package wkhtmltopdf

// Pdf is a mapping of wkhtmltoimage parameters
type Pdf struct {
	MarginTop    int `json:"margin_top"`
	MarginBottom int `json:"margin_bottom"`
	MarginLeft   int `json:"margin_left"`
	MarginRight  int `json:"margin_right"`
	PageHeight   int `json:"page_height"`
	PageWidth    int `json:"page_width"`
}
