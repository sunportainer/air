package main

import (
	"com.nzair.user/air"
	"com.nzair.user/handler"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
)

func main() {
	mux, err := handler.NewHandler()
	if err != nil {
		log.Printf("[ERR] [http,server] [message: Failed to start HTTPS server on port %s]", 3333)
		return
	}
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
			log.Printf("server closed\n")
		} else if err != nil {
			log.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()
}
