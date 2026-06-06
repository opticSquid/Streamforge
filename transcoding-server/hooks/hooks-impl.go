package hooks

import (
	"log/slog"
	"net/http"
)

func New(manager StreamManager) *Handler {
	return &Handler{manager: manager}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/on_publish", h.postOnly(h.onPublish))
	mux.HandleFunc("/on_publish_done", h.postOnly(h.onPublishDone))
}

// onPublish is called by nginx-rtmp when a stream key starts publishing.
// nginx-rtmp interprets any non-2xx response as rejection — returning 403
func (h *Handler) onPublish(w http.ResponseWriter, r *http.Request) {
	name, ok := streamName(w, r)
	if !ok {
		return
	}
	slog.Info("on_publish received", "name", name, "addr", r.PostFormValue("addr"))

	if err := h.manager.Start(name); err != nil {
		slog.Error("Failed to start stream", "name", name, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// onPublishDone is called by nginx-rtmp when a stream key disconnects.
func (h *Handler) onPublishDone(w http.ResponseWriter, r *http.Request) {
	name, ok := streamName(w, r)
	if !ok {
		return
	}
	slog.Info("on_publish_done received", "name", name)

	h.manager.Stop(name)
	w.WriteHeader(http.StatusOK)
}

// streamName extracts and validates the "name" field from the POST form.
// It writes the error response itself and returns false if the name is absent.
func streamName(w http.ResponseWriter, r *http.Request) (string, bool) {
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse form", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	name := r.PostFormValue("name")
	if name == "" {
		slog.Warn("Request missing stream name")
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	return name, true
}

// postOnly is middleware that rejects non-POST requests with 405.
func (h *Handler) postOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
