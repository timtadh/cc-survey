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

import (
	"github.com/timtadh/cc-survey/views"
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
    -l, --listen=<addr>:<port>          what to listen on
                                        default: 0.0.0.0:80
    -a, --assets=<path>                 path to asset dir
    -c, --clones=<path>                 path to clones dir
`

func Usage(code int) {
	fmt.Fprintln(os.Stderr, UsageMessage)
	if code == 0 {
		fmt.Fprintln(os.Stdout, ExtendedMessage)
		code = ErrorCodes["usage"]
	}
	os.Exit(code)
}

func main() {
	_, optargs, err := getopt.GetOpt(
		os.Args[1:],
		"hl:a:c:",
		[]string{ "help", "listen=", "assets", "clones=" },
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing command line flags", err)
		Usage(ErrorCodes["opts"])
	}

	listen := "0.0.0.0:80"
	assets := ""
	clones := ""
	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			Usage(0)
			os.Exit(0)
		case "-l", "--listen":
			listen = oa.Arg()
		case "-a", "--assets":
			assets = oa.Arg()
		case "-c", "--clones":
			clones = oa.Arg()
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag '%v'\n", oa.Opt())
			Usage(ErrorCodes["opts"])
		}
	}

	if assets == "" {
		fmt.Fprintln(os.Stderr, "You must supply a path to the assets")
		Usage(ErrorCodes["opts"])
	}
	
	if clones == "" {
		fmt.Fprintln(os.Stderr, "You must supply a path to the clones")
		Usage(ErrorCodes["opts"])
	}

	handler := views.Routes(assets, clones)

	server := &http.Server{
		Addr: listen,
		Handler: handler,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		TLSConfig: nil,
		TLSNextProto: nil,
		ConnState: nil,
		ErrorLog: nil,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

