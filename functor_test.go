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
	C1, err := kitty.NewCategory(objects, morphisms, compose)
	require.NoError(t, err)
	C2 := C1

	objMap := map[kitty.Object]kitty.Object{
		"a": "a",
		"b": "b",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "f",
	}
	_, err = kitty.NewFunctor(C1, C2, objMap, mMap)
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
	C1, err := kitty.NewCategory(objects1, morphisms1, compose1)
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
	C2, err := kitty.NewCategory(objects2, morphisms2, compose2)
	require.NoError(t, err)

	objMap := map[kitty.Object]kitty.Object{
		"a": "p",
		"b": "q",
	}
	mMap := map[kitty.MorphismID]kitty.MorphismID{
		"f": "g",
	}
	F, err := kitty.NewFunctor(C1, C2, objMap, mMap)
	require.NoError(t, err)
	require.True(t, F.IsEquivalence())
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
	C1, err := kitty.NewCategory(objects1, morphisms1, compose1)
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
	C2, err := kitty.NewCategory(objects2, morphisms2, compose2)
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
	F, err := kitty.NewFunctor(C1, C2, objMap, mMap)
	require.NoError(t, err)
	require.False(t, F.IsEquivalence())
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
	C, err := kitty.NewCategory(objects, morphisms, compose)
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
	_, err = kitty.NewFunctor(C, C, objMap, mMap)
	require.NoError(t, err)
}

func TestFunctorProperties(t *testing.T) {
	f := &kitty.Morphism{
		ID:          "f",
		Source:      "a",
		Destination: "b",
	}
	g := &kitty.Morphism{
		ID:          "g",
		Source:      "a",
		Destination: "b",
	}
	u := &kitty.Morphism{
		ID:          "u",
		Source:      "p",
		Destination: "q",
	}
	v := &kitty.Morphism{
		ID:          "v",
		Source:      "p",
		Destination: "q",
	}
	w := &kitty.Morphism{
		ID:          "w",
		Source:      "q",
		Destination: "r",
	}
	tests := map[string]struct {
		objs1                   []kitty.Object
		morphisms1              []*kitty.Morphism
		objs2                   []kitty.Object
		morphisms2              []*kitty.Morphism
		objMap                  map[kitty.Object]kitty.Object
		mMap                    map[kitty.MorphismID]kitty.MorphismID
		isFaithful              bool
		isFull                  bool
		isEssentiallySurjective bool
	}{
		"faithful, not full, essentially surjective": {
			objs1:      []kitty.Object{"a", "b"},
			morphisms1: []*kitty.Morphism{f},
			objs2:      []kitty.Object{"p", "q", "r"},
			morphisms2: []*kitty.Morphism{u, v, w, w.Inverse()},
			objMap: map[kitty.Object]kitty.Object{
				"a": "p",
				"b": "q",
			},
			mMap: map[kitty.MorphismID]kitty.MorphismID{
				"f": "u",
			},
			isFaithful:              true,
			isFull:                  false,
			isEssentiallySurjective: true,
		},
		"not faithful, full, not essentially surjective": {
			objs1:      []kitty.Object{"a", "b"},
			morphisms1: []*kitty.Morphism{f, g},
			objs2:      []kitty.Object{"p", "q", "r"},
			morphisms2: []*kitty.Morphism{u, w},
			objMap: map[kitty.Object]kitty.Object{
				"a": "p",
				"b": "q",
			},
			mMap: map[kitty.MorphismID]kitty.MorphismID{
				"f": "u",
				"g": "u",
			},
			isFaithful:              false,
			isFull:                  true,
			isEssentiallySurjective: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			compose1 := map[[2]kitty.MorphismID]kitty.MorphismID{}
			C1, err := kitty.NewCategory(test.objs1, test.morphisms1, compose1)
			require.NoError(t, err)
			compose2 := map[[2]kitty.MorphismID]kitty.MorphismID{}
			C2, err := kitty.NewCategory(test.objs2, test.morphisms2, compose2)
			require.NoError(t, err)

			F, err := kitty.NewFunctor(C1, C2, test.objMap, test.mMap)
			require.NoError(t, err)
			require.Equal(t, test.isFaithful, F.IsFaithful())
			require.Equal(t, test.isFull, F.IsFull())
			require.Equal(t, test.isEssentiallySurjective, F.IsEssentiallySurjective())
		})
	}
}
