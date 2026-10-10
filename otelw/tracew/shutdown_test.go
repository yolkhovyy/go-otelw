package tracew

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShutdown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  Config
		writers []io.Writer
	}{
		{name: "disabled", config: Config{}},
		{name: "stdout", config: Config{Enable: true}, writers: []io.Writer{io.Discard}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			instance, err := Configure(ctx, test.config, nil, test.writers...)
			require.NoError(t, err)

			require.NoError(t, instance.Shutdown(ctx))
		})
	}
}
