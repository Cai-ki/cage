// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// 主页
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 可选：检测是否来自 HTTPS（通过 Cloudflare）
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			fmt.Fprintf(w, "✅ Hello from Go! You're on HTTPS via Cloudflare.\n")
		} else {
			fmt.Fprintf(w, "⚠️ You're on HTTP (should redirect to HTTPS)\n")
		}
	})

	// 自动跳转根域名 HTTP → www HTTPS（可选但推荐）
	mux.HandleFunc("/redirect-root", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.caiki.net/", http.StatusMovedPermanently)
	})

	log.Println("Go server starting on :80")
	log.Fatal(http.ListenAndServe(":80", mux))
}
