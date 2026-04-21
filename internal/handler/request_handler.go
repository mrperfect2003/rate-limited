package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"rate-limited/internal/model"
	"rate-limited/internal/service"
)

type RequestHandler struct {
	rateLimiterService *service.RateLimiterService
}

func NewRequestHandler(rateLimiterService *service.RateLimiterService) *RequestHandler {
	return &RequestHandler{
		rateLimiterService: rateLimiterService,
	}
}

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

	accepted, job, err := h.rateLimiterService.HandleIncomingRequest(req.UserID, req.Payload)
	if err != nil {
		return c.Status(fiber.StatusTooManyRequests).JSON(model.ErrorResponse{
			Error: err.Error(),
		})
	}

	if accepted {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "request accepted",
			"user_id": req.UserID,
			"mode":    "direct",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "rate limit exceeded, request queued for retry",
		"user_id": req.UserID,
		"mode":    "queued",
		"job_id":  job.ID,
		"status":  job.Status,
	})
}

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

func (h *RequestHandler) GetJob(c *fiber.Ctx) error {
	jobID := strings.TrimSpace(c.Params("id"))
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "job id is required",
		})
	}

	job, exists := h.rateLimiterService.GetJob(jobID)
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error: "job not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(job)
}

func (h *RequestHandler) GetQueueStats(c *fiber.Ctx) error {
	stats := h.rateLimiterService.GetQueueStats()
	return c.Status(fiber.StatusOK).JSON(stats)
}

func (h *RequestHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "service is healthy",
	})
}
