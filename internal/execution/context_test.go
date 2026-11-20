package execution

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.Background())
	assert.NotNil(t, ctx)
}
