package clones

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type Clone struct {
	dir string
	pr float64
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

func (c *Clone) Img() (f *os.File, modtime time.Time, err error) {
	dot := filepath.Join(c.dir, "pattern.dot")
	path := filepath.Join(c.dir, "pattern.png")
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err = c.generateImg(dot, path)
		if err != nil {
			return nil, modtime, err
		}
		modtime = time.Now().UTC()
	} else if err != nil {
		return nil, modtime, err
	} else {
		modtime = fi.ModTime()
	}
	f, err = os.Open(path)
	if err != nil {
		return nil, modtime, err
	}
	return f, modtime, nil
}

func (c *Clone) generateImg(src, out string) error {
	dot, err := exec.LookPath("dot")
	if err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	cmd := exec.Cmd{
		Path: dot,
		Args: []string{"dot", "-Tpng", src},
		Stdout: f,
		Stderr: os.Stderr,
	}
	err = cmd.Run()
	f.Close()
	if err != nil {
		os.Remove(out)
		return err
	}
	return nil
}

