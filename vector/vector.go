package vector

type Vector struct {
	s []interface{}
}

func New() *Vector {
	return &Vector{}
}

func NewCap(n int) *Vector {
	v := &Vector{}
	v.s = make([]interface{}, 0, n)
	return v
}

func (v *Vector) Len() int {
	return len(v.s)
}

func (v *Vector) Cap() int {
	return cap(v.s)
}

func (v *Vector) Append(data ...interface{}) {
	v.s = append(v.s, data...)
}

func (v *Vector) AppendVec(o *Vector) {
	v.Append(o.s...)
}

func (v *Vector) Clone() *Vector {
	t := NewCap(v.Len())
	t.Append(v.s...)
	return t
}

func (v *Vector) Copy(src *Vector) {
	v.s = make([]interface{}, len(src.s))
	copy(v.s[:], src.s[:])
}

//Delete delete value at position i
//if i is out of range, Delete will panic
func (v *Vector) Delete(i int) {
	copy(v.s[i:], v.s[i+1:])
	v.s[len(v.s)-1] = nil
	v.s = v.s[:len(v.s)-1]
}

func (v *Vector) Insert(i int, d interface{}) {
	v.s = append(v.s, nil)
	copy(v.s[i+1:], v.s[i:])
	v.s[i] = d
}
func (v *Vector) InsertVariant(i int, d ...interface{}) {
	v.s = append(v.s[0:i], append(d, v.s[i:]...)...)
}
func (v *Vector) InsertVector(i int, d *Vector) {
	v.InsertVariant(i, d.s...)
}

//At will panic while i is out of range
func (v *Vector) At(i int) interface{} {
	return v.s[i]
}

//Extend extend j space at tail
func (v *Vector) Extend(j int) {
	v.s = append(v.s, make([]interface{}, j)...)
}

//ExtendAt extend j space after position i.
func (v *Vector) ExtendAt(i, j int) {
	v.s = append(v.s[:i], append(make([]interface{}, j), v.s[i:]...)...)
}

func (v *Vector) Pop() (r interface{}) {
	r = v.s[len(v.s)-1]
	v.Delete(len(v.s) - 1)
	return r
}

func (v *Vector) Push(d interface{}) {
	v.Append(d)
}

func (v *Vector) PopFront() (r interface{}) {
	r = v.s[0]
	v.Delete(0)
	return r
}

func (v *Vector) PushFront(d interface{}) {
	v.Insert(0, d)
}

func (v *Vector) Reverse() {
	for left, right := 0, len(v.s)-1; left < right; left, right = left+1, right-1 {
		v.s[left], v.s[right] = v.s[right], v.s[left]
	}
}

func (v *Vector) Clear() {
	for i, _ := range v.s {
		v.s[i] = nil
	}
	v.s = v.s[0:0]
}
