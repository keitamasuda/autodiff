/* Copyright (C) 2015 Philipp Benner
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

package autodiff

/* -------------------------------------------------------------------------- */

//import "fmt"

/* -------------------------------------------------------------------------- */

// True if matrix a equals b.
func Mequal(a, b Matrix) bool {
  n1, m1 := a.Dims()
  n2, m2 := b.Dims()
  if n1 != n2 || m1 != m2 {
    panic("MEqual(): matrix dimensions do not match!")
  }
  v1 := a.ToVector()
  v2 := b.ToVector()
  for i := 0; i < v1.Dim(); i++ {
    if !v1.At(i).Equals(v2.At(i)) {
      return false
    }
  }
  return true
}

/* -------------------------------------------------------------------------- */

// Element-wise addition of two matrices. The result is stored in r.
func (r *DenseMatrix) MaddM(a, b Matrix) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  n2, m2 := b.Dims()
  if n1 != n || m1 != m || n2 != n || m2 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Add(a.At(i, j), b.At(i, j))
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Add scalar b to all elements of a. The result is stored in r.
func (r *DenseMatrix) MaddS(a Matrix, b Scalar) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  if n1 != n || m1 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Add(a.At(i, j), b)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Element-wise substraction of two matrices. The result is stored in r.
func (r *DenseMatrix) MsubM(a, b Matrix) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  n2, m2 := b.Dims()
  if n1 != n || m1 != m || n2 != n || m2 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Sub(a.At(i, j), b.At(i, j))
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Substract b from all elements of a. The result is stored in r.
func (r *DenseMatrix) MsubS(a Matrix, b Scalar) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  if n1 != n || m1 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Sub(a.At(i, j), b)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Element-wise multiplication of two matrices. The result is stored in r.
func (r *DenseMatrix) MmulM(a, b Matrix) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  n2, m2 := b.Dims()
  if n1 != n || m1 != m || n2 != n || m2 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Mul(a.At(i, j), b.At(i, j))
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Multiply all elements of a with b. The result is stored in r.
func (r *DenseMatrix) MmulS(a Matrix, b Scalar) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  if n1 != n || m1 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Mul(a.At(i, j), b)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Element-wise division of two matrices. The result is stored in r.
func (r *DenseMatrix) MdivM(a, b Matrix) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  n2, m2 := b.Dims()
  if n1 != n || m1 != m || n2 != n || m2 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Div(a.At(i, j), b.At(i, j))
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Divide all elements of a by b. The result is stored in r.
func (r *DenseMatrix) MdivS(a Matrix, b Scalar) Matrix {
  n,  m  := r.Dims()
  n1, m1 := a.Dims()
  if n1 != n || m1 != m {
    panic("matrix dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Div(a.At(i, j), b)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Matrix product of a and b. The result is stored in r.
func (r *DenseMatrix) MdotM(a, b Matrix) Matrix {
  n, m := r.Dims()
  n1, m1 := a.Dims()
  n2, m2 := b.Dims()
  if n1 != n || m2 != m || n1 != m2 || m1 != n2 {
    panic("matrix dimensions do not match!")
  }
  t1 := NullScalar(a.ElementType())
  t2 := NullScalar(a.ElementType())
  if r == b {
    t3 := r.Tmp2[0:n]
    for j := 0; j < m; j++ {
      for i := 0; i < n; i++ {
        t2.Reset()
        for k := 0; k < m1; k++ {
          t1.Mul(a.At(i, k), b.At(k, j))
          t2.Add(t2, t1)
        }
        t3[i].Set(t2)
      }
      for i := 0; i < n; i++ {
        r.At(i, j).Set(t3[i])
      }
    }
  } else {
    t3 := r.Tmp2[0:m]
    for i := 0; i < n; i++ {
      for j := 0; j < m; j++ {
        t2.Reset()
        for k := 0; k < m1; k++ {
          t1.Mul(a.At(i, k), b.At(k, j))
          t2.Add(t2, t1)
        }
        t3[j].Set(t2)
      }
      for j := 0; j < m; j++ {
        r.At(i, j).Set(t3[j])
      }
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Matrix vector product of a and b. The result is stored in r.
func (r DenseVector) MdotV(a Matrix, b Vector) Vector {
  n, m := a.Dims()
  if r[0] == b.At(0) {
    panic("result and argument must be different vectors")
  }
  if r.Dim() != n || b.Dim() != m {
    panic("matrix/vector dimensions do not match!")
  }
  t := NullScalar(a.ElementType())
  for i := 0; i < n; i++ {
    r[i].Reset()
    for j := 0; j < m; j++ {
      t.Mul(a.At(i, j), b.At(j))
      r[i].Add(r[i], t)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Vector matrix product of a and b. The result is stored in r.
func (r DenseVector) VdotM(a Vector, b Matrix) Vector {
  n, m := b.Dims()
  if r[0] == a.At(0) {
    panic("result and argument must be different vectors")
  }
  if r.Dim() != m || a.Dim() != n {
    panic("matrix/vector dimensions do not match!")
  }
  t := NullScalar(a.ElementType())
  for i := 0; i < m; i++ {
    r[i].Reset()
    for j := 0; j < n; j++ {
      t.Mul(a.At(j), b.At(j, i))
      r[i].Add(r[i], t)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Outer product of two vectors. The result is stored in r.
func (r *DenseMatrix) Outer(a, b Vector) Matrix {
  n, m := r.Dims()
  if a.Dim() != n || b.Dim() != m {
    panic("matrix/vector dimensions do not match!")
  }
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).Mul(a.At(i), b.At(j))
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

// Compute the Jacobian of f at x_. The result is stored in r.
func (r *DenseMatrix) Jacobian(f func(Vector) Vector, x_ Vector) Matrix {
  n, m := r.Dims()
  x := x_.CloneVector()
  x.Variables(1)
  // compute Jacobian
  y := f(x)
  // reallocate matrix if dimensions do not match
  if r == nil || x.Dim() != m || y.Dim() != n {
     n = y.Dim()
     m = x.Dim()
    *r = *NullDenseMatrix(x_.ElementType(), n, m)
  }
  // copy derivatives
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).SetValue(y.At(i).GetDerivative(j))
    }
  }
  return r
}

// Compute the Hessian of f at x_. The result is stored in r.
func (r *DenseMatrix) Hessian(f func(Vector) Scalar, x_ Vector) Matrix {
  n, m := r.Dims()
  // reallocate matrix if dimensions do not match
  if r == nil || x_.Dim() != n || n != m {
     n = x_.Dim()
     m = x_.Dim()
    *r = *NullDenseMatrix(x_.ElementType(), n, m)
  }
  x := x_.CloneVector()
  x.Variables(2)
  // evaluate function
  y := f(x)
  // copy second derivatives
  for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
      r.At(i, j).SetValue(y.GetHessian(i, j))
    }
  }
  return r
}
