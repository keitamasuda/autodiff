/* Copyright (C) 2019 Philipp Benner
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

package saga

/* -------------------------------------------------------------------------- */

import   "fmt"
import   "math"
import   "testing"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics/vectorDistribution"

/* -------------------------------------------------------------------------- */

func hook(x ConstVector, step, lambda ConstScalar, i int) bool {
  fmt.Printf("x: %v\n", x)
  fmt.Printf("s: %v\n", step)
  fmt.Printf("d: %v\n", i)
  fmt.Println()
  return false
}

/* -------------------------------------------------------------------------- */

type proximalWrapper struct {
  ProximalOperatorType
}

func (obj proximalWrapper) Eval(x DenseFloat64Vector, w DenseFloat64Vector) {
  obj.ProximalOperatorType.Eval(x, w)
  // do not regularize intercept
  x.AT(0).SET(w.AT(0))
}

/* -------------------------------------------------------------------------- */

type jitUpdateWrapper struct {
  JitUpdateType
}

func (obj jitUpdateWrapper) Update(x, y float64, k, m int) float64 {
  // do not regularize intercept
  if k == 0 {
    return x - y
  } else {
    return obj.JitUpdateType.Update(x, y, k, m)
  }
}

/* -------------------------------------------------------------------------- */

func f_dense(class []float64, x []DenseFloat64Vector) Objective1Dense {
  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  lr, _   := NewLogisticRegression(theta_0)
  r       := NewFloat64(0.0)
  f := func(i int, theta DenseFloat64Vector) (float64, float64, DenseFloat64Vector, error) {
    y := 0.0
    w := 0.0
    if i >= len(x) {
      return y, w, nil, fmt.Errorf("index out of bounds")
    }
    if err := lr.SetParameters(theta); err != nil {
      return y, w, nil, err
    }
    if err := lr.LogPdf(r, x[i]); err != nil {
      return y, w, nil, err
    }
    if math.IsNaN(r.GetFloat64()) {
      return y, w, nil, fmt.Errorf("NaN value detected")
    }
    y = r.GetFloat64()
    w = math.Exp(r.GetFloat64()) - class[i]
    return y, w, x[i], nil
  }
  return f
}

func f_sparse(class []float64, x []SparseConstFloat64Vector) Objective1Sparse {
  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  lr, _   := NewLogisticRegression(theta_0)
  r       := NewFloat64(0.0)
  f := func(i int, theta DenseFloat64Vector) (float64, float64, SparseConstFloat64Vector, error) {
    y := 0.0
    w := 0.0
    if i >= len(x) {
      return y, w, x[i], fmt.Errorf("index out of bounds")
    }
    if err := lr.SetParameters(theta); err != nil {
      return y, w, x[i], err
    }
    if err := lr.LogPdf(r, x[i]); err != nil {
      return y, w, x[i], err
    }
    if math.IsNaN(r.GetFloat64()) {
    return y, w, x[i], fmt.Errorf("NaN value detected")
    }
    y = r.GetFloat64()
    w = math.Exp(r.GetFloat64()) - class[i]
    return y, w, x[i], nil
  }
  return f
}

/* -------------------------------------------------------------------------- */

func Test0(test *testing.T) {
  x := NewDenseFloat64Vector([]float64{1, 1})
  r := NewReal64(0.0)
  { // test positive class
    theta_0 := NewDenseReal64Vector([]float64{-1, 0.0, 0.0})
    theta_0.Variables(1)
    lr, _   := NewLogisticRegression(theta_0)

    if err := lr.ClassLogPdf(r, x, true); err != nil {
      test.Error(err); return
    }
    if math.Abs(r.GetFloat64() - -1.313262e+00) > 1e-4 {
      test.Error("test failed")
    }
    if math.Abs(r.GetDerivative(1) - 0.7310585786300048) > 1e-4 {
      test.Error("test failed")
    }
  }
  { // test negative class
    theta_0 := NewDenseReal64Vector([]float64{-1, 0.0, 0.0})
    theta_0.Variables(1)
    lr, _   := NewLogisticRegression(theta_0)

    if err := lr.ClassLogPdf(r, x, false); err != nil {
      test.Error(err); return
    }
    if math.Abs(r.GetFloat64() - -0.313262) > 1e-4 {
      test.Error("test failed")
    }
    if math.Abs(r.GetDerivative(1) - -0.268941) > 1e-4 {
      test.Error("test failed")
    }
  }
}

func Test1(test *testing.T) {

  // data
  cellSize  := []float64{
    1, 4, 1, 8, 1, 10, 1, 1, 1, 2, 1, 1, 3, 1, 7, 4, 1, 1, 7, 1}
  cellShape := []float64{
    1, 4, 1, 8, 1, 10, 1, 2, 1, 1, 1, 1, 3, 1, 5, 6, 1, 1, 7, 1}
  class := []float64{
    0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0}
  // x
  x := make([]DenseFloat64Vector, len(cellSize))
  for i := 0; i < len(cellSize); i++ {
    x[i] = NewDenseFloat64Vector([]float64{cellSize[i], cellShape[i]})
  }

  theta_0 := NewDenseReal64Vector([]float64{-1, 0.0, 0.0})
  lr, _   := NewLogisticRegression(theta_0)

  f := func(i int, theta Vector, r MagicScalar) error {
    if i >= len(cellSize) {
      return fmt.Errorf("index out of bounds")
    }
    if err := lr.SetParameters(theta); err != nil {
      return err
    }
    if err := lr.ClassLogPdf(r, x[i], class[i] == 1); err != nil {
      return err
    }
    if math.IsNaN(r.GetFloat64()) {
      return fmt.Errorf("NaN value detected")
    }
    // minimize negative log likelihood
    r.Neg(r)
    return nil
  }
  z := DenseFloat64Vector([]float64{-3.549076e+00, 1.840901e-01, 5.067003e-01})
  t := NullFloat64()

  if r, _, err := Run(WrapperDense(f), len(cellSize), theta_0, Hook{}, Gamma{1.0/20}, Epsilon{1e-8}, L1Regularization{0.0}, L2Regularization{0.0}); err != nil {
    test.Error(err)
  } else {
    if t.Vnorm(r.VsubV(r, z)); t.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
}

func Test2(test *testing.T) {

  // data
  cellSize  := []float64{
    1, 4, 1, 8, 1, 10, 1, 1, 1, 2, 1, 1, 3, 1, 7, 4, 1, 1, 7, 1}
  cellShape := []float64{
    1, 4, 1, 8, 1, 10, 1, 2, 1, 1, 1, 1, 3, 1, 5, 6, 1, 1, 7, 1}
  class := []float64{
    0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0}
  // x
  x := make([]DenseFloat64Vector, len(cellSize))
  for i := 0; i < len(cellSize); i++ {
    x[i] = DenseFloat64Vector([]float64{1.0, cellSize[i], cellShape[i]})
  }

  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  z := DenseFloat64Vector([]float64{-3.549076e+00, 1.840901e-01, 5.067003e-01})
  t := NullFloat64()

  if r, _, err := Run(Objective1Dense(f_dense(class, x)), len(cellSize), theta_0, Hook{}, Gamma{1.0/20}, Epsilon{1e-8}, L1Regularization{0.0}, L2Regularization{0.0}); err != nil {
    test.Error(err)
  } else {
    if t.Vnorm(r.VsubV(r, z)); t.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
}

func Test3(test *testing.T) {

  // data
  cellSize  := []float64{
    1, 4, 1, 8, 1, 10, 1, 1, 1, 2, 1, 1, 3, 1, 7, 4, 1, 1, 7, 1}
  cellShape := []float64{
    1, 4, 1, 8, 1, 10, 1, 2, 1, 1, 1, 1, 3, 1, 5, 6, 1, 1, 7, 1}
  class := []float64{
    0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0}
  // x
  x := make([]SparseConstFloat64Vector, len(cellSize))
  for i := 0; i < len(cellSize); i++ {
    x[i] = NewSparseConstFloat64Vector([]int{0, 1, 2}, []float64{1.0, cellSize[i], cellShape[i]}, 3)
  }

  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  z := DenseFloat64Vector([]float64{-3.549076e+00, 1.840901e-01, 5.067003e-01})
  t := NullFloat64()

  if r, _, err := Run(Objective1Sparse(f_sparse(class, x)), len(cellSize), theta_0, Hook{}, Gamma{1.0/20}, Epsilon{1e-8}, L1Regularization{0.0}, L2Regularization{0.0}); err != nil {
    test.Error(err)
  } else {
    if t.Vnorm(r.VsubV(r, z)); t.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
}

func Test4(test *testing.T) {

  // data
  cellSize  := []float64{
    1, 4, 1, 8, 1, 10, 1, 1, 1, 2, 1, 1, 3, 1, 7, 4, 1, 1, 7, 1}
  cellShape := []float64{
    1, 4, 1, 8, 1, 10, 1, 2, 1, 1, 1, 1, 3, 1, 5, 6, 1, 1, 7, 1}
  class := []float64{
    0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0}
  // x
  x := make([]SparseConstFloat64Vector, len(cellSize))
  for i := 0; i < len(cellSize); i++ {
    x[i] = NewSparseConstFloat64Vector([]int{0, 1, 2}, []float64{1.0, cellSize[i]-1.0, cellShape[i]-1.0}, 3)
  }

  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  z := DenseFloat64Vector([]float64{-2.858321e+00, 1.840900e-01, 5.067086e-01})
  t := NullFloat64()

  if r, _, err := Run(Objective1Sparse(f_sparse(class, x)), len(cellSize), theta_0, Hook{}, Gamma{1.0/20}, Epsilon{1e-8}, L1Regularization{0.0}, L2Regularization{0.0}); err != nil {
    test.Error(err)
  } else {
    if t.Vnorm(r.VsubV(r, z)); t.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
}

func Test5(test *testing.T) {

  // data
  cellSize  := []float64{
    1, 4, 1, 8, 1, 10, 1, 1, 1, 2, 1, 1, 3, 1, 7, 4, 1, 1, 7, 1}
  cellShape := []float64{
    1, 4, 1, 8, 1, 10, 1, 2, 1, 1, 1, 1, 3, 1, 5, 6, 1, 1, 7, 1}
  class := []float64{
    0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0}
  // x
  x := make([]SparseConstFloat64Vector, len(cellSize))
  for i := 0; i < len(cellSize); i++ {
    x[i] = NewSparseConstFloat64Vector([]int{0, 1, 2}, []float64{1.0, cellSize[i]-1.0, cellShape[i]-1.0}, 3)
  }

  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  z := DenseFloat64Vector([]float64{-2.76467776, 0.17584927, 0.48174453})
  t := NullFloat64()
  p := proximalWrapper{&ProximalOperatorL1{1.0/2.5}}

  if r, _, err := Run(Objective1Sparse(f_sparse(class, x)), len(cellSize), theta_0, Hook{}, Gamma{1.0/20}, Epsilon{1e-12}, ProximalOperator{p}); err != nil {
    test.Error(err)
  } else {
    if t.Vnorm(r.VsubV(r, z)); t.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
}

func Test6(test *testing.T) {

  // data
  cellSize  := []float64{
    1, 4, 1, 8, 1, 10, 1, 1, 1, 2, 1, 1, 3, 1, 7, 4, 1, 1, 7, 1}
  cellShape := []float64{
    1, 4, 1, 8, 1, 10, 1, 2, 1, 1, 1, 1, 3, 1, 5, 6, 1, 1, 7, 1}
  class := []float64{
    0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0}
  // x
  x := make([]SparseConstFloat64Vector, len(cellSize))
  for i := 0; i < len(cellSize); i++ {
    x[i] = NewSparseConstFloat64Vector([]int{0, 1, 2}, []float64{1.0, cellSize[i]-1.0, cellShape[i]-1.0}, 3)
  }

  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  z  := DenseFloat64Vector([]float64{-2.76467776, 0.17584927, 0.48174453})
  t1 := NullFloat64()
  t2 := NullDenseFloat64Vector(theta_0.Dim())
  p1 :=  proximalWrapper{&ProximalOperatorL1{1.0/2.5}}
  p2 := jitUpdateWrapper{&       JitUpdateL1{1.0/2.5}}

  trace1 := []ConstVector{}
  trace2 := []ConstVector{}
  hook1 := func(x ConstVector, step, lambda ConstScalar, i int) bool {
    // clone vector!
    trace1 = append(trace1, AsDenseFloat64Vector(x))
    return false
  }
  hook2 := func(x ConstVector, step, lambda ConstScalar, i int) bool {
    // clone vector!
    trace2 = append(trace2, AsDenseFloat64Vector(x))
    return false
  }

  if r, _, err := Run(Objective1Sparse(f_sparse(class, x)), len(cellSize), theta_0, Hook{hook1}, Gamma{1.0/20}, Epsilon{1e-12}, ProximalOperator{p1}); err != nil {
    test.Error(err)
  } else {
    if t1.Vnorm(t2.VsubV(r, z)); t1.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
  if r, _, err := Run(Objective1Sparse(f_sparse(class, x)), len(cellSize), theta_0, Hook{hook2}, Gamma{1.0/20}, Epsilon{1e-12}, JitUpdate{p2}); err != nil {
    test.Error(err)
  } else {
    if t1.Vnorm(t2.VsubV(r, z)); t1.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
  if len(trace1) == 0 || len(trace1) != len(trace2) {
    test.Error("test failed")
  } else {
    for i := 0; i < len(trace1); i++ {
      if t1.Vnorm(t2.VsubV(trace1[i], trace2[i])); t1.GetFloat64() > 1e-4 {
        test.Error("test failed")
      }
    }
  }
}

func Test7(test *testing.T) {

  // data
  cellSize  := []float64{
    1, 4, 1, 8, 1, 10, 1, 1, 1, 2, 1, 1, 3, 1, 7, 4, 1, 1, 7, 1}
  cellShape := []float64{
    1, 4, 1, 8, 1, 10, 1, 2, 1, 1, 1, 1, 3, 1, 5, 6, 1, 1, 7, 1}
  class := []float64{
    0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0}
  // x
  x := make([]SparseConstFloat64Vector, len(cellSize))
  for i := 0; i < len(cellSize); i++ {
    x[i] = NewSparseConstFloat64Vector([]int{0, 1, 2}, []float64{1.0, cellSize[i]-1.0, cellShape[i]-1.0}, 3)
  }

  theta_0 := NewDenseFloat64Vector([]float64{-1, 0.0, 0.0})
  z := DenseFloat64Vector([]float64{-2.81978662, 0.22409962, 0.45325202})
  t := NullFloat64()
  p := proximalWrapper{&ProximalOperatorTi{1.0/2.5}}

  if r, _, err := Run(Objective1Sparse(f_sparse(class, x)), len(cellSize), theta_0, Hook{}, Gamma{1.0/20}, Epsilon{1e-12}, ProximalOperator{p}); err != nil {
    test.Error(err)
  } else {
    if t.Vnorm(r.VsubV(r, z)); t.GetFloat64() > 1e-4 {
      test.Error("test failed")
    }
  }
}
