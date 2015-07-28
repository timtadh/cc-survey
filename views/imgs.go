package views

import (
	"log"
	"net/http"
	"strconv"
)

import (
)

import (
)


func (v *Views) PatternImg(c *Context) {
	cid, err := strconv.Atoi(c.p.ByName("clone"))
	if err != nil || cid < 0 || cid >= len(v.clones) {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("invalid clone id")
		}
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed parameter submitted"))
		return
	}
	f, modtime, err := v.clones[cid].Img()
	if err != nil {
		log.Println(err)
		c.rw.WriteHeader(500)
		c.rw.Write([]byte("could not generate img"))
		return
	}
	defer f.Close()
	http.ServeContent(c.rw, c.r, "pattern.png", modtime, f)
}

func (v *Views) InstanceImg(c *Context) {
	cid, err := strconv.Atoi(c.p.ByName("clone"))
	if err != nil || cid < 0 || cid >= len(v.clones) {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("invalid clone id")
		}
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed parameter submitted"))
		return
	}
	clone := v.clones[cid]
	iid, err := strconv.Atoi(c.p.ByName("instance"))
	if err != nil || iid < 0 || iid >= len(clone.Instances) {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("invalid instance id")
		}
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed parameter submitted"))
		return
	}
	f, modtime, err := clone.Instances[iid].Img()
	if err != nil {
		log.Println(err)
		c.rw.WriteHeader(500)
		c.rw.Write([]byte("could not generate img"))
		return
	}
	defer f.Close()
	http.ServeContent(c.rw, c.r, "pattern.png", modtime, f)
}

