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
