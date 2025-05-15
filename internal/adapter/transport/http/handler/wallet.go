package handler

import (
	"java-code/internal/adapter/transport/http/handler/dto"
	"java-code/internal/service"
	"java-code/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletHandler struct {
	logger  logger.Logger
	service *service.Wallet
}

func NewWalletHandler(logger logger.Logger, service *service.Wallet) *WalletHandler {
	return &WalletHandler{
		logger:  logger,
		service: service,
	}
}

func (h WalletHandler) OperationWithWallet(c *gin.Context) {
	var req dto.WalletOperationRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Error("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errResponse{Error: "Invalid request body"})
		return
	}

	oparationWallet, err := dto.ToDomainWalletOperation(req)
	if err != nil {
		h.logger.Error("invalid walletId", zap.Error(err))
		c.JSON(http.StatusBadRequest, errResponse{Error: "Invalid wallet ID"})
		return
	}

	if err := h.service.OperationWithWallet(c.Request.Context(), oparationWallet); err != nil {
		h.logger.Error("Invalid register request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errResponse{Error: "Internal server error"})
		return
	}

	h.logger.Info("Operation wallet success")

	c.JSON(http.StatusOK, successResponse{Success: true})
}

func (h WalletHandler) GetWallets(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("failed to parse id", zap.Error(err))
		c.JSON(http.StatusBadRequest, errResponse{Error: "Invalid id"})
		return
	}

	walletBalance, err := h.service.GetWalletByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to update person", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errResponse{Error: "Internal server error"})
		return
	}
	h.logger.Info("wallet get success", zap.String("id wallet:", walletBalance.ID.String()))

	c.JSON(http.StatusOK, dto.ToDtoWalletBalanceRepsonse(walletBalance))
}
