package dao

type Params []Param

func NewParams(len int) Params {
	return make([]Param, 0, len)
}

func (pp Params) Len() int {
	return len(pp)
}

func (pp Params) Append(param Param) Params {

	pp = append(pp, param)

	return pp
}

//

type Param map[string]interface{}

func NewParam(len int) Param {
	return make(map[string]interface{}, len)
}

func (pp Param) Len() int {
	return len(pp)
}

func (p Param) Add(key string, value interface{}) Param {

	p[key] = value

	return p
}
