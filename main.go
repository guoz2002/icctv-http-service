package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// 最简健康检查接口：用于判断服务是否存活
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":8080"
	log.Printf("OrangePi HTTP service listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
