package kitty

import "errors"

type NaturalTransformation struct {
	F *Functor
	G *Functor

	Components map[Object]MorphismID
}

func NewNaturalTransformation(
	f, g *Functor, comp map[Object]MorphismID) (*NaturalTransformation, error) {
	nt := &NaturalTransformation{
		F:          f,
		G:          g,
		Components: comp,
	}
	for obj, mID := range nt.Components {
		if mID == Identity {
			nt.Components[obj] = f.objectMap[obj].GetIdentityID()
		}
	}
	err := nt.validate()
	if err != nil {
		return nil, err
	}
	return nt, nil
}

// F(a) -- F(f) --> F(b)
//  |                |
// η_a              η_b
//  ↓                ↓
// G(a) -- G(f) --> G(b)

func (nt *NaturalTransformation) validate() error {
	for _, m := range nt.F.Source.Morphisms {
		a := m.Source
		b := m.Destination

		etaA := nt.Components[a]
		etaB := nt.Components[b]

		Ff := nt.F.MapMorphism(m.ID)
		Gf := nt.G.MapMorphism(m.ID)

		left, err := nt.F.Destination.Compose(etaA, Gf)
		if err != nil {
			return err
		}
		right, err := nt.F.Destination.Compose(Ff, etaB)
		if err != nil {
			return err
		}
		if left != right {
			return errors.New("naturality failed")
		}
	}
	return nil
}
