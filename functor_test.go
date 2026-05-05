package kitty_test

import (
	"testing"

	"github.com/peng225/kitty"
	"github.com/stretchr/testify/require"
)

func TestIdentityFunctor(t *testing.T) {
	objects := []kitty.Object{"a", "b"}
	morphisms := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c1, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
	c2 := c1

	objMap := map[kitty.Object]kitty.Object{
		"a": "a",
		"b": "b",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "f",
	}
	_, err = kitty.NewFunctor(c1, c2, objMap, mMap)
	require.NoError(t, err)
}

func TestSameFormFunctor(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
	}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c1, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "g",
			Source:      "p",
			Destination: "q",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c2, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "p",
		"b": "q",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "g",
	}
	_, err = kitty.NewFunctor(c1, c2, objMap, mMap)
	require.NoError(t, err)
}

func TestDifferentFormFunctor(t *testing.T) {
	objects1 := []kitty.Object{"a", "b", "c"}
	morphisms1 := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
		{
			ID:          "g",
			Source:      "b",
			Destination: "c",
		},
	}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c1, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "h",
			Source:      "p",
			Destination: "q",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c2, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "p",
		"b": "q",
		"c": "q",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "h",
		"g": kitty.Identity,
	}
	_, err = kitty.NewFunctor(c1, c2, objMap, mMap)
	require.NoError(t, err)
}

func TestLongMorphismChain(t *testing.T) {
	objects := []kitty.Object{"a", "b", "c", "d"}
	morphisms := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
		{
			ID:          "g",
			Source:      "b",
			Destination: "c",
		},
		{
			ID:          "h",
			Source:      "c",
			Destination: "d",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "f",
		"g": "g",
		"h": "h",
	}
	_, err = kitty.NewFunctor(c, c, objMap, mMap)
	require.NoError(t, err)
}
