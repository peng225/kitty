package kitty_test

import (
	"testing"

	"github.com/peng225/kitty"
	"github.com/stretchr/testify/require"
)

func TestConeWithDiscreteShape(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	shape, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q", "r"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "g1",
			Source:      "p",
			Destination: "q",
		},
		{
			ID:          "g2",
			Source:      "p",
			Destination: "r",
		},
		{
			ID:          "g3",
			Source:      "q",
			Destination: "r",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	C, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "q",
		"b": "r",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{}
	Diagram, err := kitty.NewFunctor(shape, C, objMap, mMap)
	require.NoError(t, err)

	vertex := kitty.Object("p")
	components := map[kitty.Object]kitty.MorphismID{
		"a": "g1",
		"b": "g2",
	}
	_, err = kitty.NewCone(Diagram, vertex, components)
	require.NoError(t, err)
}

func TestConeWithOneMapShape(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{
		{"f1", "a", "b"},
	}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	shape, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q", "r"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "g1",
			Source:      "p",
			Destination: "q",
		},
		{
			ID:          "g2",
			Source:      "p",
			Destination: "r",
		},
		{
			ID:          "g3",
			Source:      "q",
			Destination: "r",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"g1", "g3"}: "g2",
	}
	C, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "q",
		"b": "r",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f1": "g3",
	}
	Diagram, err := kitty.NewFunctor(shape, C, objMap, mMap)
	require.NoError(t, err)

	vertex := kitty.Object("p")
	components := map[kitty.Object]kitty.MorphismID{
		"a": "g1",
		"b": "g2",
	}
	_, err = kitty.NewCone(Diagram, vertex, components)
	require.NoError(t, err)
}
