package ports

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/wreckitral/pinjamean/internal/common/server/httperr"
	"github.com/wreckitral/pinjamean/internal/loans/app"
	"github.com/wreckitral/pinjamean/internal/loans/app/command"
	"github.com/wreckitral/pinjamean/internal/loans/app/query"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) SubmitLoan(w http.ResponseWriter, r *http.Request) {
	req := SubmitLoanRequest{}
	if err := render.Decode(r, &req); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	cmd := command.SubmitLoan{
		LoanUUID: uuid.New().String(),
		BorrowerUUID: req.BorrowerUuid.String(),
		LoanAmountIDR: req.AmountIdr,
		TermMonths: req.TermMonths,
		LoanType: loan.LoanType(req.LoanType),
	}
	if err := h.app.Commands.SubmitLoan.Handle(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.Header().Set("content-location", "/loans/"+cmd.LoanUUID)
	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetLoanByUuid(w http.ResponseWriter, r *http.Request, uuid openapi_types.UUID) {
	view, err := h.app.Queries.GetLoanByID.Handle(r.Context(), query.GetLoan{
		LoanUUID: uuid.String(),
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, view)
}

func appLoanToResponse(appLoan query.LoanView) (LoanResponse, error) {
	loanUuid, err := uuid.Parse(appLoan.UUID)
	if err != nil {
		return LoanResponse{}, fmt.Errorf("invalid loan UUID %q: %w", appLoan.UUID, err)
	}

	borrowerUuid, err := uuid.Parse(appLoan.BorrowerUUID)
	if err != nil {
		return LoanResponse{}, fmt.Errorf("invalid borrower UUID %q: %w", appLoan.BorrowerUUID, err)
	}

	apiType, err := domainLoanTypeToAPI(appLoan.LoanType)
	if err != nil {
		return LoanResponse{}, err
	}

	return LoanResponse{
		Uuid:            loanUuid,
		BorrowerUuid:    borrowerUuid,
		AmountIdr:       appLoan.LoanAmountIDR,
		TermMonths:      appLoan.TermMonths,
		LoanType:        apiType,
		InterestRateApr: float32(appLoan.InterestRateAPR),
		Status:          string(appLoan.Status),
	}, nil
}

// apiLoanTypeToDomain converts an incoming request's LoanType into the
// domain's own type, before it's handed to loan.NewLoan.
func apiLoanTypeToDomain(t LoanType) (loan.LoanType, error) {
	switch t {
	case KPR:
		return loan.TypeKPR, nil
	case KKB:
		return loan.TypeKKB, nil
	case KTA:
		return loan.TypeKTA, nil
	case KMG:
		return loan.TypeKMG, nil
	case KMK:
		return loan.TypeKMK, nil
	case KUR:
		return loan.TypeKUR, nil
	default:
		return "", fmt.Errorf("no domain mapping for API loan type %q", t)
	}
}

// domainLoanTypeToAPI converts the domain's LoanType into the wire-format
// enum generated from the OpenAPI spec.
func domainLoanTypeToAPI(t loan.LoanType) (LoanType, error) {
	switch t {
	case loan.TypeKPR:
		return KPR, nil
	case loan.TypeKKB:
		return KKB, nil
	case loan.TypeKTA:
		return KTA, nil
	case loan.TypeKMG:
		return KMG, nil
	case loan.TypeKMK:
		return KMK, nil
	case loan.TypeKUR:
		return KUR, nil
	default:
		return "", fmt.Errorf("no API mapping for domain loan type %q", t)
	}
}

