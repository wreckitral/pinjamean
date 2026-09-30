package ports

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/wreckitral/pinjamean/internal/common/server/httperr"
	"github.com/wreckitral/pinjamean/internal/identity/app"
	"github.com/wreckitral/pinjamean/internal/identity/app/command"
	"github.com/wreckitral/pinjamean/internal/identity/app/query"
	"github.com/wreckitral/pinjamean/internal/identity/domain/account"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app: app}
}

var _ ServerInterface = HttpServer{}

// LoginHandler adapts the application login query to POST /auth/login.
func LoginHandler(login query.LoginHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		view, err := login.Handle(r.Context(), query.Login{
			Email:    string(req.Email),
			Password: req.Password,
		})

		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(LoginResponse{
			AccessToken: view.Token,
			TokenType:   "Bearer",
		})
	})
}

func (h HttpServer) SaveAccount(w http.ResponseWriter, r *http.Request) {
	var req SaveAccountRequest
	if err := render.Decode(r, &req); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}
	if req.Password == nil || *req.Password == "" {
		httperr.BadRequest("missing-password", fmt.Errorf("password is required"), w, r)
		return
	}

	role, err := apiRoleToDomain(req.Role)
	if err != nil {
		httperr.BadRequest("invalid-account-role", err, w, r)
		return
	}

	cmd := command.SaveAccount{
		UUID:     req.Uuid.String(),
		Email:    string(req.Email),
		Password: *req.Password,
		Role:     role,
	}
	if err := h.app.Commands.SaveAccount.Handle(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.Header().Set("content-location", "/accounts/"+cmd.UUID)
	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetAccountByEmail(w http.ResponseWriter, r *http.Request, email openapi_types.Email) {
	view, err := h.app.Queries.GetAccountByEmail.Handle(r.Context(), query.GetAccountByEmail{Email: string(email)})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	resp, err := accountViewToResponse(view)
	if err != nil {
		httperr.InternalError("invalid-account", err, w, r)
		return
	}
	render.Respond(w, r, resp)
}

func accountViewToResponse(view query.AccountView) (AccountResponse, error) {
	parsedUUID, err := uuid.Parse(view.UUID)
	if err != nil {
		return AccountResponse{}, fmt.Errorf("invalid account UUID %q: %w", view.UUID, err)
	}
	role, err := domainRoleToAPI(view.Role)
	if err != nil {
		return AccountResponse{}, err
	}
	return AccountResponse{Uuid: openapi_types.UUID(parsedUUID), Email: openapi_types.Email(view.Email), Role: role}, nil
}

func apiRoleToDomain(role AccountRole) (account.Role, error) {
	switch role {
	case Officer:
		return account.RoleOfficer, nil
	case CreditAnalyst:
		return account.RoleCreditAnalyst, nil
	default:
		return "", fmt.Errorf("unsupported account role %q", role)
	}
}

func domainRoleToAPI(role account.Role) (AccountRole, error) {
	switch role {
	case account.RoleOfficer:
		return Officer, nil
	case account.RoleCreditAnalyst:
		return CreditAnalyst, nil
	default:
		return "", fmt.Errorf("unsupported account role %q", role)
	}
}
