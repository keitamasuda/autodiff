/* -*- mode: go; -*-
 *
 * Copyright (C) 2015-2020 Philipp Benner
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
package autodiff
/* -------------------------------------------------------------------------- */
import "fmt"
import "bytes"
import "encoding/json"
import "sort"
/* vector type declaration
 * -------------------------------------------------------------------------- */
type SparseConstInt32Vector struct {
  values []int32
  indices []int
  idxmap map[int]int
  n int
}
/* constructors
 * -------------------------------------------------------------------------- */
func UnsafeSparseConstInt32Vector(indices []int, values []int32, n int) SparseConstInt32Vector {
  r := nilSparseConstInt32Vector(n)
  r.indices = indices
  r.values = values
  return r
}
// Allocate a new vector. Scalars are set to the given values.
func NewSparseConstInt32Vector(indices []int, values []int32, n int) SparseConstInt32Vector {
  if len(indices) != len(values) {
    panic("invalid number of indices")
  }
  sort.Sort(sortIntConstInt32{indices, values})
  r := nilSparseConstInt32Vector(n)
  r.indices = indices[0:0]
  r.values = make([]int32, 0, len(values))
  for i, k := range indices {
    if k >= n {
      panic("index larger than vector dimension")
    }
    if values[i] != 0.0 {
      r.values = append(r.values, values[i])
      r.indices = append(r.indices, k)
    }
  }
  return r
}
func nilSparseConstInt32Vector(n int) SparseConstInt32Vector {
  r := SparseConstInt32Vector{}
  r.n = n
  // create map here so that no pointer receivers are needed
  r.idxmap = make(map[int]int)
  return r
}
// Convert vector type.
func AsSparseConstInt32Vector(v ConstVector) SparseConstInt32Vector {
  switch v_ := v.(type) {
  case SparseConstInt32Vector:
    return v_
  }
  indices := []int{}
  values := []int32{}
  n := v.Dim()
  for it := v.ConstIterator(); it.Ok(); it.Next() {
    indices = append(indices, it.Index())
    values = append(values, it.GetConst().GetInt32())
  }
  return NewSparseConstInt32Vector(indices, values, n)
}
/* cloning
 * -------------------------------------------------------------------------- */
func (obj SparseConstInt32Vector) Clone() SparseConstInt32Vector {
  r := nilSparseConstInt32Vector(obj.n)
  r.indices = obj.indices
  r.values = obj.values
  return r
}
/* methods specific to this type
 * -------------------------------------------------------------------------- */
func (obj SparseConstInt32Vector) GetSparseIndices() []int {
  return obj.indices
}
func (obj SparseConstInt32Vector) GetSparseValues() []int32 {
  return obj.values
}
func (obj SparseConstInt32Vector) First() (int, ConstInt32) {
  return obj.indices[0], ConstInt32(obj.values[0])
}
func (obj SparseConstInt32Vector) Last() (int, ConstInt32) {
  i := len(obj.indices) - 1
  return obj.indices[i], ConstInt32(obj.values[i])
}
func (obj SparseConstInt32Vector) CreateIndex() {
  if len(obj.idxmap) == 0 {
    for i, k := range obj.indices {
      obj.idxmap[k] = i
    }
  }
}
/* const interface
 * -------------------------------------------------------------------------- */
func (obj SparseConstInt32Vector) CloneConstVector() ConstVector {
  return obj.Clone()
}
func (obj SparseConstInt32Vector) Dim() int {
  return obj.n
}
func (obj SparseConstInt32Vector) Int8At(i int) int8 {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return int8(obj.values[k])
  } else {
    return int8(0)
  }
}
func (obj SparseConstInt32Vector) Int16At(i int) int16 {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return int16(obj.values[k])
  } else {
    return int16(0)
  }
}
func (obj SparseConstInt32Vector) Int32At(i int) int32 {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return int32(obj.values[k])
  } else {
    return int32(0)
  }
}
func (obj SparseConstInt32Vector) Int64At(i int) int64 {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return int64(obj.values[k])
  } else {
    return int64(0)
  }
}
func (obj SparseConstInt32Vector) IntAt(i int) int {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return int(obj.values[k])
  } else {
    return int(0)
  }
}
func (obj SparseConstInt32Vector) Float32At(i int) float32 {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return float32(obj.values[k])
  } else {
    return float32(0)
  }
}
func (obj SparseConstInt32Vector) Float64At(i int) float64 {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return float64(obj.values[k])
  } else {
    return float64(0)
  }
}
func (obj SparseConstInt32Vector) ConstAt(i int) ConstScalar {
  if len(obj.idxmap) == 0 {
    obj.CreateIndex()
  }
  if k, ok := obj.idxmap[i]; ok {
    return ConstFloat64(obj.values[k])
  } else {
    return ConstFloat64(0.0)
  }
}
func (obj SparseConstInt32Vector) ConstSlice(i, j int) ConstVector {
  if i == 0 {
    k1 := 0
    k2 := sort.SearchInts(obj.indices, j)
    r := nilSparseConstInt32Vector(j)
    r.values = obj.values [k1:k2]
    r.indices = obj.indices[k1:k2]
    return r
  } else {
    k1 := sort.SearchInts(obj.indices, i)
    k2 := sort.SearchInts(obj.indices, j)
    r := nilSparseConstInt32Vector(j-i)
    r.values = obj.values[k1:k2]
    r.indices = make([]int, k2-k1)
    for k := k1; k < k2; k++ {
      r.indices[k-k1] = obj.indices[k] - i
    }
    return r
  }
}
func (obj SparseConstInt32Vector) AsConstMatrix(n, m int) ConstMatrix {
  panic("not implemented")
}
/* imlement ScalarContainer
 * -------------------------------------------------------------------------- */
func (obj SparseConstInt32Vector) Reduce(f func(Scalar, ConstScalar) Scalar, r Scalar) Scalar {
  for _, v := range obj.values {
    r = f(r, ConstInt32(v))
  }
  return r
}
func (obj SparseConstInt32Vector) ElementType() ScalarType {
  return ConstInt32Type
}
/* iterator methods
 * -------------------------------------------------------------------------- */
func (obj SparseConstInt32Vector) ConstIterator() VectorConstIterator {
  return obj.ITERATOR()
}
func (obj SparseConstInt32Vector) ConstIteratorFrom(i int) VectorConstIterator {
  return obj.ITERATOR_FROM(i)
}
func (obj SparseConstInt32Vector) ConstJointIterator(b ConstVector) VectorConstJointIterator {
  return obj.JOINT_ITERATOR(b)
}
func (obj SparseConstInt32Vector) ITERATOR() *SparseConstInt32VectorIterator {
  r := SparseConstInt32VectorIterator{0, obj}
  return &r
}
func (obj SparseConstInt32Vector) ITERATOR_FROM(i int) *SparseConstInt32VectorIterator {
  k := 0
  for j, idx := range obj.indices {
    if idx >= i {
      k = j
      break
    }
  }
  r := SparseConstInt32VectorIterator{k, obj}
  return &r
}
func (obj SparseConstInt32Vector) JOINT_ITERATOR(b ConstVector) *SparseConstInt32VectorJointIterator {
  r := SparseConstInt32VectorJointIterator{}
  r.it1 = obj.ITERATOR()
  r.it2 = b.ConstIterator()
  r.idx = -1
  r.Next()
  return &r
}
/* type conversion
 * -------------------------------------------------------------------------- */
func (obj SparseConstInt32Vector) String() string {
  var buffer bytes.Buffer
  buffer.WriteString(fmt.Sprintf("%d:[", obj.n))
  first := true
  for it := obj.ConstIterator(); it.Ok(); it.Next() {
    if !first {
      buffer.WriteString(", ")
    } else {
      first = false
    }
    buffer.WriteString(fmt.Sprintf("%d:%v", it.Index(), it.GetConst()))
  }
  buffer.WriteString("]")
  return buffer.String()
}
func (obj SparseConstInt32Vector) Table() string {
  var buffer bytes.Buffer
  first := true
  for it := obj.ConstIterator(); it.Ok(); it.Next() {
    if !first {
      buffer.WriteString(" ")
    } else {
      first = false
    }
    buffer.WriteString(fmt.Sprintf("%d:%v", it.Index(), it.GetConst()))
  }
  if len(obj.indices) > 0 {
    if i := obj.indices[len(obj.indices)-1]; i != obj.n-1 {
      buffer.WriteString(fmt.Sprintf(" %d:%v", i, int32(0)))
    }
  }
  return buffer.String()
}
/* json
 * -------------------------------------------------------------------------- */
func (obj SparseConstInt32Vector) MarshalJSON() ([]byte, error) {
  k := []int{}
  v := []int32{}
  r := struct{
    Index []int
    Value []int32
    Length int}{}
  for it := obj.ConstIterator(); it.Ok(); it.Next() {
    k = append(k, it.Index())
    v = append(v, it.GetConst().GetInt32())
  }
  r.Index = k
  r.Value = v
  r.Length = obj.n
  return json.MarshalIndent(r, "", "  ")
}
/* iterator
 * -------------------------------------------------------------------------- */
type SparseConstInt32VectorIterator struct {
  i int
  v SparseConstInt32Vector
}
func (obj *SparseConstInt32VectorIterator) GetConst() ConstScalar {
  return obj.GET()
}
func (obj *SparseConstInt32VectorIterator) GET() ConstInt32 {
  return ConstInt32(obj.v.values[obj.i])
}
func (obj *SparseConstInt32VectorIterator) Ok() bool {
  return obj.i < len(obj.v.indices)
}
func (obj *SparseConstInt32VectorIterator) Index() int {
  return obj.v.indices[obj.i]
}
func (obj *SparseConstInt32VectorIterator) Next() {
  obj.i += 1
}
func (obj *SparseConstInt32VectorIterator) Clone() *SparseConstInt32VectorIterator {
  return &SparseConstInt32VectorIterator{obj.i, obj.v}
}
func (obj *SparseConstInt32VectorIterator) CloneConstIterator() VectorConstIterator {
  return &SparseConstInt32VectorIterator{obj.i, obj.v}
}
/* joint iterator
 * -------------------------------------------------------------------------- */
type SparseConstInt32VectorJointIterator struct {
  it1 *SparseConstInt32VectorIterator
  it2 VectorConstIterator
  idx int
  s1 ConstInt32
  s2 ConstScalar
}
func (obj *SparseConstInt32VectorJointIterator) Index() int {
  return obj.idx
}
func (obj *SparseConstInt32VectorJointIterator) Ok() bool {
  return !(obj.s1.GetInt32() == int32(0)) ||
         !(obj.s2.GetInt32() == int32(0))
}
func (obj *SparseConstInt32VectorJointIterator) Next() {
  ok1 := obj.it1.Ok()
  ok2 := obj.it2.Ok()
  obj.s1 = ConstInt32(0)
  obj.s2 = ConstInt32(0)
  if ok1 {
    obj.idx = obj.it1.Index()
    obj.s1 = obj.it1.GET()
  }
  if ok2 {
    switch {
    case obj.idx > obj.it2.Index() || !ok1:
      obj.idx = obj.it2.Index()
      obj.s1 = ConstInt32(0)
      obj.s2 = obj.it2.GetConst()
    case obj.idx == obj.it2.Index():
      obj.s2 = obj.it2.GetConst()
    }
  }
  if obj.s1 != ConstInt32(0) {
    obj.it1.Next()
  }
  if obj.s2 != ConstInt32(0) {
    obj.it2.Next()
  } else {
    obj.s2 = ConstInt32(0.0)
  }
}
func (obj *SparseConstInt32VectorJointIterator) GetConst() (ConstScalar, ConstScalar) {
  return obj.s1, obj.s2
}
func (obj *SparseConstInt32VectorJointIterator) GET() (ConstInt32, ConstScalar) {
  return obj.s1, obj.s2
}
func (obj *SparseConstInt32VectorJointIterator) CloneConstJointIterator() VectorConstJointIterator {
  r := SparseConstInt32VectorJointIterator{}
  r.it1 = obj.it1.Clone()
  r.it2 = obj.it2.CloneConstIterator()
  r.idx = obj.idx
  r.s1 = obj.s1
  r.s2 = obj.s2
  return &r
}
/* math
 * -------------------------------------------------------------------------- */
// Test if elements in a equal elements in b.
func (a SparseConstInt32Vector) Equals(b ConstVector, epsilon float64) bool {
  if a.Dim() != b.Dim() {
    panic("Equals(): Vector dimensions do not match!")
  }
  for it := a.ConstJointIterator(b); it.Ok(); it.Next() {
    s1, s2 := it.GetConst()
    if !s1.Equals(s2, epsilon) {
      return false
    }
  }
  return true
}
