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
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Machiel/slugify"
	"github.com/boltdb/bolt"
)

var ErrFatalNoPagesBucket = errors.New("Bucket 'pages' does not exist")
var ErrFatalNoIdsBucket = errors.New("Bucket 'ids' does not exist")

// var letters = []bytes("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lenLetters = len(letterBytes)

var pagesBucket = []byte("pages")
var idsBucket = []byte("ids")

type Page struct {
	Id       string    `json:"id"`      // e.g. "aoAAhc5i4bmKMSZk"
	Name     string    `json:"name"`    // e.g. "first-post-chzc9BkU
	Title    string    `json:"title"`   // e.g. "First Post"
	Author   string    `json:"author"`  // e.g. "Andrew Chilton"
	Content  string    `json:"content"` // e.g. "My story."
	Inserted time.Time `json:"inserted"`
	Updated  time.Time `json:"updated"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func sendJson(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func sendOk(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, msg string) {
	data := struct {
		Ok  bool   `json:"ok"`
		Msg string `json:"msg"`
	}{
		Ok:  false,
		Msg: msg,
	}

	sendJson(w, data)
}

func newPage(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// parse the JSON request
	decoder := json.NewDecoder(r.Body)
	page := Page{}
	errDecode := decoder.Decode(&page)
	if errDecode != nil {
		// http.Error(w, errDecode.Error(), http.StatusInternalServerError)
		sendError(w, "Invalid JSON")
		return
	}
	defer r.Body.Close()

	// we only need to validate if we have something in the title
	rawTitle := page.Title
	replacer := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\f", "")
	page.Title = replacer.Replace(rawTitle)

	if page.Title == "" {
		sendError(w, "Provide a title")
		return
	}

	// set the other fields
	page.Id = randStr(16)
	page.Name = slugify.Slugify(rawTitle) + "-" + randStr(8)
	now := time.Now()
	page.Inserted = now
	page.Updated = now

	errIns := db.Update(func(tx *bolt.Tx) error {
		pages := tx.Bucket(pagesBucket)
		if pages == nil {
			panic(ErrFatalNoPagesBucket)
		}

		ids := tx.Bucket(idsBucket)
		if ids == nil {
			panic(ErrFatalNoPagesBucket)
		}

		// write this page out
		bytes, errMarshal := json.Marshal(page)
		if errMarshal != nil {
			return errMarshal
		}

		errPutPage := pages.Put([]byte(page.Name), bytes)
		if errPutPage != nil {
			return errPutPage
		}

		errPutId := ids.Put([]byte(page.Id), []byte(page.Name))
		if errPutId != nil {
			return errPutId
		}

		return nil
	})

	if errIns != nil {
		http.Error(w, errIns.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Ok  bool   `json:"ok"`
		Msg string `json:"msg"`
		Id  string `json:"id"`
	}{
		Ok:  true,
		Msg: "Saved",
		Id:  page.Id,
	}

	sendJson(w, data)
}

func getPage(db *bolt.DB, name string) (*Page, error) {
	var p *Page

	// see if we can find this page
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(pagesBucket)
		if b == nil {
			panic(ErrFatalNoPagesBucket)
		}

		raw := b.Get([]byte(name))

		// see if this exists
		if raw == nil {
			log.Printf("Not Found : %s\n", name)
			return nil
		}

		// try and decode
		page := Page{}
		err := json.Unmarshal(raw, &page)
		if err != nil {
			return err
		}
		p = &page
		return nil
	})

	return p, err
}

func servePage(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// everything else
	name := r.URL.Path[1:]
	log.Printf("Page=%q\n", html.EscapeString(name))

	page, errPage := getPage(db, name)
	if errPage != nil {
		http.Error(w, errPage.Error(), http.StatusInternalServerError)
		return
	}

	if page == nil {
		log.Printf("Not Found : %s\n", name)
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	// serve the page
	fmt.Fprintf(w, "Page : %#v\n", page)
}

func savePage(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// save the page
	name := r.URL.Path[1:]
	fmt.Fprintf(w, "Page : %#v\n", name)
}

func handler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("path=%s\n", path)

		// simple routing instead of using a complicated router (until we need one)
		if path == "/" {
			if r.Method == "POST" {
				newPage(w, r, db)
				return
			}
			http.ServeFile(w, r, "./static/index.html")

		} else if path == "/favicon.ico" {
			http.ServeFile(w, r, "static/favicon.ico")

		} else if path == "/robots.txt" {
			http.ServeFile(w, r, "static/robots.txt")

		} else {
			if r.Method == "GET" {
				servePage(w, r, db)
				return
			}
			if r.Method == "POST" {
				savePage(w, r, db)
				return
			}
		}
	}
}

func main() {
	// open the db
	db, errOpen := bolt.Open("publish.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	check(errOpen)
	defer db.Close()

	errUpdate := db.Update(func(tx *bolt.Tx) error {
		var err error

		_, err = tx.CreateBucketIfNotExists(pagesBucket)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(idsBucket)
		if err != nil {
			return err
		}

		return nil
	})
	check(errUpdate)

	// set up the static file server
	static := http.FileServer(http.Dir("static"))

	// use the default mux
	http.Handle("/s/", static)
	http.HandleFunc("/", handler(db))

	// the server
	port := os.Getenv("PORT")
	log.Printf("Listening on port %s ...\n", port)
	errListen := http.ListenAndServe(":"+port, nil)
	log.Fatal(errListen)
}

// --------------------------------------------------------------------------------------------------------------------
