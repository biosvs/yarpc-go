// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package client

import "fmt"

// Status represents the result of running a behavior.
type Status string

// Different valid Statuses.
const (
	Passed  Status = "passed"
	Skipped        = "skipped"
	Failed         = "failed"
)

// Params provides access to the behavior parameters.
type Params interface {
	Param(name string) string
}

// Behavior provides access to parameters and a means to send back results for
// a behavior implementation.
type Behavior interface {
	Params

	// The following methods specify the state of the behavior. They may be
	// called any number of times to emit multiple entries for a single call.

	Skip(reason string)
	Fail(message string)
	Pass(output string)

	Skipf(format string, args ...interface{})
	Failf(format string, args ...interface{})
	Passf(format string, args ...interface{})
}

//////////////////////////////////////////////////////////////////////////////

// EntryBuilder builds entries for a behavior's results based on the calls made
// by the test.
type EntryBuilder interface {
	Skip(reason string) interface{}
	Fail(message string) interface{}
	Pass(output string) interface{}
}

// BasicEntry is the most basic form of an entry for a behavior test.
type BasicEntry struct {
	Status Status `json:"status"`
	Output string `json:"output"`
}

// basicEntryBuilder is an EntryBuilder that builds BasicEntry objects.
type basicEntryBuilder struct{}

// BasicEntryBuilder is a simple EntryBuilder that includes very little
// information.
var BasicEntryBuilder EntryBuilder = basicEntryBuilder{}

// Skip for basicEntryBuilder.
func (b basicEntryBuilder) Skip(reason string) interface{} {
	return BasicEntry{Status: Skipped, Output: reason}
}

// Fail for basicEntryBuilder.
func (b basicEntryBuilder) Fail(message string) interface{} {
	return BasicEntry{Status: Failed, Output: message}
}

// Pass for basicEntryBuilder.
func (b basicEntryBuilder) Pass(output string) interface{} {
	return BasicEntry{Status: Passed, Output: output}
}

//////////////////////////////////////////////////////////////////////////////

// BehaviorTester is the root Behavior.
type BehaviorTester struct {
	Params

	Failed  bool
	Skipped bool
	Entries []interface{}
}

// NewBehavior provides a new Behavior that may be passed into a test to record
// its results.
func (bt *BehaviorTester) NewBehavior(builder EntryBuilder) Behavior {
	return behavior{Params: bt.Params, Tester: bt, Builder: builder}
}

// putEntry records a new entry with this BehaviorTester.
func (bt *BehaviorTester) putEntry(entry interface{}, status Status) {
	switch status {
	case Failed:
		bt.Failed = true
	case Skipped:
		bt.Skipped = true
	default:
		// nothing to do
	}
	bt.Entries = append(bt.Entries, entry)
}

//////////////////////////////////////////////////////////////////////////////

type behavior struct {
	Params

	Tester  *BehaviorTester
	Builder EntryBuilder
}

func (b behavior) Skip(reason string) {
	b.Tester.putEntry(b.Builder.Skip(reason), Skipped)
}

func (b behavior) Fail(message string) {
	b.Tester.putEntry(b.Builder.Fail(message), Failed)
}

func (b behavior) Pass(output string) {
	b.Tester.putEntry(b.Builder.Pass(output), Passed)
}

func (b behavior) Skipf(format string, args ...interface{}) {
	b.Skip(fmt.Sprintf(format, args...))
}
func (b behavior) Failf(format string, args ...interface{}) {
	b.Fail(fmt.Sprintf(format, args...))
}
func (b behavior) Passf(format string, args ...interface{}) {
	b.Pass(fmt.Sprintf(format, args...))
}
