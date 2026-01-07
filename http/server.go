package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	router.Path("/tasks").Methods("POST").HandlerFunc(s.httpHendlers.HandleCreateTask)
	router.Path("/tasks/{title}").Methods("GET").HandlerFunc(s.httpHendlers.HandleGetTask)
	router.Path("/tasks").Methods("GET").Queries("completed", "true").HandlerFunc(s.httpHendlers.HandleGetAllUncompletedTasks)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.httpHendlers.HandleGetAllTasks)
	router.Path("/tasks/{title}").Methods("PATCH").HandlerFunc(s.httpHendlers.HandleCompletedTask)
	router.Path("/tasks/{title}").Methods("DELETE").HandlerFunc(s.httpHendlers.HandleDeleteTask)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
