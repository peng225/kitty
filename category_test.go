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

func TestInvalidMorphism(t *testing.T) {
	objects := []kitty.Object{"a", "b"}

	tests := map[string]struct {
		m []*kitty.Morphism
	}{
		"invalid source": {
			m: []*kitty.Morphism{
				{
					ID:          "f",
					Source:      "x",
					Destination: "b",
				},
			},
		},
		"invalid destination": {
			m: []*kitty.Morphism{
				{
					ID:          "f",
					Source:      "a",
					Destination: "x",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
			_, err := kitty.NewCategory(objects, test.m, compose)
			require.Error(t, err)
		})
	}
}

func TestComposition(t *testing.T) {
	objects := []kitty.Object{"a", "b", "c"}
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
			Source:      "b",
			Destination: "c",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
}

func TestAssociativeLawViolation(t *testing.T) {
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
		{
			ID:          "gf",
			Source:      "a",
			Destination: "c",
		},
		{
			ID:          "hgf",
			Source:      "a",
			Destination: "d",
		},
		{
			ID:          "hg",
			Source:      "b",
			Destination: "d",
		},
		{
			ID:          "hgf2",
			Source:      "a",
			Destination: "d",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"f", "g"}:  "gf",
		{"g", "h"}:  "hg",
		{"f", "hg"}: "hgf",
		{"gf", "h"}: "hgf2",
	}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.Error(t, err)
}

func TestUnknownMorphismInComposition(t *testing.T) {
	objects := []kitty.Object{"a", "b", "c"}
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
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"f", "g"}: "gf",
	}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.Error(t, err)
}

func TestInverse(t *testing.T) {
	objects := []kitty.Object{"a", "b"}
	f := &kitty.Morphism{
		ID:          "f",
		Source:      "a",
		Destination: "b",
	}
	morphisms := []*kitty.Morphism{f, f.Inverse()}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
}

func TestMorphismLoop(t *testing.T) {
	objects := []kitty.Object{"a", "b"}
	morphisms := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
		{
			ID:          "g",
			Source:      "b",
			Destination: "a",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
}

func TestComplicatedInverseChain(t *testing.T) {
	objects := []kitty.Object{"a", "b", "c"}
	f := &kitty.Morphism{
		ID:          "f",
		Source:      "a",
		Destination: "b",
	}
	g := &kitty.Morphism{
		ID:          "g",
		Source:      "b",
		Destination: "a",
	}
	morphisms := []*kitty.Morphism{
		f, f.Inverse(), g, g.Inverse(),
		{
			ID:          "h",
			Source:      "a",
			Destination: "c",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	_, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
}

func TestCompositionReductionToAnotherOne(t *testing.T) {
	objects := []kitty.Object{"a", "b", "c"}
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
			ID:          "gf",
			Source:      "a",
			Destination: "c",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"f", "g"}: "gf",
	}
	C, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)

	gf, err := C.Compose("f", "g")
	require.NoError(t, err)
	require.Equal(t, kitty.MorphismID("gf"), gf)
}

func TestCompositionReductionToIdentity(t *testing.T) {
	a := kitty.Object("a")
	b := kitty.Object("b")
	objects := []kitty.Object{a, b}
	morphisms := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
		{
			ID:          "g",
			Source:      "b",
			Destination: "a",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"f", "g"}: kitty.Identity,
	}
	C, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)

	gf, err := C.Compose("f", "g")
	require.NoError(t, err)
	require.Equal(t, a.GetIdentityID(), gf)
}

func TestCompositionComplicatedReduction(t *testing.T) {
	objects := []kitty.Object{"p", "q", "r", "s", "t"}
	g7 := &kitty.Morphism{
		ID:          "g7",
		Source:      "p",
		Destination: "q",
	}
	morphisms := []*kitty.Morphism{
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
			Source:      "p",
			Destination: "t",
		},
		{
			ID:          "g4",
			Source:      "q",
			Destination: "r",
		},
		{
			ID:          "g5",
			Source:      "q",
			Destination: "s",
		},
		{
			ID:          "g6",
			Source:      "q",
			Destination: "t",
		},
		g7, g7.Inverse(),
		{
			ID:          "g8",
			Source:      "r",
			Destination: "t",
		},
		{
			ID:          "g9",
			Source:      "s",
			Destination: "t",
		},
	}
	compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{
		{"g7", "g4"}: "g1",
		{"g7", "g5"}: "g2",
		{"g7", "g6"}: "g3",
		{"g1", "g8"}: "g3",
		{"g2", "g9"}: "g3",
		{"g4", "g8"}: "g6",
		{"g5", "g9"}: "g6",
	}
	_, err := kitty.NewCategory(objects, morphisms, compose2)
	require.NoError(t, err)
}

func TestHom(t *testing.T) {
	objects := []kitty.Object{"a", "b"}
	morphisms := []*kitty.Morphism{
		{
			ID:          "f",
			Source:      "a",
			Destination: "b",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	C, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)

	hom := C.Hom("a", "b")
	require.Len(t, hom, 1)
	require.Equal(t, kitty.MorphismID("f"), hom[0])
	hom = C.Hom("a", "a")
	require.Len(t, hom, 1)
	require.Equal(t, kitty.MorphismID("_id_a"), hom[0])
	hom = C.Hom("b", "b")
	require.Len(t, hom, 1)
	require.Equal(t, kitty.MorphismID("_id_b"), hom[0])
}

func TestIsomorphic(t *testing.T) {
	objects := []kitty.Object{"a", "b", "c"}
	f := &kitty.Morphism{
		ID:          "f",
		Source:      "a",
		Destination: "b",
	}
	morphisms := []*kitty.Morphism{
		f, f.Inverse(),
		{
			ID:          "g",
			Source:      "b",
			Destination: "c",
		},
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	C, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
	require.True(t, C.IsIsomorphic("a", "b"))
	require.True(t, C.IsIsomorphic("b", "a"))
	require.False(t, C.IsIsomorphic("b", "c"))
	require.False(t, C.IsIsomorphic("c", "b"))
}

func TestOpposite(t *testing.T) {
	objects := []kitty.Object{"a", "b", "c"}
	f := &kitty.Morphism{
		ID:          "f",
		Source:      "a",
		Destination: "b",
	}
	g := &kitty.Morphism{
		ID:          "g",
		Source:      "b",
		Destination: "c",
	}
	morphisms := []*kitty.Morphism{
		f, g,
	}
	compose := map[[2]kitty.MorphismID]kitty.MorphismID{}
	C, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
	D := C.Opposite()
	_, ok := D.Objects["a"]
	require.True(t, ok)
	_, ok = D.Objects["b"]
	require.True(t, ok)
	fop, ok := D.Morphisms["f"]
	require.True(t, ok)
	require.Equal(t, kitty.Object("b"), fop.Source)
	require.Equal(t, kitty.Object("a"), fop.Destination)
	gop, ok := D.Morphisms["g"]
	require.True(t, ok)
	require.Equal(t, kitty.Object("c"), gop.Source)
	require.Equal(t, kitty.Object("b"), gop.Destination)
}
