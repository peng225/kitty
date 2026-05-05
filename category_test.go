package kitty_test

import (
	"testing"

	"github.com/peng225/kitty"
	"github.com/stretchr/testify/require"
)

func TestOneObjectCategory(t *testing.T) {
	objects := []kitty.Object{"a"}
	morphisms := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "a",
		},
	}
	// 合成を指示するときに実装の詳細である"id_"が出てくるのは辛い。
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"f", "f"}: kitty.Identity,
	}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
}

func TestTwoObjectCategory(t *testing.T) {
	objects := []kitty.Object{"a", "b"}
	morphisms := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
}
