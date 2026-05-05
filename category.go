package kitty

import (
	"errors"
	"fmt"
	"strings"
)

const Identity = "identity"

type Object string
type MorphismID string

type Morphism struct {
	ID          MorphismID
	Source      Object
	Destination Object
}

type Category struct {
	Objects   map[Object]struct{}
	Morphisms map[MorphismID]*Morphism

	// key: f,g -> g∘f
	composeTable map[[2]MorphismID]MorphismID
}

func (o Object) GetIdentityID() MorphismID {
	return MorphismID("_id_" + o)
}

func (m *Morphism) IsIdentity() bool {
	return strings.HasPrefix(string(m.ID), "_id") &&
		m.Source == m.Destination
}

func NewCategory(
	objs []Object, morphisms []*Morphism, compose map[[2]MorphismID]MorphismID,
) (*Category, error) {
	objects := make(map[Object]struct{})
	for _, o := range objs {
		objects[o] = struct{}{}
	}
	c := &Category{
		Objects:      objects,
		Morphisms:    make(map[MorphismID]*Morphism),
		composeTable: compose,
	}

	for _, m := range morphisms {
		if strings.HasPrefix(string(m.ID), "_") {
			return nil, errors.New("Morphism name should not start with '_'.")
		}
		c.Morphisms[m.ID] = m
	}

	for o := range c.Objects {
		id := o.GetIdentityID()
		m := &Morphism{
			ID:          id,
			Source:      o,
			Destination: o,
		}
		c.Morphisms[m.ID] = m
	}

	processedCompose := c.composeTable
	for k, v := range compose {
		if v == Identity {
			// Since k[1]◦k[0] is identity, its object should be the destination of k[1].
			processedCompose[k] = c.Morphisms[k[1]].Destination.GetIdentityID()
		}
	}

	for k, v := range c.composeTable {
		for _, w := range []MorphismID{k[0], k[1], v} {
			if _, ok := c.Morphisms[w]; !ok {
				return nil, fmt.Errorf("unknown morphism found in the composition definition: %s", w)
			}
		}
	}

	err := c.constructComposition()
	if err != nil {
		return nil, err
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Category) constructComposition() error {
	initialMorphismCount := len(c.Morphisms)
	previousCTableCount := len(c.composeTable)
	for len(c.Morphisms) < initialMorphismCount*initialMorphismCount {
		toBeAddedMorphisms := make([]*Morphism, 0)
		for _, m1 := range c.Morphisms {
			for _, m2 := range c.Morphisms {
				if m1.Destination != m2.Source {
					continue
				}

				key := [2]MorphismID{m1.ID, m2.ID}
				if _, ok := c.composeTable[key]; ok {
					continue
				}
				switch {
				case m1.IsIdentity():
					// This case contains (m1.ID, m2.ID) = (identity, identity).
					composedID := m2.ID
					c.composeTable[key] = composedID
				case m2.IsIdentity():
					composedID := m1.ID
					c.composeTable[key] = composedID
				default:
					// e.g. m1.ID = f, m2.ID = g => gf
					composedID := m2.ID + m1.ID
					c.composeTable[key] = composedID
					toBeAddedMorphisms = append(toBeAddedMorphisms,
						&Morphism{
							ID:          composedID,
							Source:      m1.Source,
							Destination: m2.Destination,
						},
					)
				}
			}
		}
		for _, m := range toBeAddedMorphisms {
			c.Morphisms[m.ID] = m
		}
		if previousCTableCount == len(c.composeTable) {
			return nil
		}
		previousCTableCount = len(c.composeTable)
	}
	return errors.New("composition construction stuck detected.")
}

func (c *Category) Compose(f, g MorphismID) (MorphismID, error) {
	key := [2]MorphismID{f, g}
	res, ok := c.composeTable[key]
	if !ok {
		return "", fmt.Errorf("composition not defined for %s and %s",
			f, g)
	}
	return res, nil
}

func (c *Category) validate() error {
	// Type check
	for key, res := range c.composeTable {
		f, ok := c.Morphisms[key[0]]
		if !ok {
			return fmt.Errorf("invalid morphism %s found in the composition rule.",
				key[0])
		}
		g, ok := c.Morphisms[key[1]]
		if !ok {
			return fmt.Errorf("invalid morphism %s found in the composition rule.",
				key[1])
		}

		if f.Destination != g.Source {
			return errors.New("invalid composition domain")
		}

		r := c.Morphisms[res]
		if r.Source != f.Source || r.Destination != g.Destination {
			return errors.New("composition result mismatch")
		}
	}

	// Associative law
	for f := range c.Morphisms {
		for g := range c.Morphisms {
			for h := range c.Morphisms {
				gf, ok1 := c.Compose(f, g)
				hg, ok2 := c.Compose(g, h)

				if ok1 == nil && ok2 == nil {
					left, err1 := c.Compose(gf, h)
					right, err2 := c.Compose(f, hg)

					if err1 == nil && err2 == nil && left != right {
						return errors.New("associativity failed")
					}
				}
			}
		}
	}

	// Invalid source or destination
	for _, m := range c.Morphisms {
		if _, ok := c.Objects[m.Source]; !ok {
			return fmt.Errorf("%s has an invalid source: %s", m.ID, m.Source)
		}
		if _, ok := c.Objects[m.Destination]; !ok {
			return fmt.Errorf("%s has an invalid destination: %s", m.ID, m.Destination)
		}
	}

	return nil
}
