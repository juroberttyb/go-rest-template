package api

import (
	"net/http"

	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/service"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	c service.Order
}

func addOrderRoutes(root *gin.RouterGroup, c service.Order) {

	h := &orderHandler{
		c: c,
	}

	g := root.Group("orders")

	g.GET("", h.getOrders)
	g.GET(":order_id", h.get)
	g.POST("", h.new)
	// g.DELETE("", h.delete)
}

type getOrdersReq struct {
	pageReq
	Filter models.OrderStatus `form:"filter" validate:"optional"` // param to decide whether to return a list of recommended, attending, or attended kickstarts
	// Tag    []string            `form:"tag" validate:"optional"`
	AppID string `form:"app_id" validate:"optional"` // app id of source app where user comes from
}

// @Summary		Get a list of kickstarts
// @Description	Get a list of kickstarts where the parameter filter decides whether to return a list of general, attending, or attended kickstarts w.s.t to the given user
// @Tags			kickstart
// @Param			input	query	getOrdersReq	true	"pagination params and filter"
// @Produce		json
// @Success		200	{object}	pageResp{data=[]models.Order}
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/kickstarts [get]
// @Security		Bearer
func (h *orderHandler) getOrders(ctx *gin.Context) {
	p := getOrdersReq{}
	if err := ctx.BindQuery(&p); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	uid := ctx.GetString("user_id")
	if p.Count == 0 {
		p.Count = 10
	} else if p.Count > 50 {
		p.Count = 50
	}

	rctx := ctx.Request.Context()
	kickstarts, next, err := h.c.GetOrders(rctx, uid, p.Next, p.Count, p.Filter)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, pageResp{
		Data: kickstarts,
		Next: next,
	})
}

type getOrderReq struct {
	OrderID string `uri:"kickstart_id" binding:"required,uuid4"`
}

// @Summary		Get a kickstart
// @Description	Get a kickstart with its detailed info
// @Tags			kickstart
// @Param			kickstart_id	path	string	true	"ID of a kickstart"	example(8306778b-7287-f72b-6b26-a95316de96e4)"
// @Produce		json
// @Success		200	{object}	models.Order
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/kickstarts/{kickstart_id} [get]
// @Security		Bearer
func (h *orderHandler) get(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	p := getOrderReq{}
	if err := ctx.BindUri(&p); err != nil {
		handleError(ctx, err)
		return
	}

	rctx := ctx.Request.Context()
	kickstart, err := h.c.Get(rctx, userID, p.OrderID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, kickstart)
}

type attendReq struct {
	OrderID string  `json:"kickstart_id" binding:"required,uuid4"`
	Email   *string `json:"email"`
}

// @Summary		Attend a kickstart
// @Description	Attend a kickstart
// @Tags			kickstart
// @Param			jsonBody	body	attendReq	true	"kickstart id to attend and user's email"
// @Produce		json
// @Success		200
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/kickstarts/attendance [post]
// @Security		Bearer
func (h *orderHandler) new(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	cp := attendReq{}
	if err := ctx.BindJSON(&cp); err != nil {
		handleError(ctx, err)
		return
	}

	rctx := ctx.Request.Context()
	if err := h.c.New(rctx, userID, cp.OrderID, cp.Email); err != nil { // , p.Code
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}
