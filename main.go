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
	sparsapi "siteparser/api"
)

// type TestsStruct struct {
// 	URL       *string `json:"url"`
// }

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr

		log.Printf("%s - - [%s] %q %q %q",
			addr,
			time.Now().Format("02/Jan/2006:15:04:05.999 +0200"),
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
	var hostPort string = "127.0.0.1:8800"
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/", sparsapi.HelloHandler).Methods("GET").Name("testing")
	router.HandleFunc("/test/", sparsapi.TestRequestHandler).Methods("GET").Name("testHandler")
	router.HandleFunc("/test-db/", sparsapi.GetUserEmails).Methods("GET").Name("listuseremails")
	// routers here

	router.Use(logMiddleware)
	// db, err := utl.PostgressDBConnection()

	// env := &Env{db: db}

	hpGlVar := os.Getenv("SERVER")

	if hpGlVar != "" {
		hostPort = hpGlVar
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
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

	// router.HandleFunc("/", HelloHandler).Methods("GET").Name("testing")
	// http.Handle("/", router)
	// fmt.Println(http.ListenAndServe(hostPort, nil))

}

// func MainPage(w http.ResponseWriter, r *http.Request){

// }

// func HelloJson struct{

// }
