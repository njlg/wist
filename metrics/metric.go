// Package metrics ...
package metrics

import (
	"errors"
	"sync"
)

// Metric ...
type Metric struct {
	mutex sync.Mutex
	list  []float64
	head  int
}

// New ...
func New(size int) *Metric {
	return &Metric{
		list: make([]float64, 0, size),
	}
}

// Add ...
func (m *Metric) Add(val float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// if the list is at capacity, drop head
	if len(m.list) >= cap(m.list) {
		m.list = m.list[1:]
	}

	m.list = append(m.list, val)
}

var data map[string]*Metric

func init() {

	data = make(map[string]*Metric, 20)

	size := 120
	data["link.quality"] = New(size)
	data["link.level"] = New(size)
	data["link.noise"] = New(size)
	data["discard_packets.nwid"] = New(size)
	data["discard_packets.crypt"] = New(size)
	data["discard_packets.frag"] = New(size)
	data["discard_packets.retry"] = New(size)
	data["discard_packets.misc"] = New(size)

	data["bitrate"] = New(size)
}

// Add ...
func Add(name string, val float64) {
	if d, ok := data[name]; ok {
		d.Add(val)
	}
}

// Get ...
func Get(name string) (float64, error) {
	d, ok := data[name]

	if !ok {
		return 0, errors.New("slkfjldkf")
	}

	return d.list[len(d.list)-1], nil
}

// GetAll ...
func GetAll(name string) ([]float64, error) {
	d, ok := data[name]

	if !ok {
		return nil, errors.New("slkfjldkf")
	}

	return d.list, nil
}
