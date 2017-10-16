/* -*- mode: go; -*-
 *
 * Copyright (C) 2015-2017 Philipp Benner
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
package autodiff
/* -------------------------------------------------------------------------- */
// Test if elements in a equal elements in b.
func (a DenseBareRealVector) Equals(b ConstVector, epsilon float64) bool {
  if a.Dim() != b.Dim() {
    panic("VEqual(): Vector dimensions do not match!")
  }
  for i := 0; i < a.Dim(); i++ {
    if !a.ConstAt(i).Equals(b.ConstAt(i), epsilon) {
      return false
    }
  }
  return true
}
/* -------------------------------------------------------------------------- */
// Element-wise addition of two vectors. The result is stored in r.
func (r DenseBareRealVector) VaddV(a, b ConstVector) Vector {
  n := len(r)
  if a.Dim() != n || b.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Add(a.ConstAt(i), b.ConstAt(i))
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Element-wise addition of a vector and a scalar. The result is stored in r.
func (r DenseBareRealVector) VaddS(a ConstVector, b ConstScalar) Vector {
  n := len(r)
  if a.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Add(a.ConstAt(i), b)
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Element-wise substraction of two vectors. The result is stored in r.
func (r DenseBareRealVector) VsubV(a, b ConstVector) Vector {
  n := len(r)
  if a.Dim() != n || b.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Sub(a.ConstAt(i), b.ConstAt(i))
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Element-wise substractor of a vector and a scalar. The result is stored in r.
func (r DenseBareRealVector) VsubS(a ConstVector, b ConstScalar) Vector {
  n := len(r)
  if a.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Sub(a.ConstAt(i), b)
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Element-wise multiplication of two vectors. The result is stored in r.
func (r DenseBareRealVector) VmulV(a, b ConstVector) Vector {
  n := len(r)
  if a.Dim() != n || b.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Mul(a.ConstAt(i), b.ConstAt(i))
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Element-wise substraction of a vector and a scalar. The result is stored in r.
func (r DenseBareRealVector) VmulS(a ConstVector, s ConstScalar) Vector {
  n := len(r)
  if a.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Mul(a.ConstAt(i), s)
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Element-wise division of two vectors. The result is stored in r.
func (r DenseBareRealVector) VdivV(a, b ConstVector) Vector {
  n := len(r)
  if a.Dim() != n || b.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Div(a.ConstAt(i), b.ConstAt(i))
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Element-wise division of a vector and a scalar. The result is stored in r.
func (r DenseBareRealVector) VdivS(a ConstVector, s ConstScalar) Vector {
  n := len(r)
  if a.Dim() != n {
    panic("vector dimensions do not match")
  }
  for i := 0; i < a.Dim(); i++ {
    r[i].Div(a.ConstAt(i), s)
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Matrix vector product of a and b. The result is stored in r.
func (r DenseBareRealVector) MdotV(a Matrix, b ConstVector) Vector {
  n, m := a.Dims()
  if r.Dim() != n || b.Dim() != m {
    panic("matrix/vector dimensions do not match!")
  }
  if n == 0 || m == 0 {
    return r
  }
  if &r[0] == b.ConstAt(0) {
    panic("result and argument must be different vectors")
  }
  t := NullScalar(a.ElementType())
  for i := 0; i < n; i++ {
    r[i].Reset()
    for j := 0; j < m; j++ {
      t.Mul(a.At(i, j), b.ConstAt(j))
      r[i].Add(r.At(i), t)
    }
  }
  return r
}
/* -------------------------------------------------------------------------- */
// Vector matrix product of a and b. The result is stored in r.
func (r DenseBareRealVector) VdotM(a ConstVector, b Matrix) Vector {
  n, m := b.Dims()
  if r.Dim() != m || a.Dim() != n {
    panic("matrix/vector dimensions do not match!")
  }
  if n == 0 || m == 0 {
    return r
  }
  if r.At(0) == a.ConstAt(0) {
    panic("result and argument must be different vectors")
  }
  t := NullScalar(a.ElementType())
  for i := 0; i < m; i++ {
    r[i].Reset()
    for j := 0; j < n; j++ {
      t.Mul(a.ConstAt(j), b.At(j, i))
      r[i].Add(r.At(i), t)
    }
  }
  return r
}
