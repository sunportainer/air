package main

import (
	"com.nzair.user/air"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

//定义一个server的struct，然后里面含airMux
type AirMux struct {
}

func (p *AirMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		helloHandler(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func main() {
	//V8Example()
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/form", helloForm)
	log.Printf("[INFO] [http,server] [message: starting HTTPS server on port %s]", 3333)
	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, air.KeyServerAddr, l.Addr().String())
			return ctx
		},
	}
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	hasFirst := req.URL.Query().Has("email")
	reqEmail := req.URL.Query().Get("email")
	hasPassword := req.URL.Query().Has("pwd")
	reqPassword := req.URL.Query().Get("pwd")
	if reqPassword == "" || reqEmail == "" {
		w.Header().Set("x-missing-field", "pwd or username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := fmt.Sprintf("%s: got / request. firstname(%t)=%s, reqPassword(%t)=%s\n",
		ctx.Value(air.KeyServerAddr),
		hasFirst, reqEmail,
		hasPassword, reqPassword)
	io.WriteString(w, resp)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	resp := fmt.Sprintf("%s: got / request. body:\n%s\n",
		ctx.Value(air.KeyServerAddr), body)
	io.WriteString(w, resp)
}

// handler echoes r.URL.Header
func helloForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.PostFormValue("name")
	resp := fmt.Sprintf("%s: got / request. body:\n%s\n",
		ctx.Value(air.KeyServerAddr), name)
	io.WriteString(w, resp)
}
