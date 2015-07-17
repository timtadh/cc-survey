package clones

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"strconv"
)

type Clone struct {
	dir string
	pr float64
	img []byte
	Pattern *Subgraph
	Instances []*Subgraph
}

func loadCount(dir string) (int, error) {
	countBytes, err := ioutil.ReadFile(filepath.Join(dir, "count"))
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(string(bytes.TrimSpace(countBytes)), 10, 32)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func LoadAll(dir string) ([]*Clone, error) {
	count, err := loadCount(dir)
	if err != nil {
		return nil, err
	}
	clones := make([]*Clone, 0, count)
	for i := 0; i < count; i++ {
		p := filepath.Join(dir, fmt.Sprintf("%d", i))
		clone, err := Load(p)
		if err != nil {
			return nil, err
		}
		clones = append(clones, clone)
	}
	return clones, nil
}

func Load(dir string) (*Clone, error) {
	count, err := loadCount(dir)
	if err != nil {
		return nil, err
	}
	c := &Clone{
		dir: dir,
		Instances: make([]*Subgraph, 0, count),
	}
	c.Pattern, err = c.loadPattern()
	if err != nil {
		return nil, err
	}
	for i := 0; i < count; i++ {
		i, err := c.loadInstance(i)
		if err != nil {
			return nil, err
		}
		c.Instances = append(c.Instances, i)
	}
	return c, nil
}

func (c *Clone) loadPattern() (*Subgraph, error) {
	return LoadSubgraph(c.dir, true)
}

func (c *Clone) loadInstance(i int) (*Subgraph, error) {
	p := filepath.Join(c.dir, "instances", fmt.Sprintf("%d", i))
	return LoadSubgraph(p, false)
}

