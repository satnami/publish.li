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
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("path=%s\n", path)

	if path == "/" {
		log.Printf("here\n")
		// serve the homepage
		http.ServeFile(w, r, "./static/index.html")
	} else if path == "/favicon.ico" {
		http.ServeFile(w, r, "static/favicon.ico")
	} else {
		// everything else
		fmt.Fprintf(w, "Page=%q\n", html.EscapeString(r.URL.Path))
	}
}

func main() {
	// set up the static file server
	static := http.FileServer(http.Dir("static"))

	// use the default mux
	http.Handle("/s/", static)
	http.HandleFunc("/", handler)

	// the server
	port := os.Getenv("PORT")
	log.Printf("Listening on port %s ...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
}

// --------------------------------------------------------------------------------------------------------------------
