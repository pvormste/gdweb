package server

import (
	"crypto/tls"
	"mime"
	"net/http"
	"path/filepath"
)

func init() {
	mime.AddExtensionType(".wasm", "application/wasm")
	mime.AddExtensionType(".pck", "application/octet-stream")
}

// crossOriginHeaders sets the headers required for Godot 4 multi-threaded web exports.
func crossOriginHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		h.ServeHTTP(w, r)
	})
}

// New creates an HTTPS server that serves files from dir with Godot-required headers.
func New(dir string, certPEM, keyPEM []byte, addr string) (*http.Server, error) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	handler := crossOriginHeaders(http.FileServer(http.Dir(absDir)))

	return &http.Server{
		Addr:    addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		},
	}, nil
}
