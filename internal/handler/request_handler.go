package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"rate-limited/internal/model"
	"rate-limited/internal/service"
)

// RequestHandler handles HTTP requests for the API.
type RequestHandler struct {
	rateLimiterService *service.RateLimiterService
}

// NewRequestHandler creates a new RequestHandler.
func NewRequestHandler(rateLimiterService *service.RateLimiterService) *RequestHandler {
	return &RequestHandler{
		rateLimiterService: rateLimiterService,
	}
}

// HandleRequest handles for requests to POST /request. It checks the rate limit and responds accordingly.
func (h *RequestHandler) HandleRequest(c *fiber.Ctx) error {
	var req model.RequestPayload

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid request body",
		})
	}

	req.UserID = strings.TrimSpace(req.UserID)
	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "user_id is required",
		})
	}

	allowed := h.rateLimiterService.ProcessRequest(req.UserID)
	if !allowed {
		return c.Status(fiber.StatusTooManyRequests).JSON(model.ErrorResponse{
			Error: "rate limit exceeded: max 5 requests per user per minute",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.RequestResponse{
		Message: "request accepted",
		UserID:  req.UserID,
	})
}

// GetStats handles requests to GET /stats. It returns paginated stats about user requests.
func (h *RequestHandler) GetStats(c *fiber.Ctx) error {
	page := 1
	limit := 10

	if pageQuery := c.Query("page"); pageQuery != "" {
		parsedPage, err := strconv.Atoi(pageQuery)
		if err != nil || parsedPage <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error: "page must be a positive integer",
			})
		}
		page = parsedPage
	}

	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil || parsedLimit <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error: "limit must be a positive integer",
			})
		}
		limit = parsedLimit
	}

	stats := h.rateLimiterService.GetStats(page, limit)
	return c.Status(fiber.StatusOK).JSON(stats)
}
