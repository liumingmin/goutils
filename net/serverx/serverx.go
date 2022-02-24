package serverx

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ServerX struct {
	h1s *http.Server
}

func (s *ServerX) waitExitSignal() {
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		fmt.Fprintf(os.Stderr, "Receive %v signal...\n", sig)

		if s.h1s != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			if err := s.h1s.Shutdown(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "Shutdown server error: %v \n", err)

				os.Kill.Signal()
			}
		}
		os.Exit(0)
	}()
}

func (s *ServerX) Run(addr string, handler http.Handler) {
	fmt.Fprintf(os.Stderr, "LISTEN: %s\n", addr)

	//默认开启http2
	http1x := os.Getenv("SERVER_H1X")
	if http1x == "" {
		fmt.Fprintf(os.Stderr, "http2 init\n")
		h2s := &http2.Server{}
		s.h1s = &http.Server{
			Addr:    addr,
			Handler: h2c.NewHandler(handler, h2s),
		}
	} else {
		s.h1s = &http.Server{
			Addr:    addr,
			Handler: handler,
		}
	}

	s.waitExitSignal()

	e := s.h1s.ListenAndServe()
	if e != nil && e != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server stoped... error: %v\n", e)
	}
}
