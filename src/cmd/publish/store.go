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

var ErrFatalNoPageBucket = errors.New("Bucket 'page' does not exist")
var ErrFatalNoIdBucket = errors.New("Bucket 'id' does not exist")

var pageBucketName = []byte("page")
var idBucketName = []byte("id")

func storeGetPageUsingId(db *bolt.DB, id string) (*Page, error) {
	var p *Page

	// firstly, get this Id and see if it exists
	err := db.View(func(tx *bolt.Tx) error {
		idBucket := tx.Bucket(idBucketName)
		if idBucket == nil {
			panic(ErrFatalNoIdBucket)
		}

		// get this Id
		rawName := idBucket.Get([]byte(id))
		if rawName == nil {
			return nil
		}

		// now get this page
		name := string(rawName)

		pageBucket := tx.Bucket(pageBucketName)
		if pageBucket == nil {
			panic(ErrFatalNoPageBucket)
		}

		rawPage := pageBucket.Get([]byte(name))

		// see if this exists
		if rawPage == nil {
			log.Printf("Not Found : %s\n", name)
			return nil
		}

		// try and decode
		page := Page{}
		err := json.Unmarshal(rawPage, &page)
		if err != nil {
			return err
		}
		p = &page
		return nil
	})

	return p, err
}

func storeGetPage(db *bolt.DB, name string) (*Page, error) {
	var p *Page

	// see if we can find this page
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(pageBucketName)
		if b == nil {
			panic(ErrFatalNoPageBucket)
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
		pageBucket := tx.Bucket(pageBucketName)
		if pageBucket == nil {
			panic(ErrFatalNoPageBucket)
		}

		// write this page out
		bytes, errMarshal := json.Marshal(page)
		if errMarshal != nil {
			return errMarshal
		}

		errPutPage := pageBucket.Put([]byte(page.Name), bytes)
		if errPutPage != nil {
			return errPutPage
		}

		// and make sure we have an Id pointing to this name
		idBucket := tx.Bucket(idBucketName)
		if idBucket == nil {
			panic(ErrFatalNoIdBucket)
		}

		errPutId := idBucket.Put([]byte(page.Id), []byte(page.Name))
		if errPutId != nil {
			return errPutId
		}

		return nil
	})
}
