package examples

import (
	"testing"

	dtomerge "github.com/cardinalby/go-dto-merge"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T {
	return &v
}

func TestSimpleExample(t *testing.T) {
	t.Parallel()

	type UserConfig struct {
		Role string // for non-pointer fields zero value indicates it's not specified
		Name string
	}

	type Config struct {
		Verbose *bool // it is a pointer to distinguish between "not specified" and false
		User    UserConfig
	}

	defaults := Config{
		Verbose: ptr(true),
		User: UserConfig{
			Role: "admin",
			Name: "John",
		},
	}
	userProvided := Config{
		User: UserConfig{
			Name: "Jane",
		},
	}
	res, err := dtomerge.Merge(defaults, userProvided)
	require.NoError(t, err)

	require.Equal(t, true, *res.Verbose)
	require.Equal(t, "admin", res.User.Role)
	require.Equal(t, "Jane", res.User.Name)
}
