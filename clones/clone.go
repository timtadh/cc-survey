package clones

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
	"sync"
)

type Clone struct {
	source string
	dir string
	ext_id int
	pr float64
	Pattern *Subgraph
	Instances []*Subgraph
	lock sync.Mutex
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

func LoadAll(dir, source string) ([]*Clone, error) {
	count, err := loadCount(dir)
	if err != nil {
		return nil, err
	}
	clones := make([]*Clone, 0, count)
	for i := 0; i < count; i++ {
		p := filepath.Join(dir, fmt.Sprintf("%d", i))
		clone, err := Load(i, p, source)
		if err != nil {
			log.Println("skipping clone", i, err)
			continue
		}
		clones = append(clones, clone)
	}
	return clones, nil
}

func Load(ext_id int, dir, source string) (*Clone, error) {
	count, err := loadCount(dir)
	if err != nil {
		return nil, err
	}
	c := &Clone{
		ext_id: ext_id,
		dir: dir,
		source: source,
		Instances: make([]*Subgraph, 0, count),
	}
	err = c.loadPr()
	if err != nil {
		return nil, err
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
	return LoadSubgraph(c.source, c.dir, true)
}

func (c *Clone) loadInstance(i int) (*Subgraph, error) {
	p := filepath.Join(c.dir, "instances", fmt.Sprintf("%d", i))
	return LoadSubgraph(c.source, p, false)
}

func (c *Clone) loadPr() (error) {
	prBytes, err := ioutil.ReadFile(filepath.Join(c.dir, "pattern.pr"))
	if err != nil {
		return err
	}
	pr, err := strconv.ParseFloat(string(bytes.TrimSpace(prBytes)), 64)
	if err != nil {
		return err
	}
	c.pr = pr
	return nil
}

func (c *Clone) ExtId() int {
	return c.ext_id
}

func (c *Clone) Dir() string {
	return c.dir
}

func (c *Clone) Pr() float64 {
	return c.pr
}

func (c *Clone) Img() (f *os.File, modtime time.Time, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	dot := filepath.Join(c.dir, "pattern.dot")
	path := filepath.Join(c.dir, "pattern.png")
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err = generateImg(dot, path)
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

func generateImg(src, out string) error {
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

