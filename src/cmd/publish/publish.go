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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// open the db
	db, errOpen := bolt.Open("publish.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	check(errOpen)
	defer db.Close()

	errUpdate := db.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists(pageBucketName)
		if err1 != nil {
			return err1
		}

		_, err2 := tx.CreateBucketIfNotExists(idBucketName)
		if err2 != nil {
			return err2
		}

		return nil
	})
	check(errUpdate)

	// set up the static file server
	static := http.FileServer(http.Dir("static"))

	// use the default mux
	http.HandleFunc("/api", apiHandler(db))
	http.Handle("/s/", static)
	http.HandleFunc("/", homeHandler(db))

	// the server
	port := os.Getenv("PORT")
	log.Printf("Listening on port %s ...\n", port)
	errListen := http.ListenAndServe(":"+port, nil)
	log.Fatal(errListen)
}
