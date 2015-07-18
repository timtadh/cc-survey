package clones

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

import (
)


type Subgraph struct {
	dir string
	java string
	jimple string
	pattern bool
	V []*Vertex
	E []*Edge
}

type Vertex struct {
	Id int
	Label string
	Attrs map[string]interface{}
}

type Edge struct {
	Src, Targ int
	Label string
}

type errorList []error

func (self errorList) Error() string {
	var s []string
	for _, err := range self {
		s = append(s, err.Error())
	}
	return "[" + strings.Join(s, ",") + "]"
}

func LoadSubgraph(dir string, pattern bool) (*Subgraph, error) {
	var vegPath string
	if pattern {
		vegPath = filepath.Join(dir, "pattern.veg")
	} else {
		vegPath = filepath.Join(dir, "embedding.veg")
	}
	_, err := os.Stat(vegPath)
	if err != nil {
		return nil, err
	}
	sg := &Subgraph{
		dir: dir,
		pattern: pattern,
	}
	err = sg.load(vegPath)
	if err != nil {
		return nil, err
	}
	return sg, nil
}

func (sg *Subgraph) load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	sg.V = make([]*Vertex, 0, 10)
	sg.E = make([]*Edge, 0, 10)
	var errors errorList
	vids := make(map[int]int)
	processLines(f, func(line []byte) {
		if len(line) == 0 || !bytes.Contains(line, []byte("\t")) {
			return
		}
		line_type, data := parseLine(line)
		switch line_type {
		case "vertex":
			if err := sg.loadVertex(vids, data); err != nil {
				errors = append(errors, err)
			}
		case "edge":
			if err := sg.loadEdge(vids, data); err != nil {
				errors = append(errors, err)
			}
		default:
			errors = append(errors, fmt.Errorf("Unknown line type %v", line_type))
			return
		}
	})
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (sg *Subgraph) loadVertex(vids map[int]int, data []byte) (err error) {
	obj, err := parseJson(data)
	if err != nil {
		return err
	}
	_id, err := obj["id"].(json.Number).Int64()
	if err != nil {
		return err
	}
	label := strings.TrimSpace(obj["label"].(string))
	id := int(_id)
	vertex := &Vertex{
		Id: id,
		Label: label,
		Attrs: make(map[string]interface{}),
	}
	vids[id] = len(sg.V)
	sg.V = append(sg.V, vertex)
	for k, v := range obj {
		if k == "id" || k == "label" {
			continue
		}
		vertex.Attrs[k] = v
	}
	return nil
}

func (sg *Subgraph) loadEdge(vids map[int]int, data []byte) (err error) {
	obj, err := parseJson(data)
	if err != nil {
		return err
	}
	_src, err := obj["src"].(json.Number).Int64()
	if err != nil {
		return err
	}
	_targ, err := obj["targ"].(json.Number).Int64()
	if err != nil {
		return err
	}
	src := int(_src)
	targ := int(_targ)
	label := strings.TrimSpace(obj["label"].(string))
	edge := &Edge{Src: src, Targ: targ, Label: label}
	sg.E = append(sg.E, edge)
	return nil
}

func parseJson(data []byte) (obj map[string]interface{}, err error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func parseLine(line []byte) (line_type string, data []byte) {
	split := bytes.Split(line, []byte("\t"))
	return strings.TrimSpace(string(split[0])), bytes.TrimSpace(split[1])
}

func processLines(reader io.Reader, process func([]byte)) {

	const SIZE = 4096

	read_chunk := func() (chunk []byte, closed bool) {
		chunk = make([]byte, 4096)
		if n, err := reader.Read(chunk); err == io.EOF {
			return nil, true
		} else if err != nil {
			panic(err)
		} else {
			return chunk[:n], false
		}
	}

	parse := func(buf []byte) (obuf, line []byte, ok bool) {
		for i := 0; i < len(buf); i++ {
			if buf[i] == '\n' {
				line = buf[:i+1]
				obuf = buf[i+1:]
				return obuf, line, true
			}
		}
		return buf, nil, false
	}

	var buf []byte
	read_line := func() (line []byte, closed bool) {
		ok := false
		buf, line, ok = parse(buf)
		for !ok {
			chunk, closed := read_chunk()
			if closed || len(chunk) == 0 {
				return buf, true
			}
			buf = append(buf, chunk...)
			buf, line, ok = parse(buf)
		}
		return line, false
	}

	var line []byte
	closed := false
	for !closed {
		line, closed = read_line()
		process(line)
	}
}
