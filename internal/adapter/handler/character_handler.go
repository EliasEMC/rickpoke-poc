package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
)

func RegisterCharacterRoutes(r *gin.RouterGroup, uc usecase.FetchCharacter) {
	r.GET("/rick/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be int"})
			return
		}
		res, err := uc.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}

type CharacterHandler struct {
	storeUC usecase.StoreCharacter
	listUC  usecase.ListCharacters
}

func NewCharacterHandler(storeUC usecase.StoreCharacter, listUC usecase.ListCharacters) *CharacterHandler {
	return &CharacterHandler{
		storeUC: storeUC,
		listUC:  listUC,
	}
}

type SaveCharacterRequest struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (h *CharacterHandler) Save(c *gin.Context) {
	var req SaveCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	character := &model.Rick{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := h.storeUC.Save(c.Request.Context(), character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, character)
}

func (h *CharacterHandler) List(c *gin.Context) {
	characters, err := h.listUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, characters)
}
