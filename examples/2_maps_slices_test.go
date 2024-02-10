package examples

import (
	"testing"

	dtomerge "github.com/cardinalby/go-dto-merge"
	"github.com/stretchr/testify/require"
)

func TestSlicesMapsExample(t *testing.T) {
	t.Parallel()

	type Config struct {
		Roles       []string
		Permissions map[string]bool
	}

	defaults := Config{
		Roles: []string{"admin", "user"},
		Permissions: map[string]bool{
			"read":  true,
			"write": false,
		},
	}

	userProvided := Config{
		Roles: []string{"user", "guest"},
		Permissions: map[string]bool{
			"write": true,
		},
	}

	res, err := dtomerge.Merge(defaults, userProvided,
		dtomerge.OptIterateMaps(true),
		dtomerge.OptMergeSlices(dtomerge.SlicesMergeStrategyUnique),
	)
	require.NoError(t, err)

	require.Equal(t, []string{"admin", "user", "guest"}, res.Roles)
	require.Equal(t, map[string]bool{
		"read":  true,
		"write": true,
	}, res.Permissions)
}
