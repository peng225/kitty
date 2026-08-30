package kitty

import (
	"fmt"
	"maps"
)

// For Cocone, components is the type of:
// I -> (Diagram(I) -> Vertex)
type Cocone Cone

func NewCocone(diagram *Functor, vertex Object, components map[Object]MorphismID) (*Cocone, error) {
	cocone := &Cocone{
		Diagram:    diagram,
		Vertex:     vertex,
		Components: components,
	}
	err := cocone.validate()
	if err != nil {
		return nil, err
	}
	return cocone, nil
}

func (c *Cocone) validate() error {
	shape := c.Diagram.Source
	C := c.Diagram.Destination

	for i := range shape.Objects {
		mID, ok := c.Components[i]
		if !ok {
			return fmt.Errorf("missing cocone component for %s", i)
		}

		m, ok := C.Morphisms[mID]
		if !ok {
			return fmt.Errorf("morphism with the ID %s not found", mID)
		}

		if m.Source != c.Diagram.MapObject(i) {
			return fmt.Errorf(
				"cocone component %s has wrong source: expected = %s, actual = %s",
				mID, c.Diagram.MapObject(i), m.Source,
			)
		}

		if m.Destination != c.Vertex {
			return fmt.Errorf(
				"cocone component %s does not end at vertex: vertex = %s, destination = %s",
				mID, c.Vertex, m.Destination,
			)
		}
	}

	// Check naturality condition.
	for fID, f := range shape.Morphisms {
		i := f.Source
		j := f.Destination

		lambdaI := c.Components[i]
		lambdaJ := c.Components[j]

		Df := c.Diagram.MapMorphism(fID)

		left, err := C.Compose(Df, lambdaJ)
		if err != nil {
			return fmt.Errorf("invalid cocone composition: %w", err)
		}

		if left != lambdaI {
			return fmt.Errorf(
				"cocone condition failed for %s: expected = %s, actual = %s",
				fID, lambdaI, left,
			)
		}
	}

	return nil
}

func (cc *Cocone) ToCone() *Cone {
	Dop := cc.Diagram.Opposite()

	cone, err := NewCone(Dop, cc.Vertex, maps.Clone(cc.Components))
	if err != nil {
		panic(err)
	}

	return cone
}

func (cc *Cocone) Morphism(
	to *Cocone,
) []MorphismID {
	fromCone := cc.ToCone()
	toCone := to.ToCone()

	// Cocone:
	//     from(cc) --f--> to
	//
	// => Cone:
	//     to^op --f^op--> from^op
	return toCone.Morphism(fromCone)
}

func (cc *Cocone) IsColimit() bool {
	return cc.ToCone().IsLimit()
}

func FindColimits(Diagram *Functor) []*Cocone {
	Dop := Diagram.Opposite()
	limits := FindLimits(Dop)

	result := make([]*Cocone, 0, len(limits))

	for _, cone := range limits {
		cocone, err := NewCocone(Diagram, cone.Vertex, maps.Clone(cone.Components))
		if err != nil {
			panic(err)
		}
		result = append(result, cocone)
	}

	return result
}
