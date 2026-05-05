package kitty_test

import (
	"testing"

	"github.com/peng225/kitty"
	"github.com/stretchr/testify/require"
)

func TestTrivialNaturalTransformation(t *testing.T) {
	objects := []kitty.Object{"a", "b"}
	morphisms := []*kitty.Morphism{
		{"f", "a", "b"},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "a",
		"b": "b",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "f",
	}
	f, err := kitty.NewFunctor(c, c, objMap, mMap)
	require.NoError(t, err)

	comp := map[kitty.Object]kitty.MorphismID{
		"a": kitty.Identity,
		"b": kitty.Identity,
	}
	_, err = kitty.NewNaturalTransformation(f, f, comp)
	require.NoError(t, err)
}

func TestDifferentFunctorNaturalTransformation(t *testing.T) {
	objects1 := []kitty.Object{"a", "b", "c"}
	morphisms1 := []*kitty.Morphism{
		{"f", "a", "b"},
		{"g", "b", "c"},
	}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c1, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"a", "b"}
	morphisms2 := []*kitty.Morphism{
		{"p", "a", "b"},
		{"p^{-1}", "b", "a"},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"p", "p^{-1}"}: kitty.Identity,
		{"p^{-1}", "p"}: kitty.Identity,
	}
	c2, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap1 := map[kitty.Object]kitty.Object{
		"a": "a",
		"b": "b",
		"c": "b",
	}
	mMap1 := map[kitty.MorphismID]kitty.MorphismID{
		"f": "p",
		"g": kitty.Identity,
	}
	f, err := kitty.NewFunctor(c1, c2, objMap1, mMap1)
	require.NoError(t, err)

	objMap2 := map[kitty.Object]kitty.Object{
		"a": "a",
		"b": "a",
		"c": "b",
	}
	mMap2 := map[kitty.MorphismID]kitty.MorphismID{
		"f": kitty.Identity,
		"g": "p",
	}
	g, err := kitty.NewFunctor(c1, c2, objMap2, mMap2)
	require.NoError(t, err)

	comp := map[kitty.Object]kitty.MorphismID{
		"a": kitty.Identity,
		"b": "p^{-1}",
		"c": kitty.Identity,
	}
	_, err = kitty.NewNaturalTransformation(f, g, comp)
	require.NoError(t, err)
}
