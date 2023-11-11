package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"webserver-go/variables"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if variables.LogPath != "" {
		openLogFile(variables.LogPath)
		fmt.Printf("Logging to %v\n", variables.LogPath)
	} else {
		fmt.Printf("Logging to console\n")
	}

	hostname, _ := getHostname()
	ipAddress, _ := getIp()
	fmt.Printf("Address: %s (%s)\n", hostname, ipAddress)

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/test", testHandler)

	fmt.Printf("Listening on port %v\n", variables.HttpPort)

	err := http.ListenAndServe(fmt.Sprintf(":%d", variables.HttpPort), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func getIp() (string, error) {
	conn, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		return "", error
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	return ipAddress.IP.String(), nil
}

func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Fprintf(w, variables.NotFound404)
		return
	}
	fmt.Fprintf(w, "<h1>Hello World</h1><div>Try the path /test</div>")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/test")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}
