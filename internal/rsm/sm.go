// Copyright 2017-2019 Lei Ni (nilei81@gmail.com) and other Dragonboat authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rsm

import (
	"io"

	"github.com/lni/dragonboat/v3/config"
	pb "github.com/lni/dragonboat/v3/raftpb"
	sm "github.com/lni/dragonboat/v3/statemachine"
)

// IStateMachine is an adapter interface for underlying IStateMachine or
// IConcurrentStateMachine instances.
type IStateMachine interface {
	Open(<-chan struct{}) (uint64, error)
	Update(entries []sm.Entry) ([]sm.Entry, error)
	Lookup(query interface{}) (interface{}, error)
	NALookup(query []byte) ([]byte, error)
	Sync() error
	Prepare() (interface{}, error)
	Save(interface{},
		io.Writer, sm.ISnapshotFileCollection, <-chan struct{}) error
	Recover(io.Reader, []sm.SnapshotFile, <-chan struct{}) error
	Close() error
	GetHash() (uint64, error)
	Concurrent() bool
	OnDisk() bool
	Type() pb.StateMachineType
}

var _ IStateMachine = &RegularStateMachine{}

// RegularStateMachine is a regular state machine not capable of taking
// concurrent snapshots.
type RegularStateMachine struct {
	sm sm.IStateMachine
	h  sm.IHash
	na sm.IExtended
}

// NewRegularStateMachine creates a new RegularStateMachine instance.
func NewRegularStateMachine(s sm.IStateMachine) *RegularStateMachine {
	r := &RegularStateMachine{sm: s}
	h, ok := s.(sm.IHash)
	if ok {
		r.h = h
	}
	na, ok := s.(sm.IExtended)
	if ok {
		r.na = na
	}
	return r
}

// Open opens the state machine.
func (s *RegularStateMachine) Open(stopc <-chan struct{}) (uint64, error) {
	panic("Open() called on RegularStateMachine")
}

// Update updates the state machine.
func (s *RegularStateMachine) Update(entries []sm.Entry) ([]sm.Entry, error) {
	if len(entries) != 1 {
		panic("len(entries) != 1")
	}
	var err error
	entries[0].Result, err = s.sm.Update(entries[0].Cmd)
	return entries, err
}

// Lookup queries the state machine.
func (s *RegularStateMachine) Lookup(query interface{}) (interface{}, error) {
	return s.sm.Lookup(query)
}

// NALookup queries the state machine.
func (s *RegularStateMachine) NALookup(query []byte) ([]byte, error) {
	if s.na == nil {
		return nil, sm.ErrNotImplemented
	}
	return s.na.NALookup(query)
}

// Sync synchronizes all in-core state with that on disk.
func (s *RegularStateMachine) Sync() error {
	panic("Sync called on RegularStateMachine")
}

// Prepare makes preparations for taking concurrent snapshot.
func (s *RegularStateMachine) Prepare() (interface{}, error) {
	panic("PrepareSnapshot called on RegularStateMachine")
}

// Save saves the snapshot.
func (s *RegularStateMachine) Save(ctx interface{},
	w io.Writer, fc sm.ISnapshotFileCollection, stopc <-chan struct{}) error {
	if ctx != nil {
		panic("ctx is not nil")
	}
	return s.sm.SaveSnapshot(w, fc, stopc)
}

// Recover recovers the state machine from a snapshot.
func (s *RegularStateMachine) Recover(r io.Reader,
	fs []sm.SnapshotFile, stopc <-chan struct{}) error {
	return s.sm.RecoverFromSnapshot(r, fs, stopc)
}

// Close closes the state machine.
func (s *RegularStateMachine) Close() error {
	return s.sm.Close()
}

// GetHash returns the uint64 hash value representing the state of a state
// machine.
func (s *RegularStateMachine) GetHash() (uint64, error) {
	if s.h == nil {
		return 0, sm.ErrNotImplemented
	}
	return s.h.GetHash()
}

// Concurrent returns a boolean flag indicating whether the state machine is
// capable of taking concurrent snapshot.
func (s *RegularStateMachine) Concurrent() bool {
	return false
}

// OnDisk returns a boolean flag indicating whether this is an on disk state
// machine.
func (s *RegularStateMachine) OnDisk() bool {
	return false
}

// Type returns the type of the state machine.
func (s *RegularStateMachine) Type() pb.StateMachineType {
	return pb.RegularStateMachine
}

// ConcurrentStateMachine is an IStateMachine type capable of taking concurrent
// snapshots.
type ConcurrentStateMachine struct {
	sm sm.IConcurrentStateMachine
	h  sm.IHash
	na sm.IExtended
}

// NewConcurrentStateMachine creates a new ConcurrentStateMachine instance.
func NewConcurrentStateMachine(s sm.IConcurrentStateMachine) *ConcurrentStateMachine {
	r := &ConcurrentStateMachine{sm: s}
	h, ok := s.(sm.IHash)
	if ok {
		r.h = h
	}
	na, ok := s.(sm.IExtended)
	if ok {
		r.na = na
	}
	return r
}

// Open opens the state machine.
func (s *ConcurrentStateMachine) Open(stopc <-chan struct{}) (uint64, error) {
	panic("Open() called on RegularStateMachine")
}

// Update updates the state machine.
func (s *ConcurrentStateMachine) Update(entries []sm.Entry) ([]sm.Entry, error) {
	return s.sm.Update(entries)
}

// Lookup queries the state machine.
func (s *ConcurrentStateMachine) Lookup(query interface{}) (interface{}, error) {
	return s.sm.Lookup(query)
}

// NALookup queries the state machine.
func (s *ConcurrentStateMachine) NALookup(query []byte) ([]byte, error) {
	if s.na == nil {
		return nil, sm.ErrNotImplemented
	}
	return s.na.NALookup(query)
}

// Sync synchronizes all in-core state with that on disk.
func (s *ConcurrentStateMachine) Sync() error {
	panic("Sync called on ConcurrentStateMachine")
}

// Prepare makes preparations for taking concurrent snapshot.
func (s *ConcurrentStateMachine) Prepare() (interface{}, error) {
	return s.sm.PrepareSnapshot()
}

// Save saves the snapshot.
func (s *ConcurrentStateMachine) Save(ctx interface{},
	w io.Writer, fc sm.ISnapshotFileCollection, stopc <-chan struct{}) error {
	return s.sm.SaveSnapshot(ctx, w, fc, stopc)
}

// Recover recovers the state machine from a snapshot.
func (s *ConcurrentStateMachine) Recover(r io.Reader,
	fs []sm.SnapshotFile, stopc <-chan struct{}) error {
	return s.sm.RecoverFromSnapshot(r, fs, stopc)
}

// Close closes the state machine.
func (s *ConcurrentStateMachine) Close() error {
	return s.sm.Close()
}

// GetHash returns the uint64 hash value representing the state of a state
// machine.
func (s *ConcurrentStateMachine) GetHash() (uint64, error) {
	if s.h == nil {
		return 0, sm.ErrNotImplemented
	}
	return s.h.GetHash()
}

// Concurrent returns a boolean flag indicating whether the state machine is
// capable of taking concurrent snapshot.
func (s *ConcurrentStateMachine) Concurrent() bool {
	return true
}

// OnDisk returns a boolean flag indicating whether this is a on disk state
// machine.
func (s *ConcurrentStateMachine) OnDisk() bool {
	return false
}

// StateMachineType returns the type of the state machine.
func (s *ConcurrentStateMachine) Type() pb.StateMachineType {
	return pb.ConcurrentStateMachine
}

// ITestFS is an interface implemented by test SMs.
type ITestFS interface {
	SetTestFS(fs config.IFS)
}

// OnDiskStateMachine is the type to represent an on disk state machine.
type OnDiskStateMachine struct {
	sm     sm.IOnDiskStateMachine
	h      sm.IHash
	na     sm.IExtended
	opened bool
}

// NewOnDiskStateMachine creates and returns an on disk state machine.
func NewOnDiskStateMachine(s sm.IOnDiskStateMachine) *OnDiskStateMachine {
	r := &OnDiskStateMachine{sm: s}
	h, ok := s.(sm.IHash)
	if ok {
		r.h = h
	}
	na, ok := s.(sm.IExtended)
	if ok {
		r.na = na
	}
	return r
}

// SetTestFS injects the specified fs to the test SM.
func (s *OnDiskStateMachine) SetTestFS(fs config.IFS) {
	tfs, ok := s.sm.(ITestFS)
	if ok {
		plog.Infof("the underlying SM support test fs injection")
		tfs.SetTestFS(fs)
	}
}

// Open opens the state machine.
func (s *OnDiskStateMachine) Open(stopc <-chan struct{}) (uint64, error) {
	if s.opened {
		panic("Open called more than once on OnDiskStateMachine")
	}
	s.opened = true
	applied, err := s.sm.Open(stopc)
	if err != nil {
		return 0, err
	}
	return applied, nil
}

// Update updates the state machine.
func (s *OnDiskStateMachine) Update(entries []sm.Entry) ([]sm.Entry, error) {
	if !s.opened {
		panic("Update called before Open")
	}
	return s.sm.Update(entries)
}

// Lookup queries the state machine.
func (s *OnDiskStateMachine) Lookup(query interface{}) (interface{}, error) {
	if !s.opened {
		panic("Lookup called when not opened")
	}
	return s.sm.Lookup(query)
}

// NALookup queries the state machine.
func (s *OnDiskStateMachine) NALookup(query []byte) ([]byte, error) {
	if s.na == nil {
		return nil, sm.ErrNotImplemented
	}
	return s.na.NALookup(query)
}

// Sync synchronizes all in-core state with that on disk.
func (s *OnDiskStateMachine) Sync() error {
	if !s.opened {
		panic("Sync called when not opened")
	}
	return s.sm.Sync()
}

// Prepare makes preparations for taking concurrent snapshot.
func (s *OnDiskStateMachine) Prepare() (interface{}, error) {
	if !s.opened {
		panic("PrepareSnapshot called when not opened")
	}
	return s.sm.PrepareSnapshot()
}

// Save saves the snapshot.
func (s *OnDiskStateMachine) Save(ctx interface{},
	w io.Writer, fc sm.ISnapshotFileCollection, stopc <-chan struct{}) error {
	if !s.opened {
		panic("SaveSnapshot called when not opened")
	}
	return s.sm.SaveSnapshot(ctx, w, stopc)
}

// Recover recovers the state machine from a snapshot.
func (s *OnDiskStateMachine) Recover(r io.Reader,
	fs []sm.SnapshotFile, stopc <-chan struct{}) error {
	if !s.opened {
		panic("RecoverFromSnapshot called when not opened")
	}
	return s.sm.RecoverFromSnapshot(r, stopc)
}

// Close closes the state machine.
func (s *OnDiskStateMachine) Close() error {
	return s.sm.Close()
}

// GetHash returns the uint64 hash value representing the state of a state
// machine.
func (s *OnDiskStateMachine) GetHash() (uint64, error) {
	if s.h == nil {
		return 0, sm.ErrNotImplemented
	}
	return s.h.GetHash()
}

// Concurrent returns a boolean flag indicating whether the state machine is
// capable of taking concurrent snapshot.
func (s *OnDiskStateMachine) Concurrent() bool {
	return true
}

// OnDisk returns a boolean flag indicating whether this is an on disk state
// machine.
func (s *OnDiskStateMachine) OnDisk() bool {
	return true
}

// Type returns the type of the state machine.
func (s *OnDiskStateMachine) Type() pb.StateMachineType {
	return pb.OnDiskStateMachine
}
