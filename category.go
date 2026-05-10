package kitty

import (
	"errors"
	"fmt"
	"slices"
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

func (m *Morphism) Inverse() *Morphism {
	var inverseID MorphismID
	if strings.HasSuffix(string(m.ID), "^{-1}") {
		inverseID = MorphismID(strings.TrimSuffix(string(m.ID), "^{-1}"))
	} else {
		inverseID = m.ID + "^{-1}"
	}

	return &Morphism{
		ID:          inverseID,
		Source:      m.Destination,
		Destination: m.Source,
	}
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
		if strings.Contains(string(m.ID), "_") {
			return nil, errors.New("Morphism names with '_' cannot be used")
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
					// Avoid the infinite morphism loop.
					// e.g. m1.ID = "g_f", m2.ID = "f"
					m1IDTokens := strings.Split(string(m1.ID), "_")
					m2IDTokens := strings.Split(string(m2.ID), "_")
					if containsAsSubsequence(m1IDTokens, m2IDTokens) ||
						containsAsSubsequence(m2IDTokens, m1IDTokens) {
						continue
					}
					// Vanish the product of the pair of inverse morphisms.
					// e.g. m1.ID = fg^{-1}", m2.ID = hgf^{-1} => m1.ID = id, m2.ID = h"
					for len(m1IDTokens) != 0 && len(m2IDTokens) != 0 {
						m1FirstMorphism := c.Morphisms[MorphismID(m1IDTokens[0])]
						m2LastMorphism := c.Morphisms[MorphismID(m2IDTokens[len(m2IDTokens)-1])]
						if m1FirstMorphism.Inverse().ID == m2LastMorphism.ID {
							m1IDTokens = m1IDTokens[1:]
							m2IDTokens = m2IDTokens[:len(m2IDTokens)-1]
							continue
						}
						break
					}
					var composedID MorphismID
					switch {
					case len(m1IDTokens) == 0 && len(m2IDTokens) == 0:
						composedID = m1.Source.GetIdentityID()
					case len(m1IDTokens) == 0:
						composedID = MorphismID(strings.Join(m2IDTokens, "_"))
					case len(m2IDTokens) == 0:
						composedID = MorphismID(strings.Join(m1IDTokens, "_"))
					default:
						composedID = MorphismID(fmt.Sprintf("%s_%s",
							strings.Join(m2IDTokens, "_"), strings.Join(m1IDTokens, "_")))
					}
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
	return errors.New("composition construction stuck detected")
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

func containsAsSubsequence(tokens, subTokens []string) bool {
	if len(tokens) == 0 {
		return len(subTokens) == 0
	}
	if len(subTokens) == 0 {
		return true
	}
	i := slices.Index(tokens, subTokens[0])
	if i == -1 {
		return false
	}
	partialTokens := tokens[i:min(len(tokens), i+len(subTokens))]
	return slices.Compare(partialTokens, subTokens) == 0
}

func (c *Category) validate() error {
	// Type check
	for key, res := range c.composeTable {
		f, ok := c.Morphisms[key[0]]
		if !ok {
			return fmt.Errorf("invalid morphism %s found in the composition rule",
				key[0])
		}
		g, ok := c.Morphisms[key[1]]
		if !ok {
			return fmt.Errorf("invalid morphism %s found in the composition rule",
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
					left, err := c.Compose(gf, h)
					if err != nil {
						continue
					}
					right, err := c.Compose(f, hg)
					if err != nil {
						continue
					}
					if left != right {
						return fmt.Errorf("associativity failed: left=%s, right=%s", left, right)
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
