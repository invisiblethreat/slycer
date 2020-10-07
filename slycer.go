package slycer

import (
	"fmt"
)

type Offset struct {
	Current int
	Max     int
	Saved   map[string]int
	Skipped map[string]int
}

// NewOffsetTracker is used to initialze maps
func NewOffsetTracker() Offset {
	return Offset{Saved: make(map[string]int), Skipped: make(map[string]int)}
}

func (o *Offset) SetMax(m int) {
	o.Max = m
}

func (o *Offset) ExceedsMax(e int) bool {
	if o.Current+e > o.Max {
		return true
	}
	return false
}

// retun the current index
func (o *Offset) Index() int {
	return o.Current
}

// advance the index and return the new value
func (o *Offset) Step(s int) int {
	o.Current += s
	return o.Current
}

// advance the index in order to skip fields that we don't care about.
func (o *Offset) Skip(a int) {
	o.Current += a
}

// skipNote allows for the capture of offsets that you might not care about, but
//can then later dump as a map of constant offset locations.
func (o *Offset) SkipNote(a int, n string) {
	o.Skipped[n] = o.Current
	o.Current += a
}

func (o *Offset) SaveCurrent(s string) {
	o.Saved[s] = o.Current
}

func (o *Offset) ShowSaved() {
	o.showMap(o.Saved)
}

func (o *Offset) ShowSkipped() {
	o.showMap(o.Skipped)
}

func (o *Offset) showMap(m map[string]int) {
	fmt.Printf("const (\n")
	for k, v := range m {
		fmt.Printf("\t%s = %d\n", k, v)
	}
	fmt.Printf(")\n")
}
func (o *Offset) LoadSaved(s string) {
	if val, hit := o.Saved[s]; hit {
		o.Saved["indexBeforeLoad"] = o.Current
		o.Current = val
	} else {
		fmt.Printf("No save location called \"%s\"", s)
	}
}

func (o *Offset) GetSaved(s string) int {
	if val, hit := o.Saved[s]; hit {
		return val
	} else {
		fmt.Printf("No save location called \"%s\"", s)
		return o.Current
	}
}

func (o *Offset) RestorePrevious() {
	if val, hit := o.Saved["indexBeforeLoad"]; hit {
		swap := o.Current
		o.Current = val
		o.Saved["indexBeforeLoad"] = swap
	} else {
		fmt.Printf("No previous value to restore")
	}
}
