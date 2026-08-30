package kitty

import (
	"fmt"
	"maps"
)

type Functor struct {
	Source      *Category
	Destination *Category
	objectMap   map[Object]Object
	morphismMap map[MorphismID]MorphismID
}

func NewFunctor(
	src, dst *Category, objectMap map[Object]Object, morphismMap map[MorphismID]MorphismID,
) (*Functor, error) {
	F := &Functor{
		Source:      src,
		Destination: dst,
		objectMap:   objectMap,
		morphismMap: morphismMap,
	}
	for srcObj, dstObj := range F.objectMap {
		F.morphismMap[srcObj.GetIdentityID()] = dstObj.GetIdentityID()
	}
	for k, v := range F.morphismMap {
		if v == Identity {
			srcObj := F.Source.Morphisms[k].Source
			mappedSrcObj := F.objectMap[srcObj]
			// Here, we proceed under the assumption that F[srcObj] == F[dstObj].
			F.morphismMap[k] = mappedSrcObj.GetIdentityID()
		}
	}
	err := F.constructComposition()
	if err != nil {
		return nil, err
	}
	err = F.validate()
	if err != nil {
		return nil, err
	}
	return F, nil
}

func (F *Functor) constructComposition() error {
	previousMMapCount := len(F.morphismMap)
	for len(F.morphismMap) != len(F.Source.Morphisms) {
		for _, m1 := range F.Source.Morphisms {
			for _, m2 := range F.Source.Morphisms {
				if m1.ID == m2.ID {
					continue
				}
				composed, err := F.Source.Compose(m1.ID, m2.ID)
				if err != nil {
					// FIXME: should define the composition not found error.
					continue
				}
				if _, ok := F.morphismMap[m1.ID]; !ok {
					continue
				}
				if _, ok := F.morphismMap[m2.ID]; !ok {
					continue
				}
				fm1CircFm2, err := F.Destination.Compose(F.morphismMap[m1.ID], F.morphismMap[m2.ID])
				if err != nil {
					return fmt.Errorf("failed to construct composition: %w", err)
				}
				F.morphismMap[composed] = fm1CircFm2
			}
		}
		if previousMMapCount == len(F.morphismMap) {
			return fmt.Errorf("composition construction stuck detected")
		}
		previousMMapCount = len(F.morphismMap)
	}
	return nil
}

func (F *Functor) validate() error {
	// Type preservation
	for id, m := range F.Source.Morphisms {
		Fm, ok := F.Destination.Morphisms[F.morphismMap[id]]
		if !ok {
			return fmt.Errorf("the destination of %s is not defined", id)
		}

		if Fm.Source != F.objectMap[m.Source] {
			return fmt.Errorf("source type mismatch: src=%s, F(src)=%s",
				Fm.Source, F.objectMap[m.Source])
		}
		if Fm.Source != F.objectMap[m.Source] ||
			Fm.Destination != F.objectMap[m.Destination] {
			return fmt.Errorf("destination type mismatch: dst=%s, F(dst)=%s",
				Fm.Destination, F.objectMap[m.Destination])
		}
	}

	// Identities
	for _, m := range F.Source.Morphisms {
		if !m.IsIdentity() {
			continue
		}
		if !F.Destination.Morphisms[F.morphismMap[m.ID]].IsIdentity() {
			return fmt.Errorf("identity not preserved: src morphism=%s, dst morphism=%s",
				m.ID, F.morphismMap[m.ID])
		}
	}

	// Composition preservation
	for _, f1 := range F.Source.Morphisms {
		for _, f2 := range F.Source.Morphisms {
			// Check if f2◦f1 is a valid morphism.
			if f1.Destination != f2.Source {
				continue
			}
			sourceComp, err := F.Source.Compose(f1.ID, f2.ID)
			if err != nil {
				return err
			}

			Ff1 := F.morphismMap[f1.ID]
			Ff2 := F.morphismMap[f2.ID]
			destComp, err := F.Destination.Compose(Ff1, Ff2)
			if err != nil {
				return err
			}

			if F.morphismMap[sourceComp] != destComp {
				return fmt.Errorf("composition not preserved: sourceComp=%s, destComp=%s",
					sourceComp, destComp)
			}
		}
	}

	return nil
}

func (F *Functor) MapObject(o Object) Object {
	return F.objectMap[o]
}

func (F *Functor) MapMorphism(m MorphismID) MorphismID {
	return F.morphismMap[m]
}

func (F *Functor) ComposeWith(G *Functor) (*Functor, error) {
	// If F: A→B and G: B→C, then G∘F: A→C
	objMap := make(map[Object]Object)
	mMap := make(map[MorphismID]MorphismID)
	for obj := range F.Source.Objects {
		objMap[obj] = G.objectMap[F.objectMap[obj]]
	}
	for m := range F.Source.Morphisms {
		mMap[m] = G.morphismMap[F.morphismMap[m]]
	}

	H, err := NewFunctor(
		F.Source, G.Destination, objMap, mMap,
	)
	if err != nil {
		return nil, err
	}

	return H, nil
}

func IdentityFunctor(c *Category) (*Functor, error) {
	objMap := make(map[Object]Object)
	mMap := make(map[MorphismID]MorphismID)
	for obj := range c.Objects {
		objMap[obj] = obj
	}
	for m := range c.Morphisms {
		mMap[m] = m
	}
	F, err := NewFunctor(c, c, objMap, mMap)
	if err != nil {
		return nil, err
	}

	return F, nil
}

func (F *Functor) IsFaithful() bool {
	C := F.Source
	for a := range C.Objects {
		for b := range C.Objects {
			seen := map[MorphismID]bool{}
			for _, f := range C.Hom(a, b) {
				image := F.morphismMap[f]
				if seen[image] {
					return false
				}
				seen[image] = true
			}
		}
	}
	return true
}

func (F *Functor) IsFull() bool {
	C := F.Source
	D := F.Destination
	for a := range C.Objects {
		for b := range C.Objects {
			targetHom := D.Hom(
				F.objectMap[a],
				F.objectMap[b],
			)
			image := map[MorphismID]bool{}
			for _, f := range C.Hom(a, b) {
				image[F.morphismMap[f]] = true
			}
			for _, g := range targetHom {
				if !image[g] {
					return false
				}
			}
		}
	}
	return true
}

func (F *Functor) IsEssentiallySurjective() bool {
	D := F.Destination
	for x := range D.Objects {
		found := false
		for a := range F.Source.Objects {
			if D.IsIsomorphic(F.objectMap[a], x) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (F *Functor) IsEquivalence() bool {
	return F.IsFaithful() && F.IsFull() &&
		F.IsEssentiallySurjective()
}

func (F *Functor) Opposite() *Functor {
	Cop := F.Source.Opposite()
	Dop := F.Destination.Opposite()

	result, err := NewFunctor(
		Cop, Dop,
		maps.Clone(F.objectMap),
		maps.Clone(F.morphismMap), // F(f) = g => F^op(f^op) = g^op
	)
	if err != nil {
		panic(err)
	}

	return result
}
