package kitty

import (
	"fmt"
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
	f := &Functor{
		Source:      src,
		Destination: dst,
		objectMap:   objectMap,
		morphismMap: morphismMap,
	}
	for srcObj, dstObj := range f.objectMap {
		f.morphismMap[srcObj.GetIdentityID()] = dstObj.GetIdentityID()
	}
	for k, v := range f.morphismMap {
		if v == Identity {
			srcObj := f.Source.Morphisms[k].Source
			mappedSrcObj := f.objectMap[srcObj]
			// Here, we proceed under the assumption that F[srcObj] == F[dstObj].
			f.morphismMap[k] = mappedSrcObj.GetIdentityID()
		}
	}
	err := f.validate()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Functor) validate() error {
	// Type preservation
	for id, m := range f.Source.Morphisms {
		fm, ok := f.Destination.Morphisms[f.morphismMap[id]]
		if !ok {
			return fmt.Errorf("the destination of %s is not defined.", id)
		}

		if fm.Source != f.objectMap[m.Source] {
			return fmt.Errorf("source type mismatch: src=%s, F(src)=%s",
				fm.Source, f.objectMap[m.Source])
		}
		if fm.Source != f.objectMap[m.Source] ||
			fm.Destination != f.objectMap[m.Destination] {
			return fmt.Errorf("destination type mismatch: dst=%s, F(dst)=%s",
				fm.Destination, f.objectMap[m.Destination])
		}
	}

	// Identities
	for _, m := range f.Source.Morphisms {
		if !m.IsIdentity() {
			continue
		}
		if !f.Destination.Morphisms[f.morphismMap[m.ID]].IsIdentity() {
			return fmt.Errorf("identity not preserved: src morphism=%s, dst morphism=%s",
				m.ID, f.morphismMap[m.ID])
		}
	}

	// Composition preservation
	for _, f1 := range f.Source.Morphisms {
		for _, f2 := range f.Source.Morphisms {
			// Check if f2◦f1 is a valid morphism.
			if f1.Destination != f2.Source {
				continue
			}
			sourceComp, err := f.Source.Compose(f1.ID, f2.ID)
			if err != nil {
				return err
			}

			Ff1 := f.morphismMap[f1.ID]
			Ff2 := f.morphismMap[f2.ID]
			destComp, err := f.Destination.Compose(Ff1, Ff2)
			if err != nil {
				return err
			}

			if f.morphismMap[sourceComp] != destComp {
				return fmt.Errorf("composition not preserved: sourceComp=%s, destComp=%s",
					sourceComp, destComp)
			}
		}
	}

	return nil
}

func (f *Functor) MapMorphism(m MorphismID) MorphismID {
	return f.morphismMap[m]
}
