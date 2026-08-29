package kitty

import (
	"fmt"
	"maps"
)

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

func EnumerateCones(Diagram *Functor) []*Cone {
	shape := Diagram.Source
	C := Diagram.Destination

	var result []*Cone

	objects := make([]Object, 0, len(shape.Objects))
	for obj := range shape.Objects {
		objects = append(objects, obj)
	}

	var search func(
		int,
		Object,
		map[Object]MorphismID,
	) error

	search = func(
		i int,
		vertex Object,
		components map[Object]MorphismID,
	) error {
		if i == len(objects) {
			cone, err := NewCone(Diagram, vertex, maps.Clone(components))

			if err != nil {
				return err
			}
			result = append(result, cone)

			return nil
		}

		obj := objects[i]
		Fobj := Diagram.MapObject(obj)

		for _, mID := range C.Hom(vertex, Fobj) {
			components[obj] = mID

			err := search(
				i+1,
				vertex,
				components,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for vertex := range C.Objects {
		err := search(
			0,
			vertex,
			map[Object]MorphismID{},
		)
		if err != nil {
			panic(err)
		}
	}

	return result
}

func (c *Cone) Morphism(
	to *Cone,
) []MorphismID {
	C := c.Diagram.Destination

	var result []MorphismID

	for _, u := range C.Hom(c.Vertex, to.Vertex) {
		ok := true
		for j := range c.Diagram.Source.Objects {
			lambda := to.Components[j]
			mu := c.Components[j]

			composed, err := C.Compose(u, lambda)
			if err != nil {
				ok = false
				break
			}

			if composed != mu {
				ok = false
				break
			}
		}

		if ok {
			result = append(result, u)
		}
	}

	return result
}

func (c *Cone) IsLimit() bool {
	for _, other := range EnumerateCones(c.Diagram) {
		morphisms := other.Morphism(c)

		// Check existence.
		if len(morphisms) == 0 {
			return false
		}

		// Check uniqueness.
		if len(morphisms) != 1 {
			return false
		}
	}

	return true
}

func FindLimits(Diagram *Functor) []*Cone {
	var result []*Cone

	for _, cone := range EnumerateCones(Diagram) {
		if cone.IsLimit() {
			result = append(result, cone)
		}
	}

	return result
}
