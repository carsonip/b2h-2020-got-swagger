package hack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatRoute(t *testing.T) {
	assert.Equal(t, "/api/{id}/s/{sid}", formatRoute("/api/:id/s/:sid"))
}
