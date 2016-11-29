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
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

var ErrFatalNoPagesBucket = errors.New("Bucket 'pages' does not exist")

var pagesBucket = []byte("pages")

func storeGetPage(db *bolt.DB, name string) (*Page, error) {
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

func storePutPage(db *bolt.DB, page Page) error {
	return db.Update(func(tx *bolt.Tx) error {
		pages := tx.Bucket(pagesBucket)
		if pages == nil {
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

		return nil
	})
}
