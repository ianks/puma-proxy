package main

import (
	"flag"
	"github.com/tv42/httpunix"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var socketPath = flag.String("sock", "/tmp/puma-proxy.sock", "Path for the Puma socket")
var listenAddr = flag.String("listen", "localhost:3000", "Address to listen for requests")
var cmd *exec.Cmd

func proxyToUnixSocket(w http.ResponseWriter, req *http.Request) {
	u := &httpunix.Transport{
		DialTimeout:           1000 * time.Millisecond,
		RequestTimeout:        30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}
	u.RegisterLocation("localhost", *socketPath)

	req.URL.Scheme = "http+unix"
	req.URL.Host = "localhost"
	req.Header.Add("X-Forwarded-For", req.RemoteAddr)

	resp, err := u.RoundTrip(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	copyHeader(w.Header(), resp.Header)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func runCommand() {
	args := flag.Args()
	log.Println("Running command: $", strings.Join(args[:], " "))
	execCmd, execArgs := args[0], args[1:]
	cmd = exec.Command(execCmd, execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		log.Fatal("failed to start process: ", err)
	}
}

func setupStaticAssetsServer() {
	args := flag.Args()

	if len(args) == 0 {
		return
	}

	publicFolder, err := os.Stat("./public")
	if err == nil && publicFolder.IsDir() {
		cwd, err := os.Getwd()

		if err == nil {
			log.Println("Serving static assets from:", filepath.Join(cwd, publicFolder.Name()))
		}
		fs := http.FileServer(http.Dir("./public"))
		http.Handle("/assets/", fs)
	}
}

func setupProxy() {
	http.HandleFunc("/", proxyToUnixSocket)
}

func main() {
	flag.Parse()

	runCommand()
	setupStaticAssetsServer()
	setupProxy()

	log.Println("Proxying to unix socket at", *socketPath)
	log.Println("Puma proxy is ready to handle requests at", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))

	if cmd != nil {
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	}
}
