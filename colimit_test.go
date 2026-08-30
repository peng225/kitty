package kitty_test

import (
	"testing"

	"github.com/peng225/kitty"
	"github.com/stretchr/testify/require"
)

func TestCoconeWithDiscreteShape(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	shape, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q", "r"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "g1",
			Source:      "q",
			Destination: "p",
		},
		{
			ID:          "g2",
			Source:      "r",
			Destination: "p",
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
	_, err = kitty.NewCocone(Diagram, vertex, components)
	require.NoError(t, err)
}

func TestCoconeWithOneMapShape(t *testing.T) {
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
			Source:      "q",
			Destination: "p",
		},
		{
			ID:          "g2",
			Source:      "r",
			Destination: "p",
		},
		{
			ID:          "g3",
			Source:      "q",
			Destination: "r",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"g3", "g2"}: "g1",
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
	_, err = kitty.NewCocone(Diagram, vertex, components)
	require.NoError(t, err)
}

func TestColimit(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	shape, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q", "r", "s"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "g1",
			Source:      "r",
			Destination: "p",
		},
		{
			ID:          "g2",
			Source:      "s",
			Destination: "p",
		},
		{
			ID:          "g3",
			Source:      "r",
			Destination: "q",
		},
		{
			ID:          "g4",
			Source:      "s",
			Destination: "q",
		},
		{
			ID:          "g5",
			Source:      "p",
			Destination: "q",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"g1", "g5"}: "g3",
		{"g2", "g5"}: "g4",
	}
	C, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "r",
		"b": "s",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{}
	Diagram, err := kitty.NewFunctor(shape, C, objMap, mMap)
	require.NoError(t, err)

	vertex := kitty.Object("p")
	components := map[kitty.Object]kitty.MorphismID{
		"a": "g1",
		"b": "g2",
	}
	cocone1, err := kitty.NewCocone(Diagram, vertex, components)
	require.NoError(t, err)

	vertex = kitty.Object("q")
	components = map[kitty.Object]kitty.MorphismID{
		"a": "g3",
		"b": "g4",
	}
	cocone2, err := kitty.NewCocone(Diagram, vertex, components)
	require.NoError(t, err)

	mIDs := cocone1.Morphism(cocone2)
	require.Len(t, mIDs, 1)
	require.Equal(t, kitty.MorphismID("g5"), mIDs[0])

	require.True(t, cocone1.IsColimit())
	require.False(t, cocone2.IsColimit())

	colimits := kitty.FindColimits(Diagram)
	require.Len(t, colimits, 1)
	colimit := colimits[0]
	require.Equal(t, cocone1.Vertex, colimit.Vertex)
	require.Equal(t, cocone1.Components, colimit.Components)
}
