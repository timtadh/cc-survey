package clones

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

import (
)


type Subgraph struct {
	source string
	dir string
	java template.HTML
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

type Position struct {
	StartLine, EndLine int
	StartColumn, EndColumn int
}

type Positions []*Position

type errorList []error

func (self errorList) Error() string {
	var s []string
	for _, err := range self {
		s = append(s, err.Error())
	}
	return "[" + strings.Join(s, ",") + "]"
}

func LoadSubgraph(source, dir string, pattern bool) (*Subgraph, error) {
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
		source: source,
		dir: dir,
		pattern: pattern,
	}
	err = sg.load(vegPath)
	if err != nil {
		return nil, err
	}
	return sg, nil
}

func (sg *Subgraph) Class() string {
	return sg.V[0].Class()
}

func (sg *Subgraph) PathToJava() string {
	return sg.V[0].PathToJava()
}

func (sg *Subgraph) Java() template.HTML {
	if sg.java != "" {
		return sg.java
	}
	pathFrag := sg.V[0].PathToJava()
	path := filepath.Join(sg.source, pathFrag)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return "Could not load instance"
	}
	defer f.Close()
	positions := make(Positions, 0, len(sg.V))
	min := -1
	max := 0
	for _, v := range sg.V {
		if pathFrag != v.PathToJava() {
			log.Println(fmt.Errorf("instance spread across multiple files"))
			return "Could not load instance"
		}
		p, err := v.Position()
		if err != nil {
			log.Println(err)
			return "Could not load instance"
		}
		positions = append(positions, p)
		if min == -1 || p.StartLine < min {
			min = p.StartLine
		}
		if p.EndLine > max {
			max = p.EndLine
		}
	}
	lines := make([]string, 0, max-min)
	i := 0
	processLines(f, func(l []byte) {
		if i < min - 5 || i > max + 5 {
			i++
			return
		}
		var line string
		if positions.ContainsLine(i+1) {
			line = fmt.Sprintf(`<div class="line">%d: <span class="highlight">%v</span></div>`, i, string(l))
		} else {
			line = fmt.Sprintf(`<div class="line">%d: %v</div>`, i, string(l))
		}
		lines = append(lines, line)
		i++
	})
	sg.java = template.HTML(strings.Join(lines, ""))
	return sg.java
}

func (v *Vertex) Class() string {
	return v.Attrs["class_name"].(string)
}

func (v *Vertex) PackageName() string {
	return v.Attrs["package_name"].(string)
}

func (v *Vertex) PathToJava() string {
	pkg := v.PackageName()
	name := v.Attrs["source_file"].(string)
	return filepath.Join(strings.Replace(pkg, ".", string(filepath.Separator), -1), name)
}

func (positions Positions) ContainsLine(line int) bool {
	for _, p := range positions {
		if line >= p.StartLine && line <= p.EndLine {
			return true
		}
	}
	return false
}

func (v *Vertex) Position() (p *Position, err error) {
	sl, err := v.Attrs["start_line"].(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	el, err := v.Attrs["end_line"].(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	sc, err := v.Attrs["start_column"].(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	ec, err := v.Attrs["end_column"].(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	p = &Position{
		StartLine: int(sl),
		EndLine: int(el),
		StartColumn: int(sc),
		EndColumn: int(ec),
	}
	return p, nil
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

