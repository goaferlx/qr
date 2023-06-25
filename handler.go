package main

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gorilla/mux"

	"golang.org/x/exp/slog"
)

type Response struct {
	Text  string
	Code  string
	Error string
}

type Handler struct {
	http.Handler
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	h := Handler{
		logger: logger,
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", h.ShowForm(nil)).Methods(http.MethodGet)
	router.HandleFunc("/", h.GenerateCode()).Methods(http.MethodPost)
	h.Handler = router
	return &h
}

// GET /qr
func (h *Handler) ShowForm(data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("").ParseFiles("qr.html")
		if err != nil {
			h.logger.Error("failed to parse template files", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "qr", data); err != nil {
			h.logger.Error("failed to execute template", "error", err, "template", "qr")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// POST /qr
func (h *Handler) GenerateCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Info("failed to parse form", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			h.ShowForm(Response{Error: "Something went wrong"}).ServeHTTP(w, r)
			return
		}

		code, err := qr.Encode(r.PostFormValue("text"), qr.H, qr.Auto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Info("failed to encode to qr", "error", err)
			h.ShowForm(Response{Error: "Something went wrong"}).ServeHTTP(w, r)
			return
		}

		scaledCode, err := barcode.Scale(code, 600, 600)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Info("failed to scale qr code", "error", err)
			h.ShowForm(Response{Error: "Something went wrong"}).ServeHTTP(w, r)
			return
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, scaledCode); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to encode to png", "error", err)
			h.ShowForm(Response{Error: "Something went wrong"}).ServeHTTP(w, r)
			return
		}

		encodedQR := base64.RawStdEncoding.EncodeToString(buf.Bytes())
		h.ShowForm(Response{Text: r.PostFormValue("text"), Code: encodedQR}).ServeHTTP(w, r)
	}
}
