package kitty

import (
	"errors"
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
	Objects   []Object
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
	objects []Object, morphisms []*Morphism, compose map[[2]MorphismID]MorphismID,
) (*Category, error) {
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

	for _, o := range c.Objects {
		id := o.GetIdentityID()
		m := &Morphism{
			ID:          id,
			Source:      o,
			Destination: o,
		}
		c.Morphisms[m.ID] = m
	}

	for _, o := range c.Objects {
		for _, m := range c.Morphisms {
			if m.Source == o {
				key := [2]MorphismID{o.GetIdentityID(), m.ID}
				c.composeTable[key] = m.ID
			} else if m.Destination == o {
				key := [2]MorphismID{m.ID, o.GetIdentityID()}
				c.composeTable[key] = m.ID
			}
		}
	}

	processedCompose := c.composeTable
	for k, v := range compose {
		if v == Identity {
			// Since k[1]◦k[0] is identity, its object should be the destination of k[1].
			processedCompose[k] = c.Morphisms[k[1]].Destination.GetIdentityID()
		}
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Category) Compose(f, g MorphismID) (MorphismID, error) {
	key := [2]MorphismID{f, g}
	res, ok := c.composeTable[key]
	if !ok {
		return "", errors.New("composition not defined")
	}
	return res, nil
}

func (c *Category) validate() error {
	// Type check
	for key, res := range c.composeTable {
		f := c.Morphisms[key[0]]
		g := c.Morphisms[key[1]]

		if f.Destination != g.Source {
			return errors.New("invalid composition domain")
		}

		r := c.Morphisms[res]
		if r.Source != f.Source || r.Destination != g.Destination {
			return errors.New("composition result mismatch")
		}
	}

	// Associative law
	for fID := range c.Morphisms {
		for gID := range c.Morphisms {
			for hID := range c.Morphisms {

				f, ok1 := c.Compose(fID, gID)
				g, ok2 := c.Compose(gID, hID)

				if ok1 == nil && ok2 == nil {
					left, err1 := c.Compose(f, hID)
					right, err2 := c.Compose(fID, g)

					if err1 == nil && err2 == nil && left != right {
						return errors.New("associativity failed")
					}
				}
			}
		}
	}

	return nil
}
