package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"webserver-go/variables"
)

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-c
		// Run Cleanup
		slog.Error("Stopping...", "signal", sig)
		os.Exit(1)
	}()
}

func main() {
	fmt.Println("PID", os.Getpid())

	// Define and check parameters
	httpPort := flag.Int("port", variables.HttpPortDefault, "Listening TCP Port")
	logFile := flag.String("log", "", "JSON logfile that is used for logging")
	flag.Parse()
	// End of: Define and check parameters

	if *logFile != "" {
		openLogFile(*logFile)
		fmt.Printf("Logging to %v\n", *logFile)
		slog.Info("Starting...")
	} else {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
		slog.SetDefault(logger)
		fmt.Printf("Logging to console\n")
	}

	hostname, _ := getHostname()
	ipAddress, _ := getIp()
	fmt.Printf("Address: %s (%s)\n", hostname, ipAddress)
	fmt.Printf("Port: %d\n", *httpPort)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/test", testHandler)

	slog.Info("Webserver listening", "port", *httpPort)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), logRequest(http.DefaultServeMux))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
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
		w.WriteHeader(404)
		fmt.Fprintf(w, variables.NotFound404)
		//slog.Error("Resquest 404", "RemoteAddr", r.RemoteAddr, "Method", r.Method, "URL", r.URL)
		return
	}
	fmt.Fprintf(w, "<h1>Hello World</h1><div>Try the path /test</div>")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/test")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		if variables.LogRequests {
			slog.Info("Request", "RemoteAddr", r.RemoteAddr, "Method", r.Method, "URL", r.URL)
		}
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			slog.Error("OpenLogfile: os.OpenFile", "ERROR", err)
			os.Exit(1)
		}
		logger := slog.New(slog.NewJSONHandler(lf, &slog.HandlerOptions{}))
		slog.SetDefault(logger)
	}
}
