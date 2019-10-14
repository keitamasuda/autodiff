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

package vectorEstimator

/* -------------------------------------------------------------------------- */

import   "fmt"
import   "math"
import   "math/rand"

import . "github.com/pbenner/autodiff"
import   "github.com/pbenner/autodiff/algorithm/saga"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/vectorDistribution"
import . "github.com/pbenner/autodiff/logarithmetic"

import . "github.com/pbenner/threadpool"


/* -------------------------------------------------------------------------- */

type logisticRegression struct {
  Theta DenseBareRealVector
}

/* -------------------------------------------------------------------------- */

func (obj logisticRegression) Dim() int {
  return len(obj.Theta)-1
}

func (obj logisticRegression) LogPdfDense(x DenseConstRealVector) float64 {
  // set r to first element of theta
  r := float64(obj.Theta[0])
  n := x.Dim()
  for i := 1; i < n; i++ {
    r += float64(x[i])*float64(obj.Theta[i])
  }
  return -LogAdd(0.0, -r)
}

func (obj logisticRegression) LogPdfSparse(v SparseConstRealVector) float64 {
  x     := v.GetSparseValues ()
  index := v.GetSparseIndices()
  // set r to first element of theta
  r := float64(obj.Theta[0])
  // loop over x
  i := 0
  n := len(index)
  // skip first element
  if index[i] == 0 {
    i++
  }
  for ; i < n; i++ {
    r += float64(x[i])*float64(obj.Theta[index[i]])
  }
  return -LogAdd(0.0, -r)
}

/* -------------------------------------------------------------------------- */

type LogisticRegression struct {
  logisticRegression
  sparse     bool
  n          int
  x_sparse []SparseConstRealVector
  x_dense  [] DenseConstRealVector
  x        []ConstVector
  c        []bool
  stepSize   float64
  // optional parameters
  Balance         bool
  Epsilon         float64
  L1Reg           float64
  AutoReg         int
  Eta          [2]float64
  L2Reg           float64
  TiReg           float64
  StepSizeFactor  float64
  MaxIterations   int
  ClassWeights [2]float64
  Seed            int64
  Hook            func(x ConstVector, step, lambda ConstScalar, i int) bool
  sagaLogisticRegressionL1
}

/* -------------------------------------------------------------------------- */

func NewLogisticRegression(n int, sparse bool) (*LogisticRegression, error) {
  r := LogisticRegression{}
  r.logisticRegression.Theta = NullDenseBareRealVector(n)
  r.Epsilon         = 1e-5
  r.MaxIterations   = int(^uint(0) >> 1)
  r.ClassWeights[0] = 1.0
  r.ClassWeights[1] = 1.0
  r.Eta[0]          = 1.1
  r.Eta[1]          = 0.9
  r.StepSizeFactor  = 1.0
  r.sparse          = sparse
  return &r, nil
}

/* -------------------------------------------------------------------------- */

func (obj *LogisticRegression) Clone() *LogisticRegression {
  r := LogisticRegression{}
  // copy data and optional arguments
  r  = *obj
  r.logisticRegression.Theta = obj.logisticRegression.Theta.Clone()
  return &r
}

func (obj *LogisticRegression) CloneVectorEstimator() VectorEstimator {
  return obj.Clone()
}

/* -------------------------------------------------------------------------- */

func (obj *LogisticRegression) ScalarType() ScalarType {
  return BareRealType
}

func (obj *LogisticRegression) GetParameters() Vector {
  return obj.Theta
}

func (obj *LogisticRegression) SetParameters(x Vector) error {
  obj.Theta = AsDenseBareRealVector(x)
  return nil
}

/* -------------------------------------------------------------------------- */

func (obj *LogisticRegression) GetData() ([]ConstVector, int) {
  return obj.x, obj.n
}

// x_i = (1.0, x_i1, x_i2, ..., x_im, class_label)
func (obj *LogisticRegression) SetData(x []ConstVector, n int) error {
  if len(x) == 0 {
    return nil
  }
  if k := obj.logisticRegression.Dim()+2; x[0].Dim() != k {
    return fmt.Errorf("LogisticRegression: data has invalid dimension: got data of dimension `%d' but expected dimension `%d'", x[0].Dim(), k)
  }
  if obj.sparse {
    x_sparse := make([]ConstVector, len(x))
    for i, _ := range x {
      if x[i].Dim() != x[0].Dim() {
        return fmt.Errorf("data has inconsistent dimensions")
      }
      t := x[i].ConstSlice(0, x[i].Dim()-1)
      switch a := t.(type) {
      case SparseConstRealVector:
        x_sparse[i] = a
      default:
        x_sparse[i] = AsSparseConstRealVector(t)
      }
    }
    obj.SetSparseData(x_sparse, nil, n)
    obj.setStepSize()
  } else {
    x_dense := make([]ConstVector, len(x))
    for i, _ := range x {
      if x[i].Dim() != x[0].Dim() {
        return fmt.Errorf("data has inconsistent dimensions")
      }
      t := x[i].ConstSlice(0, x[i].Dim()-1)
      switch a := t.(type) {
      case DenseConstRealVector:
        x_dense[i] = a
      default:
        x_dense[i] = AsDenseConstRealVector(t)
      }
    }
    obj.SetDenseData(x_dense, nil, n)
    obj.setStepSize()
  }
  obj.c = make([]bool       , len(x))
  obj.x = make([]ConstVector, len(x))
  for i, _ := range x {
    switch a := x[i].(type) {
    case SparseConstRealVector:
      // do not use ValueAt to prevent that an index
      // for the sparse vector is constructed
      if j, v := a.First(); j != 0 || v != 1.0 {
        return fmt.Errorf("first element of data vector must be set to one")
      }
      if j, v := a.Last (); j != a.Dim()-1 {
        // last entry is not the class label =>
        // class is zero
        obj.c[i] = false
      } else {
        switch v {
        case 1.0: obj.c[i] = true
        case 0.0: obj.c[i] = false
        default : return fmt.Errorf("invalid class label `%f'", v)
        }
      }
    default:
      if x[i].ValueAt(0) != 1.0 {
        return fmt.Errorf("first element of data vector must be set to one")
      }
      v := x[i].ValueAt(x[i].Dim()-1)
      switch v {
      case 1.0: obj.c[i] = true
      case 0.0: obj.c[i] = false
      default : return fmt.Errorf("invalid class label `%f'", v)
      }
    }
    obj.x[i] = x[i]
  }
  obj.setLabels(obj.c)
  return nil
}

func (obj *LogisticRegression) setLabels(c []bool) {
  obj.c = c
  if obj.Balance {
    n1 := 0
    n0 := 0
    for i := 0; i < len(obj.c); i++ {
      switch obj.c[i] {
      case true : n1++
      case false: n0++
      }
    }
    obj.ClassWeights[1] = float64(n0+n1)/float64(2*n1)
    obj.ClassWeights[0] = float64(n0+n1)/float64(2*n0)
  }
}

func (obj *LogisticRegression) SetSparseData(x []ConstVector, c []bool, n int) error {
  obj.n        = n
  obj.x        = nil
  obj.x_sparse = make([]SparseConstRealVector, len(x))
  obj.x_dense  = nil
  obj.sparse   = true
  for i, _ := range x {
    if k := obj.logisticRegression.Dim()+1; x[i].Dim() != k {
      return fmt.Errorf("LogisticRegression.SetSparseData: data has invalid dimension: got data of dimension `%d' but expected dimension `%d'", x[i].Dim(), k)
    }
    switch a := x[i].(type) {
    case SparseConstRealVector:
      obj.x_sparse[i] = a
    default:
      return fmt.Errorf("data is not of type SparseConstRealVector")
    }
  }
  obj.setStepSize()
  obj.setLabels(c)
  return nil
}

func (obj *LogisticRegression) SetDenseData(x []ConstVector, c []bool, n int) error {
  obj.n        = n
  obj.x        = nil
  obj.x_sparse = nil
  obj.x_dense  = make([]DenseConstRealVector, len(x))
  obj.sparse   = false
  for i, _ := range x {
    switch a := x[i].(type) {
    case DenseConstRealVector:
      obj.x_dense[i] = a
    default:
      return fmt.Errorf("data is not of type DenseConstRealVector")
    }
  }
  obj.setStepSize()
  obj.setLabels(c)
  return nil
}

/* -------------------------------------------------------------------------- */

func (obj *LogisticRegression) Estimate(gamma ConstVector, p ThreadPool) error {
  if gamma != nil {
    panic("internal error")
  }
  { m := 0
    if obj.L1Reg != 0.0 { m++ }
    if obj.L2Reg != 0.0 { m++ }
    if obj.TiReg != 0.0 { m++ }
    if m > 1 {
      return fmt.Errorf("multiple regularizations are not supported")
    }
  }
  var proxop    saga.ProximalOperatorType
  var jitUpdate saga.JitUpdateType
  switch {
  case obj.sparse && obj.L2Reg == 0.0 && obj.TiReg == 0.0:
    // use specialized saga implementation
    if err := obj.sagaLogisticRegressionL1.Initialize(saga.Objective1Sparse(obj.f_sparse), len(obj.x_sparse), obj.Theta,
      saga.L1Regularization{obj.L1Reg},
      saga.AutoReg         {obj.AutoReg},
      saga.Gamma           {obj.stepSize},
      saga.Seed            {obj.Seed}, p); err != nil {
      return err
    }
    if r, s, err := obj.sagaLogisticRegressionL1.Execute(
      saga.Epsilon         {obj.Epsilon},
      saga.Eta             {obj.Eta},
      saga.MaxIterations   {obj.MaxIterations},
      saga.Hook            {obj.Hook}); err != nil {
      return err
    } else {
      obj.Seed = s
      obj.SetParameters(r)
      return nil
    }
  case obj.L1Reg != 0.0: proxop = proximalWrapper{&saga.ProximalOperatorL1{obj.L1Reg}}
  case obj.L2Reg != 0.0: proxop = proximalWrapper{&saga.ProximalOperatorL2{obj.L2Reg}}
  case obj.TiReg != 0.0: proxop = proximalWrapper{&saga.ProximalOperatorTi{obj.TiReg}}
  }
  if obj.sparse {
    if r, s, err := saga.Run(saga.Objective1Sparse(obj.f_sparse), len(obj.x_sparse), obj.Theta,
      saga.Hook            {obj.Hook},
      saga.AutoReg         {obj.AutoReg},
      saga.Gamma           {obj.stepSize},
      saga.Epsilon         {obj.Epsilon},
      saga.Eta             {obj.Eta},
      saga.MaxIterations   {obj.MaxIterations},
      saga.Seed            {obj.Seed},
      saga.ProximalOperator{proxop},
      saga.JitUpdate       {jitUpdate}); err != nil {
      return err
    } else {
      obj.Seed = s
      obj.SetParameters(r)
    }
  } else {
    if r, s, err := saga.Run(saga.Objective1Dense(obj.f_dense), len(obj.x_dense), obj.Theta,
      saga.Hook            {obj.Hook},
      saga.AutoReg         {obj.AutoReg},
      saga.Gamma           {obj.stepSize},
      saga.Eta             {obj.Eta},
      saga.Epsilon         {obj.Epsilon},
      saga.MaxIterations   {obj.MaxIterations},
      saga.Seed            {obj.Seed},
      saga.ProximalOperator{proxop},
      saga.JitUpdate       {jitUpdate}); err != nil {
      return err
    } else {
      obj.Seed = s
      obj.SetParameters(r)
    }
  }
  return nil
}

func (obj *LogisticRegression) EstimateOnData(x []ConstVector, gamma ConstVector, p ThreadPool) error {
  if err := obj.SetData(x, len(x)); err != nil {
    return err
  }
  return obj.Estimate(gamma, p)
}

func (obj *LogisticRegression) GetEstimate() (VectorPdf, error) {
  return vectorDistribution.NewLogisticRegression(obj.Theta)
}

/* -------------------------------------------------------------------------- */

func (obj *LogisticRegression) setStepSize() {
  max_squared_sum := 0.0
  if obj.sparse {
    for _, x := range obj.x_sparse {
      r  := 0.0
      it := x.ConstIterator()
      // skip first element
      if it.Ok() {
        it.Next()
      }
      for ; it.Ok(); it.Next() {
        r += it.GetValue()*it.GetValue()
      }
      if r > max_squared_sum {
        max_squared_sum = r
      }
    }
  } else {
    for _, x := range obj.x_dense {
      r  := 0.0
      it := x.ConstIterator()
      // skip first element
      if it.Ok() {
        it.Next()
      }
      for ; it.Ok(); it.Next() {
        r += it.GetValue()*it.GetValue()
      }
      if r > max_squared_sum {
        max_squared_sum = r
      }
    }
  }
  L := (0.25*(max_squared_sum + 1.0) + obj.L2Reg/float64(obj.n))
  L *= math.Max(obj.ClassWeights[0], obj.ClassWeights[1])
  obj.stepSize  = 1.0/(2.0*L + math.Min(2.0*obj.L2Reg, L))
  obj.stepSize *= obj.StepSizeFactor
}

/* -------------------------------------------------------------------------- */

type proximalWrapper struct {
  saga.ProximalOperatorType
}

func (obj proximalWrapper) Eval(x DenseBareRealVector, w DenseBareRealVector, t *BareReal) {
  obj.ProximalOperatorType.Eval(x, w, t)
  // do not regularize intercept
  x.AT(0).SET(w.AT(0))
}

/* -------------------------------------------------------------------------- */

type jitUpdateWrapper struct {
  saga.JitUpdateType
}

func (obj jitUpdateWrapper) Update(x, y BareReal, k, m int) BareReal {
  // do not regularize intercept
  if k == 0 {
    return x - BareReal(m)*y
  } else {
    return obj.JitUpdateType.Update(x, y, k, m)
  }
}

/* -------------------------------------------------------------------------- */

func (obj *LogisticRegression) f_dense(i int, theta DenseBareRealVector) (ConstReal, ConstReal, DenseConstRealVector, error) {
  x := obj.x_dense
  y := ConstReal(0.0)
  w := ConstReal(0.0)
  if i >= len(x) {
    return y, w, x[i], fmt.Errorf("index out of bounds")
  }
  obj.logisticRegression.Theta = theta

  r := obj.logisticRegression.LogPdfDense(x[i])

  if math.IsNaN(r) {
    return y, w, x[i], fmt.Errorf("NaN value detected")
  }
  y = ConstReal(r)
  if obj.c[i] {
    w = ConstReal(obj.ClassWeights[1]*(math.Exp(r) - 1.0))
  } else {
    w = ConstReal(obj.ClassWeights[0]*(math.Exp(r)))
  }
  return y, w, x[i], nil
}

func (obj *LogisticRegression) f_sparse(i int, theta DenseBareRealVector) (ConstReal, ConstReal, SparseConstRealVector, error) {
  x := obj.x_sparse
  y := ConstReal(0.0)
  w := ConstReal(0.0)
  if len(theta) == 0 {
    return y, w, x[i], nil
  }
  if i >= len(x) {
    return y, w, x[i], fmt.Errorf("index out of bounds")
  }
  obj.logisticRegression.Theta = theta

  r := obj.logisticRegression.LogPdfSparse(x[i])

  if math.IsNaN(r) {
    return y, w, x[i], fmt.Errorf("NaN value detected")
  }
  y = ConstReal(r)
  if obj.c[i] {
    w = ConstReal(obj.ClassWeights[1]*(math.Exp(r) - 1.0))
  } else {
    w = ConstReal(obj.ClassWeights[0]*(math.Exp(r)))
  }
  return y, w, x[i], nil
}

/* -------------------------------------------------------------------------- */

type sagaJitUpdateL1 struct {
  saga.JitUpdateL1
}

func (obj sagaJitUpdateL1) Update(x, y BareReal, k, m int) BareReal {
  // do not regularize intercept
  if k == 0 {
    return x - BareReal(m)*y
  } else {
    return obj.JitUpdateL1.Update(x, y, k, m)
  }
}

/* -------------------------------------------------------------------------- */

type gradientJit struct {
  G SparseConstRealVector
  W ConstReal
}

func (obj gradientJit) Add(v DenseBareRealVector) {
  g := obj.G.GetSparseValues()
  for i, k := range obj.G.GetSparseIndices() {
    v[k] = v[k] + BareReal(obj.W.GetValue()*g[i])
  }
}

func (obj *gradientJit) Set(w ConstReal, g SparseConstRealVector) {
  obj.G = g
  obj.W = w
}

func (g1 gradientJit) Update(g2 gradientJit, v DenseBareRealVector) {
  v1 := g1.G.GetSparseValues()
  v2 := g2.G.GetSparseValues()
  if v1 != nil && &v1[0] == &v2[0] {
    for i, k := range g1.G.GetSparseIndices() {
      v[k] += (g2.W - g1.W)*ConstReal(v1[i])
    }
  } else {
    for i, k := range g1.G.GetSparseIndices() {
      v[k] -= g1.W*ConstReal(v1[i])
    }
    for i, k := range g2.G.GetSparseIndices() {
      v[k] += g2.W*ConstReal(v2[i])
    }
  }
}

/* -------------------------------------------------------------------------- */

type sagaLogisticRegressionL1worker struct {
  indices       []int
  f               saga.Objective1Sparse
  x0              DenseBareRealVector
  x1              DenseBareRealVector
  xs            []bool
  xk            []int
  d               int
  ns              int
  dict          []gradientJit
  s               DenseBareRealVector
  cumulative_sums DenseBareRealVector
  t_n             BareReal
  t_g             BareReal
  jit             sagaJitUpdateL1
  rand           *rand.Rand
}

func (obj *sagaLogisticRegressionL1worker) Initialize(
  f saga.Objective1Sparse,
  indices []int,
  x DenseBareRealVector,
  l1reg  saga.L1Regularization,
  gamma saga.Gamma,
  seed saga.Seed) error {

  n := len(indices)

  obj.f  = f
  obj.x0 = AsDenseBareRealVector(x)
  obj.x1 = AsDenseBareRealVector(x)
  obj.xk = make([]int,  x.Dim())
  obj.xs = make([]bool, n)
  obj.ns = 0
  obj.cumulative_sums = NullDenseBareRealVector(n)
  obj.indices         = indices

  // length of gradient
  obj.d = x.Dim()

  // some constants
  obj.t_n = BareReal(0.0)
  obj.t_g = BareReal(gamma.Value)

  obj.jit.SetLambda(l1reg.Value*gamma.Value/float64(n))

  // sum of gradients
  obj.s = NullDenseBareRealVector(obj.d)
  // initialize s and d
  obj.dict = make([]gradientJit, n)
  if seed.Value != -1 {
    obj.rand = rand.New(rand.NewSource(seed.Value))
  }
  return nil
}

func (obj *sagaLogisticRegressionL1worker) jitUpdates(i_, j int) error {
  if _, _, gt, err := obj.f(obj.indices[j], nil); err != nil {
    return err
  } else {
    for _, k := range gt.GetSparseIndices() {
      if m := i_-obj.xk[k]; m > 0 {
        cum_sum := obj.cumulative_sums[i_-1]
        if obj.xk[k] != 0 {
          cum_sum -= obj.cumulative_sums[obj.xk[k]-1]
        }
        obj.x1[k] = obj.jit.Update(obj.x1[k], cum_sum*obj.s[k]/BareReal(m), k, m)
      }
    }
  }
  return nil
}

func (obj *sagaLogisticRegressionL1worker) jitUpdatesMissing(n int) {
  for k := 0; k < obj.x1.Dim(); k++ {
    if m := n-obj.xk[k]; m > 0 {
      cum_sum := obj.cumulative_sums[n-1]
      if obj.xk[k] != 0 {
        cum_sum -= obj.cumulative_sums[obj.xk[k]-1]
      }
      obj.x1[k] = obj.jit.Update(obj.x1[k], cum_sum*obj.s[k]/BareReal(m), k, m)
    }
    // reset xk
    obj.xk[k] = 0
  }
}

func (obj *sagaLogisticRegressionL1worker) gradientUpdates(i_, j int, g1, g2 gradientJit) {
  v1 := g1.G.GetSparseValues()
  v2 := g2.G.GetSparseValues()
  if v1 == nil || &v1[0] != &v2[0] {
    // data vectors are different
    for i, k := range g1.G.GetSparseIndices() {
      obj.x1[k] = obj.x1[k] + obj.t_g*(1.0 - 1.0/obj.t_n)*g1.W*BareReal(v1[i])
      obj.xk[k] = i_
    }
    for i, k := range g2.G.GetSparseIndices() {
      obj.x1[k] = obj.x1[k] - obj.t_g*(1.0 - 1.0/obj.t_n)*g2.W*BareReal(v2[i])
      obj.xk[k] = i_
    }
  } else {
    // data vectors are identical
    c := BareReal(g2.W - g1.W)
    for i, k := range g2.G.GetSparseIndices() {
      obj.x1[k] = obj.x1[k] - obj.t_g*(1.0 - 1.0/obj.t_n)*c*BareReal(v2[i])
      obj.xk[k] = i_
    }
  }
}

func (obj *sagaLogisticRegressionL1worker) Iterate(epoch int) error {
  n := len(obj.xs)
  var g1 gradientJit
  var g2 gradientJit
  for i_ := 0; i_ < n; i_++ {
    j := i_
    if obj.rand != nil {
      j = obj.rand.Intn(n)
    }
    if !obj.xs[j] {
      obj.xs[j] = true
      obj.ns   += 1
      obj.t_n   = BareReal(obj.ns)
    }
    if i_ == 0 {
      obj.cumulative_sums[0 ] = obj.t_g/obj.t_n
    } else {
      obj.cumulative_sums[i_] = obj.cumulative_sums[i_-1] + obj.t_g/obj.t_n
    }
    // get old gradient
    g1 = obj.dict[j]
    // perform jit updates for all x_i where g_i != 0
    if err := obj.jitUpdates(i_, j); err != nil {
      return nil
    }
    // evaluate objective function
    if _, w, gt, err := obj.f(obj.indices[j], obj.x1); err != nil {
      return err
    } else {
      g2.Set(w, gt)
    }
    // perform actual gradient step
    obj.gradientUpdates(i_, j, g1, g2)
    // update gradient avarage
    g1.Update(g2, obj.s)
    // update dictionary
    obj.dict[j].Set(g2.W, g2.G)
  }
  // compute missing updates of x1
  obj.jitUpdatesMissing(n)
  return nil
}

func (obj *sagaLogisticRegressionL1worker) ComputeLambda() float64 {
  return float64(len(obj.indices))*obj.jit.GetLambda()/obj.t_g.GetValue()
}

/* -------------------------------------------------------------------------- */

type sagaLogisticRegressionL1 struct {
  Workers []sagaLogisticRegressionL1worker
  Indices []int
  Pool      ThreadPool
  autoReg   saga.AutoReg
  n_x_old   int
  n_x_new   int
  l1_step   float64
}

func (obj *sagaLogisticRegressionL1) Initialize(
  f saga.Objective1Sparse,
  n int,
  x DenseBareRealVector,
  l1reg  saga.L1Regularization,
  autoReg saga.AutoReg,
  gamma saga.Gamma,
  seed saga.Seed,
  pool ThreadPool) error {

  obj.autoReg = autoReg
  // prevent that in auto-lambda mode the step size is initialized to zero
  if autoReg.Value > 0 && l1reg.Value == 0.0 {
    l1reg.Value = 1.0
  }
  // number of non-zero parameters used for auto-lambda mode
  obj.n_x_old = 0
  obj.n_x_new = 0
  // step size for auto-lambda mode
  obj.l1_step = 0.5*l1reg.Value*gamma.Value/float64(n)
  // slice of data indices
  obj.Indices = make([]int, n)
  for i := 0; i < n; i++ {
    obj.Indices[i] = i
  }
  m := n/pool.NumberOfThreads()
  // create a slice for every worker
  indices := make([][]int, pool.NumberOfThreads())
  for i, k := 0, 0; i < pool.NumberOfThreads(); i++ {
    if i+1 == pool.NumberOfThreads() {
      indices[i] = obj.Indices[k:n]
    } else {
      indices[i] = obj.Indices[k:k+m]; k += m
    }
  }
  obj.Workers = make([]sagaLogisticRegressionL1worker, pool.NumberOfThreads())
  for i := 0; i < pool.NumberOfThreads(); i++ {
    if err := obj.Workers[i].Initialize(f, indices[i], x, l1reg, gamma, seed); err != nil {
      return err
    }
  }
  obj.Pool = pool
  return nil
}

func (obj *sagaLogisticRegressionL1) Execute(
  epsilon saga.Epsilon,
  eta saga.Eta,
  maxIterations saga.MaxIterations,
  hook saga.Hook) (DenseBareRealVector, int64, error) {

  x0 := obj.Workers[0].x0
  x1 := obj.Workers[0].x1
  for epoch := 0; epoch < maxIterations.Value; epoch++ {
    // copy initial value
    for i := 1; i < len(obj.Workers); i++ {
      obj.Workers[i].x0.SET(x0)
    }
    if err := obj.Pool.RangeJob(0, len(obj.Workers), func(i int, pool ThreadPool, erf func() error) error {
      return obj.Workers[i].Iterate(epoch)
    }); err != nil {
      seed := int64(-1)
      if obj.Workers[0].rand != nil {
        seed = obj.Workers[0].rand.Int63()
      }
      return x1, seed, err
    }
    // compute mean
    if len(obj.Workers) > 1 {
      for i := 1; i < len(obj.Workers); i++ {
        x1.VADDV(x1, obj.Workers[i].x1)
      }
      t := ConstReal(len(obj.Workers))
      x1.VDIVS(x1, &t)
    }
    // check convergence
    if stop, delta, err := saga.EvalStopping(x0, x1, epsilon.Value*obj.Workers[0].t_g.GetValue()); stop {
      return x1, obj.Workers[0].rand.Int63(), err
    } else {
      // execute hook if available
      if hook.Value != nil && hook.Value(x1, ConstReal(delta), ConstReal(obj.Workers[0].ComputeLambda()), epoch) {
        break
      }
    }
    // update lambda
    if obj.autoReg.Value > 0 {
      obj.n_x_new = 0
      // count number of non-zero entries
      for k := 1; k < x1.Dim(); k++ {
        if x1[k] != 0.0 {
          obj.n_x_new += 1
        }
      }
      switch {
      case obj.n_x_old < obj.autoReg.Value && obj.n_x_new < obj.autoReg.Value: fallthrough
      case obj.n_x_old > obj.autoReg.Value && obj.n_x_new > obj.autoReg.Value:
        obj.l1_step = eta.Value[0]*obj.l1_step
      default:
        obj.l1_step = eta.Value[1]*obj.l1_step
      }
      for i, _ := range obj.Workers {
        if obj.n_x_new < obj.autoReg.Value {
          obj.Workers[i].jit.SetLambda(obj.Workers[i].jit.GetLambda() - obj.l1_step)
        } else
        if obj.n_x_new > obj.autoReg.Value {
          obj.Workers[i].jit.SetLambda(obj.Workers[i].jit.GetLambda() + obj.l1_step)
        }
        if obj.Workers[i].jit.GetLambda() < 0.0 {
          obj.Workers[i].jit.SetLambda(0.0)
        }
      }
      // swap old and new counts
      obj.n_x_old, obj.n_x_new = obj.n_x_new, obj.n_x_old
    }
    // shuffle all indices
    if len(obj.Workers) > 1 && obj.Workers[0].rand != nil {
      obj.Workers[0].rand.Shuffle(len(obj.Indices), func(i, j int) {
        obj.Indices[i], obj.Indices[j] = obj.Indices[j], obj.Indices[i]
      })
    }
    x0.SET(x1)
  }
  seed := int64(-1)
  if obj.Workers[0].rand != nil {
    seed = obj.Workers[0].rand.Int63()
  }
  return x1, seed, nil
}
