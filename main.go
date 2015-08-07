package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
    -s, --src=<path>                    path to the source
    --private-ssl-key=<path>
    --ssl-cert=<path>
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
		"hl:a:c:s:",
		[]string{ "help", "listen=", "assets", "clones=", "src=",
		          "private-ssl-key=", "ssl-cert=" },
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing command line flags", err)
		Usage(ErrorCodes["opts"])
	}

	listen := "0.0.0.0:80"
	assets := ""
	clones := ""
	source := ""
	ssl_key := ""
	ssl_cert := ""
	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			Usage(0)
			os.Exit(0)
		case "-l", "--listen":
			listen = oa.Arg()
		case "-a", "--assets":
			assets, err = filepath.Abs(oa.Arg())
			if err != nil {
				fmt.Fprintf(os.Stderr, "assets path was bad: %v", err)
				Usage(ErrorCodes["opts"])
			}
		case "-c", "--clones":
			clones, err = filepath.Abs(oa.Arg())
			if err != nil {
				fmt.Fprintf(os.Stderr, "clones path was bad: %v", err)
				Usage(ErrorCodes["opts"])
			}
		case "-s", "--src":
			source, err = filepath.Abs(oa.Arg())
			if err != nil {
				fmt.Fprintf(os.Stderr, "source path was bad: %v", err)
				Usage(ErrorCodes["opts"])
			}
		case "--private-ssl-key":
			ssl_key, err = filepath.Abs(oa.Arg())
			if err != nil {
				fmt.Fprintf(os.Stderr, "private-ssl-key path was bad: %v", err)
				Usage(ErrorCodes["opts"])
			}
			_, err = os.Stat(ssl_key)
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "private-ssl-key path does not exist. %v", ssl_key)
				Usage(ErrorCodes["opts"])
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "private-ssl-key path was bad: %v", err)
				Usage(ErrorCodes["opts"])
			}
		case "--ssl-cert":
			log.Println("ssl-cert", oa.Arg())
			ssl_cert, err = filepath.Abs(oa.Arg())
			if err != nil {
				fmt.Fprintf(os.Stderr, "ssl-cert path was bad: %v", err)
				Usage(ErrorCodes["opts"])
			}
			_, err = os.Stat(ssl_cert)
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "ssl-cert path does not exist. %v", ssl_cert)
				Usage(ErrorCodes["opts"])
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "ssl-cert path was bad: %v", err)
				Usage(ErrorCodes["opts"])
			}
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

	if source == "" {
		fmt.Fprintln(os.Stderr, "You must supply a path to the source")
		Usage(ErrorCodes["opts"])
	}

	if (ssl_key == "" && ssl_cert != "") || (ssl_key != "" && ssl_cert == "") {
		fmt.Fprintln(os.Stderr, "To use ssl you must supply key and cert")
		Usage(ErrorCodes["opts"])
	}

	handler := views.Routes(assets, clones, source)

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

	if ssl_key == "" {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println(ssl_cert, ssl_key)
		err = server.ListenAndServeTLS(ssl_cert, ssl_key)
		if err != nil {
			log.Panic(err)
		}
	}
}

