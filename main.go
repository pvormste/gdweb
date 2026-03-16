package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pvormste/gdweb/internal/cert"
	"github.com/pvormste/gdweb/internal/network"
	"github.com/pvormste/gdweb/internal/server"
)

func main() {
	port := flag.Int("port", 8443, "HTTPS port to listen on")
	dir := flag.String("dir", ".", "Directory to serve")
	host := flag.String("host", "0.0.0.0", "Address to bind to")
	open := flag.Bool("open", false, "Open browser on startup")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gdweb: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "gdweb: directory does not exist: %s\n", absDir)
		os.Exit(1)
	}

	dnsNames := []string{"localhost"}
	ipAddresses := []net.IP{
		net.IPv4(127, 0, 0, 1),
		net.IPv6loopback,
	}
	for _, ip := range network.LocalIPs() {
		ipAddresses = append(ipAddresses, ip)
	}

	certPEM, keyPEM, err := cert.Generate(dnsNames, ipAddresses)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gdweb: failed to generate certificate: %v\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	urls := buildURLs(*port, network.LocalIPs())
	srv, err := server.New(absDir, certPEM, keyPEM, addr, urls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gdweb: %v\n", err)
		os.Exit(1)
	}

	printBanner(absDir, *port, network.LocalIPs(), urls)

	if *open {
		go openBrowser(*port)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gdweb: %v\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	tlsLn := tls.NewListener(ln, srv.TLSConfig)
	if err := srv.Serve(tlsLn); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "gdweb: %v\n", err)
		os.Exit(1)
	}
}

func buildURLs(port int, localIPs []net.IP) []string {
	urls := []string{"https://localhost:" + fmt.Sprint(port)}
	for _, ip := range localIPs {
		urls = append(urls, fmt.Sprintf("https://%s:%d", ip.String(), port))
	}
	return urls
}

func printBanner(serveDir string, port int, localIPs []net.IP, urls []string) {
	fmt.Println("gdweb - Godot Web Export Server")
	fmt.Printf("Serving:  %s\n", serveDir)
	fmt.Printf("Port:     %d\n\n", port)
	fmt.Println("  Local:   https://localhost:" + fmt.Sprint(port))
	for _, ip := range localIPs {
		fmt.Printf("  Network: https://%s:%d\n", ip.String(), port)
	}
	if len(urls) > 0 {
		fmt.Printf("\n  QR codes: %s/qr\n", urls[0])
	}
	fmt.Println("\nNote: Self-signed certificate. Accept the browser warning to proceed.")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()
}

func openBrowser(port int) {
	url := fmt.Sprintf("https://localhost:%d", port)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
