package map_index

import "errors"

type Index map[string][]byte

func NewMapIndex() Index {
	return make(map[string][]byte)
}

func (i Index) Get(k []byte) ([]byte, error) {
	v, ok := i[string(k)]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (i Index) Set(k []byte, v []byte) error {
	i[string(k)] = v
	return nil
}

func (i Index) Delete(k []byte) error {
	delete(i, string(k))
	return nil
}

func (i Index) Batch(ks [][]byte, vs [][]byte) error {
	if len(ks) != len(vs) {
		return errors.New("not euqal batch")
	}
	for j := range ks {
		i[string(ks[j])] = vs[j]
	}
	return nil
}
