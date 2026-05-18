package httpapi

import (
	"net/http"
	"strconv"

	authusecase "kerjadekat/backend/internal/auth/usecase"
	agenthttp "kerjadekat/backend/internal/agent/delivery/http"
	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/internal/httpapi/middleware"
	"kerjadekat/backend/internal/order/delivery/ws"
	orderusecase "kerjadekat/backend/internal/order/usecase"
	skillusecase "kerjadekat/backend/internal/skill/usecase"
	userusecase "kerjadekat/backend/internal/user/usecase"
	workerusecase "kerjadekat/backend/internal/worker/usecase"
	"kerjadekat/backend/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Deps wires HTTP delivery to application use cases.
type Deps struct {
	Auth    *authusecase.Auth
	Users   *userusecase.Users
	Skills  *skillusecase.Skills
	Orders  *orderusecase.Orders
	Workers *workerusecase.Workers
	Agents  *agenthttp.Handler
	Tokens  *token.Issuer
	WSHub   *ws.Hub
}

func Mount(r *gin.Engine, d Deps) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	if d.WSHub != nil {
		v1.GET("/ws", d.WSHub.HandleWS)
	}

	v1.POST("/auth/otp/request", func(c *gin.Context) {
		var body struct {
			PhoneNumber string `json:"phone_number" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		if err := d.Auth.RequestOTP(c.Request.Context(), body.PhoneNumber); err != nil {
			WriteError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})

	v1.POST("/auth/otp/verify", func(c *gin.Context) {
		var body struct {
			PhoneNumber string `json:"phone_number" binding:"required"`
			Code        string `json:"code" binding:"required"`
			Role        string `json:"role"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		tokens, err := d.Auth.VerifyOTP(c.Request.Context(), body.PhoneNumber, body.Code, body.Role)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, tokens)
	})

	v1.POST("/auth/social", func(c *gin.Context) {
		var body struct {
			Provider string `json:"provider" binding:"required"`
			Subject  string `json:"subject" binding:"required"`
			Email    string `json:"email"`
			Name     string `json:"name"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		tokens, err := d.Auth.SocialLogin(c.Request.Context(), authusecase.SocialLoginInput{
			Provider: body.Provider,
			Subject:  body.Subject,
			Email:    body.Email,
			Name:     body.Name,
			Role:     body.Role,
		})
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, tokens)
	})

	v1.POST("/auth/refresh", func(c *gin.Context) {
		var body struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		tokens, err := d.Auth.Refresh(c.Request.Context(), body.RefreshToken)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, tokens)
	})

	authed := v1.Group("", middleware.JWTAuth(d.Tokens))

	if d.Agents != nil {
		agentRoutes := authed.Group("/agent", middleware.RequireRoles(domain.RoleAgent, domain.RoleAdmin))
		agentRoutes.GET("/territories", d.Agents.ListTerritories)
		agentRoutes.POST("/workers", d.Agents.RegisterWorker)
	}

	authed.GET("/me", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		u, err := d.Users.Me(c.Request.Context(), cl.UserID)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, u)
	})

	authed.GET("/skill-categories", func(c *gin.Context) {
		rows, err := d.Skills.ListCategories(c.Request.Context())
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": rows})
	})

	authed.POST("/orders", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if cl.Role != domain.RoleConsumer {
			WriteError(c, domain.ErrForbidden)
			return
		}
		var body struct {
			SkillID           int     `json:"skill_id" binding:"required"`
			Description       *string `json:"description"`
			Latitude          float64 `json:"latitude" binding:"required"`
			Longitude         float64 `json:"longitude" binding:"required"`
			ConsumerAddress   *string `json:"consumer_address"`
			PaymentMethodFee  *string `json:"payment_method_fee"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		ord, err := d.Orders.Create(c.Request.Context(), orderusecase.CreateOrderInput{
			ConsumerID:       cl.UserID,
			SkillID:          body.SkillID,
			Description:      body.Description,
			Latitude:         body.Latitude,
			Longitude:        body.Longitude,
			ConsumerAddress:  body.ConsumerAddress,
			PaymentMethodFee: body.PaymentMethodFee,
		})
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusCreated, ord)
	})

	authed.GET("/orders", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		rows, err := d.Orders.ListMine(c.Request.Context(), cl.UserID, cl.Role, limit, offset)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": rows})
	})

	authed.GET("/orders/:id", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		ord, err := d.Orders.Get(c.Request.Context(), id, cl.UserID, cl.Role)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, ord)
	})

	authed.POST("/orders/:id/accept", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if cl.Role != domain.RoleWorker {
			WriteError(c, domain.ErrForbidden)
			return
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var body struct {
			AgreedRate *float64 `json:"agreed_rate"`
		}
		_ = c.ShouldBindJSON(&body)
		ord, err := d.Orders.Accept(c.Request.Context(), id, orderusecase.AcceptOrderInput{
			WorkerUserID: cl.UserID,
			AgreedRate:   body.AgreedRate,
		})
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, ord)
	})

	authed.POST("/orders/:id/start", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if cl.Role != domain.RoleWorker {
			WriteError(c, domain.ErrForbidden)
			return
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		ord, err := d.Orders.Start(c.Request.Context(), id, cl.UserID)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, ord)
	})

	authed.POST("/orders/:id/complete", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		ord, err := d.Orders.Complete(c.Request.Context(), id, cl.UserID)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, ord)
	})

	authed.POST("/orders/:id/cancel", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var body struct {
			Reason *string `json:"reason"`
		}
		_ = c.ShouldBindJSON(&body)
		ord, err := d.Orders.Cancel(c.Request.Context(), id, cl.UserID, body.Reason)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, ord)
	})

	authed.GET("/workers/nearby", func(c *gin.Context) {
		lat, err := strconv.ParseFloat(c.Query("lat"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lat required"})
			return
		}
		lng, err := strconv.ParseFloat(c.Query("lng"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lng required"})
			return
		}
		var radius float64 = 5000
		if r := c.Query("radius"); r != "" {
			radius, err = strconv.ParseFloat(r, 64)
			if err != nil || radius <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid radius"})
				return
			}
		}
		var skillID *int
		if s := c.Query("skill"); s != "" {
			id, err := strconv.Atoi(s)
			if err != nil || id <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill"})
				return
			}
			skillID = &id
		}
		items, err := d.Workers.Nearby(c.Request.Context(), workerusecase.NearbyInput{
			Latitude:     lat,
			Longitude:    lng,
			RadiusMeters: radius,
			SkillID:      skillID,
		})
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	authed.GET("/workers/me", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if cl.Role != domain.RoleWorker {
			WriteError(c, domain.ErrForbidden)
			return
		}
		p, err := d.Workers.Me(c.Request.Context(), cl.UserID)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, p)
	})

	authed.PATCH("/workers/me/availability", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if cl.Role != domain.RoleWorker {
			WriteError(c, domain.ErrForbidden)
			return
		}
		var body struct {
			Availability string `json:"availability" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		if err := d.Workers.SetAvailability(c.Request.Context(), cl.UserID, body.Availability); err != nil {
			WriteError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})

	authed.PATCH("/workers/me/location", func(c *gin.Context) {
		cl, ok := middleware.Claims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if cl.Role != domain.RoleWorker {
			WriteError(c, domain.ErrForbidden)
			return
		}
		var body struct {
			Latitude  float64 `json:"latitude" binding:"required"`
			Longitude float64 `json:"longitude" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		if err := d.Workers.SetLocation(c.Request.Context(), cl.UserID, body.Latitude, body.Longitude); err != nil {
			WriteError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})
}
