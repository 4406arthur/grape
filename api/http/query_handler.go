package http

import (
	"net/http"
	"strconv"

	"github.com/4406arthur/grape/api"
	"github.com/4406arthur/grape/repository"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type QueryHandler struct {
	QUsecase api.QueryUsecase
}

type ResponseError struct {
	Message string `json:"message"`
}

type Blocks struct {
	Blocks []repository.Block `json:"blocks"`
}
type BlockWithTXs struct {
	BlockNum   int64    `json:"block_num"`
	BlockHash  string   `json:"block_hash"`
	ParentHash string   `json:"parent_hash"`
	BlockTime  uint64   `json:"block_time"`
	TxList     []string `json:"transactions"`
}

func NewQueryHandler(e *echo.Echo, q api.QueryUsecase) {
	handler := &QueryHandler{
		QUsecase: q,
	}
	e.GET("/blocks", handler.GetLatestNBlock)
	e.GET("/blocks/:id", handler.GetBlockWithTXsByID)
	e.GET("/transaction/:txHash", handler.GetTXByHash)

}

func (q *QueryHandler) GetLatestNBlock(c echo.Context) error {
	limit := c.QueryParam("limit")
	num, _ := strconv.Atoi(limit)
	ctx := c.Request().Context()

	blocks, err := q.QUsecase.GetLatestNBlock(ctx, num)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	// log.Debugf("+%v", blocks)
	return c.JSON(http.StatusOK, &Blocks{Blocks: blocks})
}

func (q *QueryHandler) GetBlockWithTXsByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()

	block, txs, err := q.QUsecase.GetFullBlockByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	// log.Debugf("+%v", block)
	return c.JSON(http.StatusOK,
		&BlockWithTXs{
			BlockNum:   block.BlockNum,
			BlockHash:  block.BlockHash,
			ParentHash: block.ParentHash,
			BlockTime:  block.BlockTime,
			TxList:     fetchTXHash(txs),
		})
}

func (q *QueryHandler) GetTXByHash(c echo.Context) error {
	hash := c.Param("txHash")
	log.Debugf("the input: %s", hash)
	ctx := c.Request().Context()

	tx, err := q.QUsecase.GetTXByHash(ctx, hash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	// log.Debugf("+%v", blocks)
	return c.JSON(http.StatusOK, tx)
}

func fetchTXHash(txs []repository.Tx) []string {
	res := make([]string, 0)
	for _, v := range txs {
		res = append(res, v.TxHash)
	}
	return res
}
