/* Copyright (C) 2015-2018 Philipp Benner
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

import "fmt"
import "bytes"
import "unsafe"

/* -------------------------------------------------------------------------- */

type ConstRealMatrix struct {
  values   []float64
  rows       int
  cols       int
  rowOffset  int
  rowMax     int
  colOffset  int
  colMax     int
  transposed bool
}

/* constructors
 * -------------------------------------------------------------------------- */

func NewConstRealMatrix(rows, cols int, values []float64) ConstRealMatrix {
  m := ConstRealMatrix{}
  m.values    = values
  m.rows      = rows
  m.cols      = cols
  m.rowOffset = 0
  m.rowMax    = rows
  m.colOffset = 0
  m.colMax    = cols
  return m
}

func NullConstRealMatrix(rows, cols int) ConstRealMatrix {
  m := ConstRealMatrix{}
  m.values    = make([]float64, rows*cols)
  m.rows      = rows
  m.cols      = cols
  m.rowOffset = 0
  m.rowMax    = rows
  m.colOffset = 0
  m.colMax    = cols
  return m
}

/* -------------------------------------------------------------------------- */

func (matrix ConstRealMatrix) index(i, j int) int {
  if i < 0 || j < 0 || i >= matrix.rows || j >= matrix.cols {
    panic(fmt.Errorf("index (%d,%d) out of bounds for matrix of dimension %dx%d", i, j, matrix.rows, matrix.cols))
  }
  if matrix.transposed {
    return (matrix.colOffset + j)*matrix.rowMax + (matrix.rowOffset + i)
  } else {
    return (matrix.rowOffset + i)*matrix.colMax + (matrix.colOffset + j)
  }
}

func (matrix ConstRealMatrix) storageLocation() uintptr {
  return uintptr(unsafe.Pointer(&matrix.values[0]))
}

func (matrix ConstRealMatrix) ElementType() ScalarType {
  return BareRealType
}

func (matrix ConstRealMatrix) Dims() (int, int) {
  return matrix.rows, matrix.cols
}

func (matrix ConstRealMatrix) ConstAt(i, j int) ConstScalar {
  return ConstReal(matrix.values[matrix.index(i, j)])
}

func (matrix ConstRealMatrix) ConstSlice(rfrom, rto, cfrom, cto int) ConstMatrix {
  m := matrix
  m.rowOffset += rfrom
  m.rows       = rto - rfrom
  m.colOffset += cfrom
  m.cols       = cto - cfrom
  return m
}

func (matrix ConstRealMatrix) ConstRow(i int) ConstVector {
  return matrix.ROW(i)
}

func (matrix ConstRealMatrix) ROW(i int) ConstRealVector {
  var v []float64
  if matrix.transposed {
    v = make([]float64, matrix.cols)
    for j := 0; j < matrix.cols; j++ {
      v[j] = matrix.values[matrix.index(i, j)]
    }
  } else {
    i = matrix.index(i, 0)
    v = matrix.values[i:i + matrix.cols]
  }
  return ConstRealVector(v)
}

func (matrix ConstRealMatrix) ConstCol(i int) ConstVector {
  return matrix.COL(i)
}

func (matrix ConstRealMatrix) COL(j int) ConstRealVector {
  var v []float64
  if matrix.transposed {
    j = matrix.index(0, j)
    v = matrix.values[j:j + matrix.rows]
  } else {
    v = make([]float64, matrix.rows)
    for i := 0; i < matrix.rows; i++ {
      v[i] = matrix.values[matrix.index(i, j)]
    }
  }
  return v
}

func (matrix ConstRealMatrix) ConstDiag() ConstVector {
  return matrix.DIAG()
}

func (matrix ConstRealMatrix) DIAG() ConstRealVector {
  n, m := matrix.Dims()
  if n != m {
    panic("Diag(): not a square matrix!")
  }
  v := make([]float64, n)
  for i := 0; i < n; i++ {
    v[i] = matrix.values[matrix.index(i, i)]
  }
  return ConstRealVector(v)
}

func (matrix ConstRealMatrix) GetValues() []float64 {
  return matrix.values
}

func (matrix ConstRealMatrix) AsConstVector() ConstVector {
  return ConstRealVector(matrix.values)
}

/* -------------------------------------------------------------------------- */

func (m ConstRealMatrix) String() string {
  var buffer bytes.Buffer
  buffer.WriteString("[")
  for i := 0; i < m.rows; i++ {
    if i != 0 {
      buffer.WriteString(",\n ")
    }
    buffer.WriteString("[")
    for j := 0; j < m.cols; j++ {
      if j != 0 {
        buffer.WriteString(", ")
      }
      buffer.WriteString(m.ConstAt(i,j).String())
    }
    buffer.WriteString("]")
  }
  buffer.WriteString("]")
  return buffer.String()
}

func (a ConstRealMatrix) Table() string {
  var buffer bytes.Buffer
  n, m := a.Dims()
  for i := 0; i < n; i++ {
    if i != 0 {
      buffer.WriteString("\n")
    }
    for j := 0; j < m; j++ {
      if j != 0 {
        buffer.WriteString(" ")
      }
      buffer.WriteString(a.ConstAt(i,j).String())
    }
  }
  return buffer.String()
}

/* math
 * -------------------------------------------------------------------------- */

func (a ConstRealMatrix) Equals(b ConstMatrix, epsilon float64) bool {
  n1, m1 := a.Dims()
  n2, m2 := b.Dims()
  if n1 != n2 || m1 != m2 {
    panic("MEqual(): matrix dimensions do not match!")
  }
  for i := 0; i < n1; i++ {
    for j := 0; j < m1; j++ {
      if !a.ConstAt(i, j).Equals(b.ConstAt(i, j), epsilon) {
        return false
      }
    }
  }
  return true
}
