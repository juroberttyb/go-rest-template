package api

import (
	"net/http"

	"github.com/A-pen-app/kickstart/api/middleware"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/service"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	c    service.Order
	auth service.Auth
}

func addOrderRoutes(root *gin.RouterGroup, c service.Order, auth service.Auth) {
	h := &orderHandler{
		c:    c,
		auth: auth,
	}

	root.GET("board", h.getBoard)

	// FIXME: temporary token generator for testing
	root.GET("token", h.getToken)

	g := root.Group("orders")
	g.Use(middleware.AuthUser(auth))
	g.Use(middleware.NeedPermission(models.Audience))

	// FIXME: need to do pagination and filter, pagination should start from latest taker price and grow up and down
	// FIXME: implement this get method
	// g.GET(":order_id", h.get)
	g.POST("make", h.make)
	g.PATCH("take", h.take)
	g.DELETE(":order_id", h.delete)
}

type GetTokenReq struct {
	UserID string `form:"user_id" binding:"required,uuid4"`
}

type GetTokenResp struct {
	Token string `json:"token"`
}

// @Summary		temporary token generator for testing
// @Description	temporary token generator for testing which return a JWT once verified.
// @Tags		order
// @Param		user_id	path	string	true	"ID of user to get token of"
// @Produce		json
// @Success		200	{object} 	GetTokenResp
// @Failure		400	{object}	errorResp "cannot generate token"
// @Failure		500	{object}	errorResp
// @Router			/token [get]
func (h *orderHandler) getToken(ctx *gin.Context) {
	p := GetTokenReq{}
	if err := ctx.BindQuery(&p); err != nil {
		handleError(ctx, err)
		return
	}

	rctx := ctx.Request.Context()
	token, err := h.auth.IssueToken(rctx, p.UserID, models.Audience)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &GetTokenResp{
		Token: token,
	})
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
// @Router			/board [get]
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

// FIXME: need to consider integer overflow here, for example price*quantity > int max value
// FIXME: should use fixed type in64 or int32 instead of int to avoid overflow
type makeOrderBody struct {
	Action   models.OrderAction `json:"action" binding:"required" example:"buy"`
	Price    int                `json:"price" binding:"required,min=1" example:"10"`
	Quantity int                `json:"quantity" binding:"required,min=1" example:"100"`
}

// @Summary		Make a order
// @Description	Make a order
// @Tags			order
// @Param			jsonBody	body	makeOrderBody	true	"order id to attend and user's email"
// @Produce		json
// @Success		200
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/orders/make [post]
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
		b.Quantity,
	); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}

type takeOrderBody struct {
	// FIXME:user_id should be retrieved from the user's jwt token
	Action   models.OrderAction `json:"action" binding:"required" example:"buy"`
	Quantity int                `json:"quantity" binding:"required,min=1" example:"100"`
}

// @Summary		Take a order
// @Description	Take a order
// @Tags			order
// @Param			jsonBody	body	takeOrderBody	true	"order id to attend and user's email"
// @Produce		json
// @Success		200
// @Failure		400	{object}	errorResp
// @Failure		500	{object}	errorResp
// @Router			/orders/take [patch]
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
		b.Quantity,
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
