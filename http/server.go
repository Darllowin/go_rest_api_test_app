package http

import (
	"errors"
	"net/http"
)

type HTTPserver struct {
	httpHendlers *HTTPHandlers
}

func NewHTTPServer(httpHandler *HTTPHandlers) *HTTPserver {
	return &HTTPserver{
		httpHendlers: httpHandler,
	}
}

func (s *HTTPserver) StartServer() error {
	router := http.NewServeMux()

	router.HandleFunc("POST /tasks", s.httpHendlers.HandleCreateTask)
	router.HandleFunc("GET /tasks/{title}", s.httpHendlers.HandleGetTask)

	router.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("completed") == "false" {
			s.httpHendlers.HandleGetAllUncompletedTasks(w, r)
			return
		}
		s.httpHendlers.HandleGetAllTasks(w, r)
	})

	router.HandleFunc("PATCH /tasks/{title}", s.httpHendlers.HandleCompletedTask)
	router.HandleFunc("DELETE /tasks/{title}", s.httpHendlers.HandleDeleteTask)

	server := &http.Server{
		Addr:    ":9091",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}
