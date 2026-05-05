package kitty

import (
	"fmt"
)

type Functor struct {
	source      *Category
	destination *Category
	objectMap   map[Object]Object
	morphismMap map[MorphismID]MorphismID
}

func NewFunctor(
	src, dst *Category, objectMap map[Object]Object, morphismMap map[MorphismID]MorphismID,
) (*Functor, error) {
	f := &Functor{
		source:      src,
		destination: dst,
		objectMap:   objectMap,
		morphismMap: morphismMap,
	}
	for srcObj, dstObj := range f.objectMap {
		f.morphismMap[srcObj.GetIdentityID()] = dstObj.GetIdentityID()
	}
	err := f.validate()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Functor) validate() error {
	// Type preservation
	for id, m := range f.source.Morphisms {
		fm, ok := f.destination.Morphisms[f.morphismMap[id]]
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
	for _, m := range f.source.Morphisms {
		if !m.IsIdentity() {
			continue
		}
		if !f.destination.Morphisms[f.morphismMap[m.ID]].IsIdentity() {
			return fmt.Errorf("identity not preserved: src morphism=%s, dst morphism=%s",
				m.ID, f.morphismMap[m.ID])
		}
	}

	// Composition preservation
	for _, f1 := range f.source.Morphisms {
		for _, f2 := range f.source.Morphisms {
			// Check if f2◦f1 is a valid morphism.
			if f1.Destination != f2.Source {
				continue
			}
			sourceComp, err := f.source.Compose(f1.ID, f2.ID)
			if err != nil {
				return err
			}

			Ff1 := f.morphismMap[f1.ID]
			Ff2 := f.morphismMap[f2.ID]
			destComp, err := f.destination.Compose(Ff1, Ff2)
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
