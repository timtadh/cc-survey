package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

import (
	"github.com/timtadh/getopt"
)


var ErrorCodes map[string]int = map[string]int{
	"usage":   0,
	"version": 2,
	"opts":    3,
	"badint":  5,
	"baddir":  6,
	"badfile": 7,
}

var UsageMessage string = "cc-survey --help"
var ExtendedMessage string = `
cc-survey

Options
    -h, --help                          view this message
`

func Usage(code int) {
	fmt.Fprintln(os.Stderr, UsageMessage)
	if code == 0 {
		fmt.Fprintln(os.Stdout, ExtendedMessage)
		code = ErrorCodes["usage"]
	} else {
		fmt.Fprintln(os.Stderr, "Try -h or --help for help")
	}
	os.Exit(code)
}


func main() {
	_, optargs, err := getopt.GetOpt(
		os.Args[1:],
		"hl:",
		[]string{ "help", "listen=" },
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing command line flags", err)
		Usage(ErrorCodes["opts"])
	}

	listen := "0.0.0.0:80"
	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			Usage(0)
			os.Exit(0)
		case "-l", "--listen":
			listen = oa.Arg()
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag '%v'\n", oa.Opt())
			Usage(ErrorCodes["opts"])
		}
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr: listen,
		Handler: mux,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		TLSConfig: nil,
		TLSNextProto: nil,
		ConnState: nil,
		ErrorLog: nil,
	}

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("hello"))
	})

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

