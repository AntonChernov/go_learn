package main

import (
	"context"
	_ "encoding/json"
	"flag"
	"fmt"
	"os/signal"
	"time"

	_ "html/template"
	"log"
	"net/http"
	"os"

	// logMiddlewar "utils"
	"github.com/gorilla/mux"
	// "github.com/mattn/go-sqlite3"

	// utl "siteparser/utils"
	sparsapi "go_learn/api"
	utl "go_learn/utils"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + " [DEBUG] " + string(bytes))
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr

		log.Printf("%s - - [%s] %q %q %q",
			addr,
			time.Now().Format("2006-01-02T15:04:05Z07:00"),
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
			r.Referer(),
			r.UserAgent())
		// Do stuff here
		// log.Printf("Method: %v URI: %v", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	var hostPort string
	// Add all CLI arguments after this comment and before flag.Parse()

	hostPortCLI := flag.String("hostport", "127.0.0.1:8800", "Host:port for start the server!")
	flag.Parse()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/", sparsapi.HelloHandler).Methods("GET").Name("testing")
	router.HandleFunc("/test/", sparsapi.TestRequestHandler).Methods("GET").Name("testHandler")
	router.HandleFunc("/test-db/", sparsapi.GetUserEmails).Methods("GET").Name("listuseremails")
	// routers here

	router.Use(logMiddleware)

	hpGlVar := os.Getenv("SERVER")

	if hpGlVar != "" {
		hostPort = hpGlVar
	} else {
		hostPort = *hostPortCLI
	}

	log.Printf("Servers started at: %v", hostPort)
	// http.HandleFunc("/", HelloHandler)

	srv := &http.Server{
		Addr: hostPort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	utl.DB.Close()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
