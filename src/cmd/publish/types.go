// --------------------------------------------------------------------------------------------------------------------
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// --------------------------------------------------------------------------------------------------------------------

package main

import (
	"html/template"
	"time"
)

type Page struct {
	Id        string        `json:"id"`        // e.g. "aoAAhc5i4bmKMSZk"
	Name      string        `json:"name"`      // e.g. "first-post-chzc9BkU
	Title     string        `json:"title"`     // e.g. "First Post"
	Author    string        `json:"author"`    // e.g. "Andrew Chilton"
	Website   string        `json:"website"`   // e.g. "https://chilts.org"
	Twitter   string        `json:"twitter"`   // e.g. "andychilton"
	Facebook  string        `json:"facebook"`  // e.g. "andrew.peter.chilton"
	GitHub    string        `json:"github"`    // e.g. "chilts"
	Instagram string        `json:"instagram"` // e.g. "thechilts"
	Content   string        `json:"content"`   // e.g. "My story."
	Html      template.HTML `json:"html"`      // i.e. the transformed Markdown into HTML
	Inserted  time.Time     `json:"inserted"`  // i.e. The inserted time
	Updated   time.Time     `json:"updated"`   // i.e. The updated time
}
