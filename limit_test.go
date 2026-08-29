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

func TestEnumerateCones(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
	}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	C1, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q", "r", "s"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "g1",
			Source:      "p",
			Destination: "r",
		},
		{
			ID:          "g2",
			Source:      "p",
			Destination: "s",
		},
		{
			ID:          "g3",
			Source:      "q",
			Destination: "r",
		},
		{
			ID:          "g4",
			Source:      "q",
			Destination: "s",
		},
		{
			ID:          "g5",
			Source:      "r",
			Destination: "s",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"g1", "g5"}: "g2",
		{"g3", "g5"}: "g4",
	}
	C2, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "r",
		"b": "s",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "g5",
	}
	F, err := kitty.NewFunctor(C1, C2, objMap, mMap)
	require.NoError(t, err)

	cones := kitty.EnumerateCones(F)
	require.Len(t, cones, 3)
	for _, cone := range cones {
		require.Len(t, cone.Components, 2)
		switch cone.Vertex {
		case "p":
			require.Equal(t, kitty.MorphismID("g1"), cone.Components["a"])
			require.Equal(t, kitty.MorphismID("g2"), cone.Components["b"])
		case "q":
			require.Equal(t, kitty.MorphismID("g3"), cone.Components["a"])
			require.Equal(t, kitty.MorphismID("g4"), cone.Components["b"])
		case "r":
			require.True(t, C2.Morphisms[cone.Components["a"]].IsIdentity())
			require.Equal(t, kitty.MorphismID("g5"), cone.Components["b"])
		default:
			t.FailNow()
		}
	}
}

func TestLimit(t *testing.T) {
	objects1 := []kitty.Object{"a", "b"}
	morphisms1 := []*kitty.Morphism{}
	compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
	shape, err := kitty.NewCategory(objects1, morphisms1, compose1)
	require.NoError(t, err)

	objects2 := []kitty.Object{"p", "q", "r", "s"}
	morphisms2 := []*kitty.Morphism{
		{
			ID:          "g1",
			Source:      "p",
			Destination: "r",
		},
		{
			ID:          "g2",
			Source:      "p",
			Destination: "s",
		},
		{
			ID:          "g3",
			Source:      "q",
			Destination: "r",
		},
		{
			ID:          "g4",
			Source:      "q",
			Destination: "s",
		},
		{
			ID:          "g5",
			Source:      "p",
			Destination: "q",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"g5", "g3"}: "g1",
		{"g5", "g4"}: "g2",
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
	cone1, err := kitty.NewCone(Diagram, vertex, components)
	require.NoError(t, err)

	vertex = kitty.Object("q")
	components = map[kitty.Object]kitty.MorphismID{
		"a": "g3",
		"b": "g4",
	}
	cone2, err := kitty.NewCone(Diagram, vertex, components)
	require.NoError(t, err)

	mIDs := cone1.Morphism(cone2)
	require.Len(t, mIDs, 1)
	require.Equal(t, kitty.MorphismID("g5"), mIDs[0])

	require.False(t, cone1.IsLimit())
	require.True(t, cone2.IsLimit())

	limits := kitty.FindLimits(Diagram)
	require.Len(t, mIDs, 1)
	limit := limits[0]
	require.Equal(t, cone2.Vertex, limit.Vertex)
	require.Equal(t, cone2.Components, limit.Components)
}
