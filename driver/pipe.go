// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package driver

import (
	"errors"
	"net"
)

const (
	inBasePort  = 8000
	outBasePort = 8000
)

// Node alloc Pipe
type Node struct {
	// IP is the Svr IP
	IP net.IP

	// Prefix to identity services
	Prefix int

	all map[int]*pipe
}

// pipe is pipeline shared between workers
type pipe struct {
	inPorts []int

	outIP []net.IP

	outPorts [][]int

	inWorkers  Worker
	outWorkers []Worker
}

// Session is src or dst for workers
type Session struct {
	IP net.IP

	Ports []int
}

var (
	errNodeBadInput = errors.New("Bad input for node")
)

func helperPort(base int, prefix int, id int) []int {
	return []int{base + prefix + 4*id, base + prefix + 4*id + 2}
}

// Create a svr
func (n *Node) Create() {
	n.all = make(map[int]*pipe)
}

// AllocPull alloc one pull
func (n *Node) AllocPull(id int, w Worker) error {
	var p *pipe

	if p = n.all[id]; p == nil {
		p = &pipe{inPorts: helperPort(inBasePort, n.Prefix, id)}
		n.all[id] = p
	}

	if w == nil || !IsWorkerDec(w) {
		return errNodeBadInput
	}

	for _, exists := range p.outWorkers {
		if exists == w {
			return nil
		}
	}

	wid := GetWorkerWorkerID(w)
	// IP := GetWorkerWorkerIP(w)

	ses := Session{Ports: helperPort(outBasePort, n.Prefix, wid)}
	if err := SetDecodeSes(w, &ses); err != nil {
		return err
	}

	// TODO: start here

	p.outWorkers = append(p.outWorkers, w)

	return nil
}

// FreePull free one pull
func (n *Node) FreePull(id int, w Worker) error {
	var p *pipe
	if p = n.all[id]; p == nil {
		return nil
	}

	if w == nil || !IsWorkerDec(w) {
		return errNodeBadInput
	}

	var k int
	var exists Worker
	for k, exists = range p.outWorkers {
		if exists == w {
			break
		}
	}

	if exists != w {
		return nil
	}

	// TODO: free here

	p.outWorkers = remove(p.outWorkers, k)

	return nil
}

// AllocPush alloc one push
func (n *Node) AllocPush(id int, w Worker) error {
	var p *pipe

	if p = n.all[id]; p == nil {
		p = &pipe{inPorts: helperPort(inBasePort, n.Prefix, id)}
		n.all[id] = p
	}

	if w == nil || !IsWorkerEnc(w) {
		return errNodeBadInput
	}

	if exists := p.inWorkers; exists != nil {
		if exists == w {
			return nil
		}
		// TODO re-do
	}

	ses := Session{IP: n.IP,
		Ports: p.inPorts}

	if err := SetEncodeSes(w, &ses); err != nil {
		return err
	}

	// TODO: push here

	p.inWorkers = w

	return nil
}

// FreePush free one push
func (n *Node) FreePush(id int) error {
	var p *pipe
	if p = n.all[id]; p == nil {
		return nil
	}

	// TODO: free here

	p.inWorkers = nil

	return nil
}

// https://yourbasic.org/golang/delete-element-slice/
func remove(ws []Worker, i int) []Worker {
	// Remove the element at index i from a.
	ws[i] = ws[len(ws)-1] // Copy last element to index i.
	ws[len(ws)-1] = nil   // Erase last element (write zero value).
	ws = ws[:len(ws)-1]   // Truncate slice.

	return ws
}
