package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case http.MethodGet:
		var prevID, size int64
		var err error

		q := r.URL.Query()
		if q.Get("prev_id") != "" {
			prevID, err = strconv.ParseInt(q.Get("prev_id"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		size = 5
		if q.Get("size") != "" {
			size, err = strconv.ParseInt(q.Get("size"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		todoReq := model.ReadTODORequest{
			PrevID: prevID,
			Size:   size,
		}

		todoRes, err := h.Read(r.Context(), &todoReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	case http.MethodPost:
		var todoReq model.CreateTODORequest
		json.NewDecoder(r.Body).Decode(&todoReq)

		if todoReq.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todoRes, err := h.Create(r.Context(), &todoReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	case http.MethodPut:
		var todoReq model.UpdateTODORequest
		json.NewDecoder(r.Body).Decode(&todoReq)

		if todoReq.Subject == "" || todoReq.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todoRes, err := h.Update(r.Context(), &todoReq)
		if err != nil {
			var verr model.ErrNotFound
			if errors.Is(err, &verr) {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	case http.MethodDelete:
		var todoReq model.DeleteTODORequest
		json.NewDecoder(r.Body).Decode(&todoReq)

		if len(todoReq.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todoRes, err := h.Delete(r.Context(), &todoReq)
		if err != nil {
			var verr model.ErrNotFound
			if errors.Is(err, &verr) {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if err := json.NewEncoder(w).Encode(todoRes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{
		TODO: todo,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{
		TODOs: todos,
	}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{
		TODO: todo,
	}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	if err := h.svc.DeleteTODO(ctx, req.IDs); err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
