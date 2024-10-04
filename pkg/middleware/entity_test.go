package middleware

import (
	"fmt"
	"github.com/mikestefanello/pagoda/pkg/db/sqlc"
	"testing"

	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUser(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	ctx.SetParamNames("user")
	ctx.SetParamValues(fmt.Sprintf("%d", usr.ID))
	_ = tests.ExecuteMiddleware(ctx, LoadUser(c.Queries))
	ctxUsr, ok := ctx.Get(context.UserKey).(*sqlc.User)
	require.True(t, ok)
	assert.Equal(t, usr.ID, ctxUsr.ID)
}
