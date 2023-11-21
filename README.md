# Final Project Pelatihan Scalable Web Service with Golang

Implementasi Websocket sederhana pada golang untuk aplikasi web chat sederhana.

[Link](https://github.com/asrofilfachrulr/simple-webchat-frontend) Repo Frontend

# Penggunaan

Install package yang dibutuhkan dan jalankan webserver lewat terminal

```bash
go get
go run main.go

```

# Objektif

### Utama
- [x] Buatlah Chat dengan websocket golang yang dapat berinteraksi dengan lebih dari 2 user  dalam satu room chat.

### Opsional
- [x] Gorilla Websocket
- [x] Sync
- [x] encoding/json
- [x] Goroutine and channel
- [x] Gorilla Mux

### Pengembangan Mandiri  
- [x] Dockerfile
- [x] CI-CD dengan [github actions](https://github.com/asrofilfachrulr/simple-webchat-backend/blob/main/.github/workflows/main.yml) untuk build docker image dan push ke [dockerhub](https://hub.docker.com/repository/docker/cucumber1420/go-websocket/general)
