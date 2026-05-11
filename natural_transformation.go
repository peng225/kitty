package kitty

import (
	"fmt"
)

type NaturalTransformation struct {
	Source      *Functor
	Destination *Functor

	Components map[Object]MorphismID
}

func NewNaturalTransformation(
	f, g *Functor, comp map[Object]MorphismID) (*NaturalTransformation, error) {
	nt := &NaturalTransformation{
		Source:      f,
		Destination: g,
		Components:  comp,
	}
	for obj, mID := range nt.Components {
		if mID == Identity {
			nt.Components[obj] = f.MapObject(obj).GetIdentityID()
		}
	}
	err := nt.validate()
	if err != nil {
		return nil, err
	}
	return nt, nil
}

// If F=Source and G=Destination, then:
// F(a) -- F(f) --> F(b)
//  |                |
// nt_a             nt_b
//  ↓                ↓
// G(a) -- G(f) --> G(b)

func (nt *NaturalTransformation) validate() error {
	for _, m := range nt.Source.Source.Morphisms {
		a := m.Source
		b := m.Destination

		etaA := nt.Components[a]
		etaB := nt.Components[b]

		Ff := nt.Source.MapMorphism(m.ID)
		Gf := nt.Destination.MapMorphism(m.ID)

		left, err := nt.Source.Destination.Compose(etaA, Gf)
		if err != nil {
			return err
		}
		right, err := nt.Source.Destination.Compose(Ff, etaB)
		if err != nil {
			return err
		}
		if left != right {
			return fmt.Errorf("naturality failed: left=%s, right=%s", left, right)
		}
	}
	return nil
}

// Calculate nt2◦nt.
// If F=nt.Source, G=nt.Destination, F'=nt2.Source, and G'=nt2.Destination, then:
// F(a) -- nt_a --> G(a)
// F'(a) -- nt2_a --> G'(a)
// It is supposed that G(a) == F'(a),
// and nt2_a ◦ nt_a: F(a)↦G'(a) is calculated for each 'a'.
func (nt *NaturalTransformation) ComposeWith(nt2 *NaturalTransformation) (*NaturalTransformation, error) {
	components := make(map[Object]MorphismID)
	for obj := range nt.Components {
		a := nt.Components[obj]
		b, ok := nt2.Components[obj]
		if !ok {
			return nil, fmt.Errorf("nt2 does not have %s component", obj)
		}
		comp, err := nt.Source.Destination.Compose(a, b)
		if err != nil {
			return nil, err
		}
		components[obj] = comp
	}

	res, err := NewNaturalTransformation(nt.Source, nt2.Destination, components)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Calculate F◦nt.
// If G=nt.Source and H=nt.Destination,
// then the resulting natural transformation is:
// F◦nt: F∘G ⇒ F∘H
func (nt *NaturalTransformation) ApplyFunctorLeft(F *Functor) (*NaturalTransformation, error) {
	FAndNTSource, err := nt.Source.ComposeWith(F)
	if err != nil {
		return nil, err
	}
	FAndNTDestination, err := nt.Destination.ComposeWith(F)
	if err != nil {
		return nil, err
	}

	components := make(map[Object]MorphismID)
	for obj, m := range nt.Components {
		// The 'a' component of F◦nt is F(nt_a).
		components[obj] = F.MapMorphism(m)
	}
	return NewNaturalTransformation(FAndNTSource, FAndNTDestination, components)
}

// Calculate nt◦F.
// If G=nt.Source and H=nt.Destination,
// then the resulting natural transformation is:
// nt◦F: G∘F ⇒ H∘F
func (nt *NaturalTransformation) ApplyFunctorRight(F *Functor) (*NaturalTransformation, error) {
	NTSourceAndF, err := F.ComposeWith(nt.Source)
	if err != nil {
		return nil, err
	}
	NTDestinationAndF, err := F.ComposeWith(nt.Destination)
	if err != nil {
		return nil, err
	}

	components := make(map[Object]MorphismID)
	for obj := range F.Source.Objects {
		// The 'a' component of nt◦F is nt_F(a).
		imageObj := F.MapObject(obj)
		components[obj] = nt.Components[imageObj]
	}
	return NewNaturalTransformation(NTSourceAndF, NTDestinationAndF, components)
}

func IdentityNaturalTransformation(F *Functor) (*NaturalTransformation, error) {
	comp := make(map[Object]MorphismID)
	for obj := range F.Source.Objects {
		// The 'a' component of 1_F is 1_F(a).
		comp[obj] = F.MapMorphism(obj.GetIdentityID())
	}
	nt, err := NewNaturalTransformation(F, F, comp)
	if err != nil {
		return nil, err
	}

	return nt, nil
}

func (nt *NaturalTransformation) IsIsomorphic() bool {
	for _, mID := range nt.Components {
		C := nt.Destination.Destination
		m := C.Morphisms[mID]
		if !C.IsIsomorphic(m.Source, m.Destination) {
			return false
		}
	}
	return true
}
