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

	root.GET("board", h.getBoard)

	// FIXME: need to consider order status modification under heavy concurrent access (taking an order already removed, two takes at the same order which exceeds order provided amount...)
	g := root.Group("orders")
	// FIXME: need to do pagination and filter, pagination should start from latest taker price and grow up and down
	// FIXME: implement this get method
	// g.GET(":order_id", h.get)
	g.POST("", h.make)
	g.PATCH("", h.take)
	g.DELETE(":order_id", h.delete)
}

type getOrdersReq struct {
	// FIXME: add pagination with var below
	// pageReq
	// if true, return orders created by the user
	BoardType models.OrderBoardType `form:"board_type" validate:"optional"`
}

// @Summary		Get a order board
// @Description	Get a order board
// @Tags			order
// @Param			input	query	getOrdersReq	true	"related parameters"
// @Produce		json
// @Success		200	{object}	pageResp{data=[]models.Order}
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/orders [get]
// @Security		Bearer
func (h *orderHandler) getBoard(ctx *gin.Context) {
	p := getOrdersReq{}
	if err := ctx.BindQuery(&p); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	board, next, err := h.c.GetBoard(
		ctx.Request.Context(),
		p.BoardType,
	)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, pageResp{
		Data: board,
		Next: next,
	})
}

// FIXME: need to consider integer overflow here, for example price*amount > int max value
// FIXME: should use fixed type in64 or int32 instead of int to avoid overflow
type makeOrderBody struct {
	Action models.OrderAction `form:"action" binding:"required" example:"buy"`
	Price  int                `json:"price" binding:"required,min=1" example:"10"`
	Amount int                `json:"amount" binding:"required,min=1" example:"100"`
}

// @Summary		Make a order
// @Description	Make a order
// @Tags			order
// @Param			jsonBody	body	makeOrderBody	true	"order id to attend and user's email"
// @Produce		json
// @Success		200
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/orders [post]
// @Security		Bearer
func (h *orderHandler) make(ctx *gin.Context) {
	b := makeOrderBody{}
	if err := ctx.BindJSON(&b); err != nil {
		handleError(ctx, err)
		return
	}

	if err := h.c.Make(
		ctx.Request.Context(),
		b.Action,
		b.Price,
		b.Amount,
	); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}

type takeOrderBody struct {
	// FIXME:user_id should be retrieved from the user's jwt token
	UserID string             `form:"user_id" binding:"required,uuid" example:"uuid"`
	Action models.OrderAction `form:"action" binding:"required" example:"buy"`
	Amount int                `json:"amount" binding:"required,min=1" example:"100"`
}

// @Summary		Take a order
// @Description	Take a order
// @Tags			order
// @Param			jsonBody	body	takeOrderBody	true	"order id to attend and user's email"
// @Produce		json
// @Success		200
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/orders [patch]
// @Security		Bearer
func (h *orderHandler) take(ctx *gin.Context) {
	b := takeOrderBody{}
	if err := ctx.BindJSON(&b); err != nil {
		handleError(ctx, err)
		return
	}

	if err := h.c.Take(
		ctx.Request.Context(),
		b.Action,
		b.Amount,
		b.UserID,
	); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}

type deleteOrderUri struct {
	OrderID string `uri:"order_id" binding:"required,uuid4"`
}

// @Summary		Delete a order
// @Description	Delete a order
// @Tags			order
// @Param		order_id	path	string	true	"ID of order"
// @Produce		json
// @Success		200
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/orders [delete]
// @Security		Bearer
func (h *orderHandler) delete(ctx *gin.Context) {
	// FIXME: a order must only be deleted by the creator

	u := deleteOrderUri{}
	if err := ctx.BindUri(&u); err != nil {
		handleError(ctx, err)
		return
	}

	if err := h.c.Delete(
		ctx.Request.Context(),
		u.OrderID,
	); err != nil { // , p.Code
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}
