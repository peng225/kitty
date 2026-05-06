package kitty_test

import (
	"testing"

	"github.com/peng225/kitty"
	"github.com/stretchr/testify/require"
)

func TestTrivialAdjunction(t *testing.T) {
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
	F, err := kitty.NewFunctor(c, c, objMap, mMap)
	require.NoError(t, err)

	G := F
	GF, err := F.ComposeWith(G)
	require.NoError(t, err)
	FG, err := G.ComposeWith(F)
	require.NoError(t, err)

	idFSource, err := kitty.IdentityFunctor(F.Source)
	require.NoError(t, err)
	idGSource, err := kitty.IdentityFunctor(G.Source)
	require.NoError(t, err)
	comp := map[kitty.Object]kitty.MorphismID{
		"a": kitty.Identity,
		"b": kitty.Identity,
	}
	eta, err := kitty.NewNaturalTransformation(idFSource, GF, comp)
	require.NoError(t, err)
	epsilon, err := kitty.NewNaturalTransformation(FG, idGSource, comp)
	require.NoError(t, err)

	err = kitty.CheckAdjunction(F, G, eta, epsilon)
	require.NoError(t, err)
}

func TestNonTrivialAdjunction(t *testing.T) {
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

	objects2 := []kitty.Object{"x", "y"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "p",
			Source:      "x",
			Destination: "y",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	c2, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap1 := map[kitty.Object]kitty.Object{
		"a": "x",
		"b": "x",
		"c": "y",
	}
	mMap1 := map[kitty.MorphismID]kitty.MorphismID{
		"f": kitty.Identity,
		"g": "p",
	}
	F, err := kitty.NewFunctor(c1, c2, objMap1, mMap1)
	require.NoError(t, err)

	objMap2 := map[kitty.Object]kitty.Object{
		"x": "b",
		"y": "c",
	}
	mMap2 := map[kitty.MorphismID]kitty.MorphismID{
		"p": "g",
	}
	G, err := kitty.NewFunctor(c2, c1, objMap2, mMap2)
	require.NoError(t, err)

	GF, err := F.ComposeWith(G)
	require.NoError(t, err)
	FG, err := G.ComposeWith(F)
	require.NoError(t, err)

	idFSource, err := kitty.IdentityFunctor(F.Source)
	require.NoError(t, err)
	idGSource, err := kitty.IdentityFunctor(G.Source)
	require.NoError(t, err)
	comp1 := map[kitty.Object]kitty.MorphismID{
		"a": "f",
		"b": kitty.Identity,
		"c": kitty.Identity,
	}
	eta, err := kitty.NewNaturalTransformation(idFSource, GF, comp1)
	require.NoError(t, err)
	comp2 := map[kitty.Object]kitty.MorphismID{
		"x": kitty.Identity,
		"y": kitty.Identity,
	}
	epsilon, err := kitty.NewNaturalTransformation(FG, idGSource, comp2)
	require.NoError(t, err)

	err = kitty.CheckAdjunction(F, G, eta, epsilon)
	require.NoError(t, err)
}

func TestNotAdjunction(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "a",
		},
		{
			ID:          "g",
			Source:      "b",
			Destination: "b",
		},
	}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"f", "f"}: kitty.Identity,
		{"g", "g"}: kitty.Identity,
	}
	c1, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"x", "y"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "p",
			Source:      "x",
			Destination: "x",
		},
		{
			ID:          "q",
			Source:      "y",
			Destination: "y",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"p", "p"}: kitty.Identity,
		{"q", "q"}: kitty.Identity,
	}
	c2, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap1 := map[kitty.Object]kitty.Object{
		"a": "y",
		"b": "x",
	}
	mMap1 := map[kitty.MorphismID]kitty.MorphismID{
		"f": "q",
		"g": "p",
	}
	F, err := kitty.NewFunctor(c1, c2, objMap1, mMap1)
	require.NoError(t, err)

	objMap2 := map[kitty.Object]kitty.Object{
		"x": "b",
		"y": "a",
	}
	mMap2 := map[kitty.MorphismID]kitty.MorphismID{
		"p": "g",
		"q": "f",
	}
	G, err := kitty.NewFunctor(c2, c1, objMap2, mMap2)
	require.NoError(t, err)

	GF, err := F.ComposeWith(G)
	require.NoError(t, err)
	FG, err := G.ComposeWith(F)
	require.NoError(t, err)

	idFSource, err := kitty.IdentityFunctor(F.Source)
	require.NoError(t, err)
	idGSource, err := kitty.IdentityFunctor(G.Source)
	require.NoError(t, err)
	comp1 := map[kitty.Object]kitty.MorphismID{
		"a": "f",
		"b": "g",
	}
	eta, err := kitty.NewNaturalTransformation(idFSource, GF, comp1)
	require.NoError(t, err)
	comp2 := map[kitty.Object]kitty.MorphismID{
		"x": kitty.Identity,
		"y": kitty.Identity,
	}
	epsilon, err := kitty.NewNaturalTransformation(FG, idGSource, comp2)
	require.NoError(t, err)

	err = kitty.CheckAdjunction(F, G, eta, epsilon)
	require.Error(t, err)
	require.ErrorIs(t, err, kitty.ErrorTriangleIdentityViolation)
}
