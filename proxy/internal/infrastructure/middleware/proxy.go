package middleware

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy/internal/modules/swagger"
	"strings"
)

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/metrics") && !strings.HasPrefix(r.URL.Path, "/debug/pprof") {

			url, err := url.Parse(rp.host + ":" + rp.port)
			if err != nil {
				log.Fatal(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(url)
			r.Host = "hugo:1313"

			if strings.HasPrefix(r.URL.Path, "/swagger") {
				swagger.SwaggerUI(w, r)

				return
			}

			if strings.HasPrefix(r.URL.Path, "/public") {
				http.ServeFile(w, r, "./public/swagger.json")
				return
			}

			proxy.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}

	})

}
