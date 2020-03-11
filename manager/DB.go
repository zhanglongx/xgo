// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package manager

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zhanglongx/Aqua/comm"
)

// DBVER is DB File Version
const DBVER string = "1.0.0"

// DB contains all path' config. It's degsinged to be easily
// exported to file (like JSON).
// set() and get() are not thread-safe, it's caller's
// responsibility to ensure that. To ensure data returned by
// get() will not get rewritten, data from pathRow should be
// copied before any unlock method.
type DB struct {
	// Version should be used to check DB's compatibility
	Version string

	// Store contains all path params
	Store map[string]*Params
}

// loadFromFile load JSON file to Cfg
func (d *DB) loadFromFile(JFile string) error {

	buf, err := ioutil.ReadFile(JFile)
	if err != nil {
		comm.Error.Printf("Read DB file %s failed", JFile)
		return err
	}

	err = json.Unmarshal(buf, d)
	if err != nil {
		comm.Error.Printf("Decode DB file %s failed", JFile)
		return err
	}

	// FIXME: more compatible
	if d.Version != DBVER {
		comm.Error.Printf("DB file ver error: %s", d.Version)
		comm.Error.Printf("Discarding old file: %s", JFile)
		d.Store = make(map[string]*Params, 0)
		return nil
	}

	return nil
}

// saveToFile save JSON file to Cfg
func (d *DB) saveToFile(JFile string) error {

	buf, err := json.Marshal(d)
	if err != nil {
		comm.Error.Printf("Encode DB %s failed", JFile)
		return err
	}

	err = ioutil.WriteFile(JFile, buf, 0644)
	if err != nil {
		comm.Error.Printf("Write DB file %s failed", JFile)
		return err
	}

	return nil
}

// set set a new pathRow in DB, DB store *ONLY* the pointer
// passed in, so make sure passing a whole new pathRow{}
// everytime
func (d *DB) set(ID string, p *Params) error {

	d.Store[ID] = p

	return nil
}

// get return a *pathRow in DB. Because set() and get() are
// not thread-safe, you should get data in *pathRow copied,
// before using any unlock method
func (d *DB) get(ID string) *Params {

	if d.Store[ID] == nil {
		return nil
	}

	return d.Store[ID]
}
