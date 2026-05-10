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
	C, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "a",
		"b": "b",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "f",
	}
	F, err := kitty.NewFunctor(C, C, objMap, mMap)
	require.NoError(t, err)

	comp := map[kitty.Object]kitty.MorphismID{
		"a": kitty.Identity,
		"b": kitty.Identity,
	}
	_, err = kitty.NewNaturalTransformation(F, F, comp)
	require.NoError(t, err)
}

func TestDifferentFunctorNaturalTransformation(t *testing.T) {
	objects1 := []kitty.Object{"a", "b", "c"}
	morphisms1 := []*kitty.Morphism{
		{"f", "a", "b"},
		{"g", "b", "c"},
	}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	C1, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"a", "b"}
	p := &kitty.Morphism{
		ID:          "p",
		Source:      "a",
		Destination: "b",
	}
	morphisms2 := []*kitty.Morphism{
		p,
		p.Inverse(),
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	C2, err := kitty.NewCategory(objects2, morphisms2, compose2)
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
	F, err := kitty.NewFunctor(C1, C2, objMap1, mMap1)
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
	G, err := kitty.NewFunctor(C1, C2, objMap2, mMap2)
	require.NoError(t, err)

	comp := map[kitty.Object]kitty.MorphismID{
		"a": kitty.Identity,
		"b": p.Inverse().ID,
		"c": kitty.Identity,
	}
	_, err = kitty.NewNaturalTransformation(F, G, comp)
	require.NoError(t, err)
}
