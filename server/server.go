package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var redirectUrlTemplate = os.Getenv("redirectUrlTemplate")

func redirecter(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")

	if state == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	parts := strings.Split(state, ":")

	if len(parts) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	redirectUrl := strings.ReplaceAll(redirectUrlTemplate, "*", parts[0]) + r.URL.Path + "?state=" + parts[1]

	// copy over all query parameters
	for k, v := range r.URL.Query() {
		if k == "state" { // exclude state, as already included above
			continue
		}

		for _, value := range v {
			redirectUrl += "&" + k + "=" + value
		}
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func main() {
	if redirectUrlTemplate == "" {
		fmt.Sprintln("redirectUrlTemplate environment variable not set")
		os.Exit(1)
	}

	port := os.Getenv("port");

	//TODO maybe also support file pointers?
	certPem := os.Getenv("certPem")
	keyPem := os.Getenv("keyPem")

	if certPem != "" && keyPem != "" {
		cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
		if err != nil {
			log.Fatal(err)
		}

		// CHECKME not 100% sure this structure: went with extending http.Server to allow passing certPemand keyPem strings, insetad of file pointers
		router := http.NewServeMux()
		router.Handle("/", func() http.Handler {
			return http.HandlerFunc(redirecter)
		}())

		if port == "" {
			port = "443"
		}

		cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
		srv := &http.Server{
			TLSConfig: cfg,
			Handler:   router,
			Addr: ":" + port,
		}

		log.Fatal(srv.ListenAndServeTLS("", ""))
	} else {
		http.HandleFunc("/", redirecter)

		if port == "" {
			port = "8080"
		}

		log.Fatal(http.ListenAndServe(":" + port), nil))
	}
}
