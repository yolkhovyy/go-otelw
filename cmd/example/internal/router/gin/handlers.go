package gin

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (c *Controller) echo(gctx *gin.Context) {
	const defaultCount int = 5
	count := defaultCount

	scount := gctx.Query("count")
	if scount != "" {
		cnt, err := strconv.Atoi(scount)
		if err != nil {
			gctx.String(http.StatusBadRequest, "invalid count query parameter: %s", scount)

			return
		}

		count = cnt
	}

	body, err := io.ReadAll(gctx.Request.Body)
	if err != nil {
		gctx.String(http.StatusBadRequest, "read body: %w", err)

		return
	}

	input := string(body)
	if input == "" {
		gctx.String(http.StatusBadRequest, "empty request body")

		return
	}

	slog.InfoContext(
		gctx.Request.Context(),
		"echo",
		slog.String("input", input),
	)

	output, err := c.domain.Echo(gctx.Request.Context(), input, count)
	if err != nil {
		gctx.String(http.StatusInternalServerError, "echo: %w", err)

		return
	}

	gctx.String(http.StatusOK, "%s", output)
}
