package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/hiiamanop/simple_bank/db/sqlc"
	"github.com/hiiamanop/simple_bank/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username         string    `json:"username"`
	Fullname         string    `json:"fullname"`
	Email            string    `json:"email"`
	PaswordChangedAt time.Time `json:"pasword_changed_at"`
	CreatedAt        time.Time `json:"created_at"`
}

type GetUserResponse struct {
	Username         string    `json:"username"`
	Fullname         string    `json:"fullname"`
	Email            string    `json:"email"`
	PaswordChangedAt time.Time `json:"pasword_changed_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Fullname:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username:         user.Username,
		Fullname:         user.Fullname,
		Email:            user.Email,
		PaswordChangedAt: user.PaswordChangedAt,
		CreatedAt:        user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUser(ctx *gin.Context) {
	var req struct {
		Username string `uri:"username" binding:"required,alphanum,min=1"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := GetUserResponse{
		Username:         user.Username,
		Fullname:         user.Fullname,
		Email:            user.Email,
		PaswordChangedAt: user.PaswordChangedAt,
		CreatedAt:        user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
