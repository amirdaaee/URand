package supervisor

import rng "URand/RNG"

func NewNS(name string, len uint) (*NameSpace, error) {
	ns := NameSpace{Name: name, Len: len, generator: rng.GetGenerator(len)}
	if err := ns.checkFill(); err != nil {
		return nil, err
	}
	return &ns, nil
}
