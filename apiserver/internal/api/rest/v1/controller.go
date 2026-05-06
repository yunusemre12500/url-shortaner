package v1

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.yunus-emre.dev/url-shortaner/model"
	apierrors "go.yunus-emre.dev/url-shortaner/pkg/api/errors"
	"go.yunus-emre.dev/url-shortaner/pkg/util/httputil"
	"go.yunus-emre.dev/url-shortaner/pkg/util/slug"
	"go.yunus-emre.dev/url-shortaner/storage"
)

type Controller struct {
	storage storage.Storage
	tracer  trace.Tracer
}

func New(storage storage.Storage) *Controller {
	return &Controller{
		storage: storage,
		tracer:  otel.Tracer("api-server/rest/v1"),
	}
}

func (c *Controller) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /v1/links", otelhttp.NewHandler(http.HandlerFunc(c.createLink), "POST /v1/links"))
	mux.Handle("GET /v1/links/{slug}", otelhttp.NewHandler(http.HandlerFunc(c.getLink), "GET /v1/links/{slug}"))
}

func (c *Controller) createLink(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "create link")

	defer span.End()

	var body *CreateLinkRequestBodyParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.RespondWithJSON(w, http.StatusBadRequest, apierrors.InvalidJSONBody)

		return
	}

	if errs := ValidateCreateLinkRequestBodyParams(body); errs != nil {
		httputil.RespondWithJSON(w, http.StatusBadRequest, errs.Error())

		return
	}

	createLinkParams := model.CreateLinkParams{
		ExpiresAt:   body.ExpiresAt,
		OriginalURL: body.OriginalURL,
		Slug:        body.Slug,
	}

	if body.Slug == "" {
		slug, err := slug.Generate(8)

		if err != nil {
			span.SetStatus(codes.Error, "failed to generate slug")
			span.RecordError(err)

			httputil.RespondWithJSON(w, http.StatusInternalServerError, apierrors.InternalServerError)

			return
		}

		createLinkParams.Slug = slug
	}

	link := model.CreateLink(createLinkParams)

	if err := c.storage.CreateLink(ctx, link); err != nil {
		span.SetStatus(codes.Error, "failed to create link")
		span.RecordError(err)

		if err == storage.ErrConflict {
			httputil.RespondWithJSON(w, http.StatusConflict, apierrors.Conflict)

			return
		}

		httputil.RespondWithJSON(w, http.StatusInternalServerError, apierrors.InternalServerError)

		return
	}

	response := &CreateLinkResponseBodyParams{
		ClickCount:  link.ClickCount,
		CreatedAt:   link.CreatedAt,
		ExpiresAt:   link.ExpiresAt,
		ID:          link.ID,
		OriginalURL: link.OriginalURL,
		Slug:        link.Slug,
	}

	span.SetStatus(codes.Ok, "link created")

	httputil.RespondWithJSON(w, http.StatusCreated, response)
}

func (c *Controller) getLink(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "get link")

	defer span.End()

	link, err := c.storage.GetLinkBySlug(ctx, r.PathValue("slug"))

	if err != nil {
		span.SetStatus(codes.Error, "failed to get link from storage")
		span.RecordError(err)

		httputil.RespondWithJSON(w, http.StatusInternalServerError, apierrors.InternalServerError)

		return
	}

	if link == nil {
		span.SetStatus(codes.Error, "link not found")

		httputil.RespondWithJSON(w, http.StatusNotFound, apierrors.NotFound)

		return
	}

	if link.Expired() {
		span.SetStatus(codes.Error, "link expired")

		httputil.RespondWithJSON(w, http.StatusGone, apierrors.Expired)

		return
	}

	if err = c.storage.IncrementClickCountByLinkID(ctx, link.ID); err != nil {
		span.SetStatus(codes.Error, "failed to increment click count")
		span.RecordError(err)

		httputil.RespondWithJSON(w, http.StatusInternalServerError, apierrors.InternalServerError)

		return
	}

	span.SetStatus(codes.Ok, "link retrieved")

	http.Redirect(w, r, link.OriginalURL, http.StatusMovedPermanently)
}
