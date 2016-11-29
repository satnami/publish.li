// --------------------------------------------------------------------------------------------------------------------
//
// This file is part of https://github.com/appsattic/publish.li
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
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Machiel/slugify"
	"github.com/boltdb/bolt"
)

var tmpl *template.Template

func init() {
	tmpl1, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl = tmpl1
}

func render(w http.ResponseWriter, templateName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func apiPut(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page{}

		// parse the incoming JSON request
		decoder := json.NewDecoder(r.Body)
		errDecode := decoder.Decode(&page)
		if errDecode != nil {
			log.Printf("Error: %v\n", errDecode)
			sendError(w, "Invalid JSON")
			return
		}
		defer r.Body.Close()

		// check that the title has something in it (other than whitespace)
		rawTitle := page.Title

		slug := slugify.Slugify(rawTitle)
		if slug == "" {
			sendError(w, "Provide a title")
			return
		}

		// fill in the other fields to save this page
		now := time.Now()
		page.Id = randStr(16)
		page.Name = slug + "-" + randStr(8)
		page.Inserted = now
		page.Updated = now

		errIns := storePutPage(db, page)
		if errIns != nil {
			http.Error(w, errIns.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Ok   bool   `json:"ok"`
			Msg  string `json:"msg"`
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			Ok:   true,
			Msg:  "Saved",
			Id:   page.Id,
			Name: page.Name,
		}

		sendJson(w, data)
	}
}

func apiPost(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page{}

		// parse the incoming JSON request
		decoder := json.NewDecoder(r.Body)
		errDecode := decoder.Decode(&page)
		if errDecode != nil {
			log.Printf("Error: %v\n", errDecode)
			sendError(w, "Invalid JSON")
			return
		}
		defer r.Body.Close()

		// using the page.Name, retrieve this page then check it's Id is correct
		existPage, errGet := storeGetPage(db, page.Name)
		if errGet != nil {
			log.Printf("Error: %v\n", errGet)
			sendError(w, "Internal Error. Please try again later.")
			return
		}

		if existPage == nil {
			sendError(w, "This page name does not exist.")
			return
		}

		// check that this page has this Id
		if existPage.Id != page.Id {
			sendError(w, "Permission denied.")
			return
		}

		// check that the title has something in it (other than whitespace)
		rawTitle := page.Title
		replacer := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\f", "")
		page.Title = replacer.Replace(rawTitle)

		if page.Title == "" {
			sendError(w, "Provide a title")
			return
		}

		// We don't trust what is in `page`, but we know `existPage` is fine, so we'll just update a couple of fields
		// there, then re-save.
		now := time.Now()
		existPage.Title = page.Title
		existPage.Author = page.Author
		existPage.Content = page.Content
		existPage.Updated = now

		errIns := storePutPage(db, *existPage)
		if errIns != nil {
			http.Error(w, errIns.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Ok   bool   `json:"ok"`
			Msg  string `json:"msg"`
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			Ok:   true,
			Msg:  "Saved",
			Id:   page.Id,
			Name: page.Name,
		}

		sendJson(w, data)
	}
}

func apiGet(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page{}

		// Must pass a skeleton page consisting of an `Id` and a `Name`

		// parse the incoming JSON request
		decoder := json.NewDecoder(r.Body)
		errDecode := decoder.Decode(&page)
		if errDecode != nil {
			log.Printf("Error: %v\n", errDecode)
			sendError(w, "Invalid JSON")
			return
		}
		defer r.Body.Close()

		// retrieve this page
		existPage, errGet := storeGetPage(db, page.Name)
		if errGet != nil {
			log.Printf("Error: %v\n", errGet)
			sendError(w, "Internal Error. Please try again later.")
			return
		}

		if existPage == nil {
			sendError(w, "This page name does not exist.")
			return
		}

		// check that this page has this Id
		if existPage.Id != page.Id {
			sendError(w, "Permission denied.")
			return
		}

		data := struct {
			Ok   bool   `json:"ok"`
			Msg  string `json:"msg"`
			Data *Page
		}{
			Ok:   true,
			Msg:  "Saved",
			Data: existPage,
		}

		sendJson(w, data)
	}
}

func apiHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	insPost := apiPut(db)
	savePost := apiPost(db)
	getPost := apiGet(db)

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path != "/api" {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

		if r.Method == "PUT" {
			insPost(w, r)
			return
		}

		if r.Method == "POST" {
			savePost(w, r)
			return
		}

		if r.Method == "GET" {
			getPost(w, r)
			return
		}
	}
}

func servePage(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// everything else
	name := r.URL.Path[1:]
	log.Printf("Page=%q\n", html.EscapeString(name))

	page, errPage := storeGetPage(db, name)
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
	render(w, "page.html", page)
}

func homeHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("path=%s\n", path)

		// simple routing instead of using a complicated router (until we need one)
		if path == "/" {
			http.ServeFile(w, r, "./static/index.html")

		} else if path == "/favicon.ico" {
			http.ServeFile(w, r, "static/favicon.ico")

		} else if path == "/robots.txt" {
			http.ServeFile(w, r, "static/robots.txt")

		} else {
			servePage(w, r, db)
		}
	}
}
