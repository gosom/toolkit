package httpext

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"toolkit/pkg/errorsext"
	"toolkit/pkg/logger"
)

type Option func(*HTTPServer) error

type HTTPServer struct {
	srv *http.Server

	addr    string
	domains []string

	port      int
	log       logger.Logger
	certFile  string
	certCache string
	keyFile   string
}

func New(router http.Handler, opts ...Option) (*HTTPServer, error) {
	ans := HTTPServer{}

	for _, opt := range opts {
		if err := opt(&ans); err != nil {
			return nil, errorsext.WithStack(err)
		}
	}

	setupDefaults(&ans)

	srv := &http.Server{
		Addr:              ans.addr,
		Handler:           router,
		ReadTimeout:       5 * time.Second,  //nolint:gomnd // TODO
		WriteTimeout:      10 * time.Second, //nolint:gomnd // TODO
		IdleTimeout:       5 * time.Second,  //nolint:gomnd // TODO
		ReadHeaderTimeout: 5 * time.Second,  //nolint:gomnd // TODO
		MaxHeaderBytes:    1 << 20,          //nolint:gomnd // TODO
	}

	if ans.port == httpsPort {
		if ans.certFile == "" || ans.keyFile == "" {
			autoTLSManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				Cache:      autocert.DirCache(ans.certCache),
				HostPolicy: autocert.HostWhitelist(ans.domains...),
			}
			// https://ssl-config.mozilla.org/#server=go&version=1.22.0&config=intermediate&guideline=5.7
			srv.TLSConfig = &tls.Config{
				MinVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
				CurvePreferences: []tls.CurveID{
					tls.CurveP256,
					tls.X25519,
				},
				GetCertificate: autoTLSManager.GetCertificate,
				NextProtos: []string{
					"h2", "http/1.1", // enable HTTP/2
					acme.ALPNProto, // enable tls-alpn ACME challenges
				},
			}
		}
	}

	ans.srv = srv

	return &ans, nil
}

func WithLogger(log logger.Logger) Option {
	return func(s *HTTPServer) error {
		s.log = log

		return nil
	}
}

func WithDomains(domains ...string) Option {
	return func(s *HTTPServer) error {
		s.domains = domains

		return nil
	}
}

func WithCertFiles(certFile, keyFile string) Option {
	return func(s *HTTPServer) error {
		s.certFile = certFile
		s.keyFile = keyFile

		if _, err := os.Stat(certFile); err != nil {
			return errorsext.WithStack(err)
		}

		if _, err := os.Stat(keyFile); err != nil {
			return errorsext.WithStack(err)
		}

		return nil
	}
}

func WithAddr(addr string) Option {
	return func(s *HTTPServer) error {
		var err error

		s.addr = addr
		s.port, err = extractPortFromAddr(addr)

		return err
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		h.log.Info(ctx, "shutting down http server")

		shutdownCtx, shutdownStop := context.WithTimeout(
			context.WithoutCancel(ctx),
			5*time.Second, //nolint:gomnd // this is a reasonable timeout
		)

		defer shutdownStop()

		go func() {
			<-shutdownCtx.Done()

			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				h.srv.Close()
			}
		}()

		_ = h.srv.Shutdown(shutdownCtx)
	}()

	var err error

	if h.port == httpsPort {
		if h.certFile != "" && h.keyFile != "" {
			h.log.Info(ctx, "starting https server", "addr", h.addr, "certFile", h.certFile, "keyFile", h.keyFile)
			err = h.srv.ListenAndServeTLS(h.certFile, h.keyFile)
		} else {
			h.log.Info(ctx, "starting https server (autotls)", "addr", h.addr)
			err = h.srv.ListenAndServeTLS("", "")
		}
	} else {
		h.log.Info(ctx, "starting https server", "addr", h.addr)
		err = h.srv.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		return errorsext.WithStack(err)
	}

	h.log.Info(ctx, "http server stopped")

	return nil
}

func extractPortFromAddr(addr string) (int, error) {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(portStr)

	return port, errorsext.WithStack(err)
}

const (
	defaultAddr      = ":80"
	defaultCertCache = "/.cache/certcache"
	defaultHost      = "localhost"
	httpsPort        = 443
)

func setupDefaults(s *HTTPServer) {
	if s.log == nil {
		s.log = logger.Default()
	}

	if s.addr == "" {
		s.addr = defaultAddr
	}

	if len(s.domains) == 0 {
		s.domains = append(s.domains, defaultHost)
	}

	if s.certCache == "" {
		s.certCache = defaultCertCache
	}
}
