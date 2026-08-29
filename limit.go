package kitty

import "fmt"

type Cone struct {
	// shape -> C
	Diagram *Functor
	Vertex  Object
	// I -> (Vertex -> Diagram(I))
	Components map[Object]MorphismID
}

func NewCone(diagram *Functor, vertex Object, components map[Object]MorphismID) (*Cone, error) {
	cone := &Cone{
		Diagram:    diagram,
		Vertex:     vertex,
		Components: components,
	}
	err := cone.validate()
	if err != nil {
		return nil, err
	}
	return cone, nil
}

func (c *Cone) validate() error {
	shape := c.Diagram.Source
	C := c.Diagram.Destination

	for i := range shape.Objects {
		mID, ok := c.Components[i]
		if !ok {
			return fmt.Errorf("missing cone component for %s", i)
		}

		m, ok := C.Morphisms[mID]
		if !ok {
			return fmt.Errorf("morphism with the ID %s not found", mID)
		}

		if m.Source != c.Vertex {
			return fmt.Errorf(
				"cone component %s does not start at vertex: vertex = %s, source = %s",
				mID, c.Vertex, m.Source,
			)
		}

		if m.Destination != c.Diagram.MapObject(i) {
			return fmt.Errorf(
				"cone component %s has wrong Destination: expected = %s, actual = %s",
				mID, c.Diagram.MapObject(i), m.Destination,
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

		left, err := C.Compose(lambdaI, Df)
		if err != nil {
			return fmt.Errorf("invalid cone composition: %w", err)
		}

		if left != lambdaJ {
			return fmt.Errorf(
				"cone condition failed for %s: expected = %s, actual = %s",
				fID, lambdaJ, left,
			)
		}
	}

	return nil
}
