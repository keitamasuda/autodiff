## Documentation

Autodiff is a numerical optimization and linear algebra library for the Go / Golang programming language. It implements basic automatic differentation for many mathematical routines. The documentation of this package can be found [here](https://godoc.org/github.com/pbenner/autodiff).

## Scalars

Autodiff defines three different scalar types. A *Scalar* contains a single mutable value that can be the result of a mathematical operation, whereas the value of a *ConstScalar* is constant and fixed when the scalar is created. Automatic differentiation is implemented by *MagicScalar* types and allow to compute first and second order derivatives. Autodiff implements the following scalars:

| Scalar       | Implemented interfaces
|--------------|------------------------------------------------------ |
| ConstInt8    | ConstScalar                                           |
| ConstInt16   | ConstScalar                                           |
| ConstInt32   | ConstScalar                                           |
| ConstInt64   | ConstScalar                                           |
| ConstInt     | ConstScalar                                           |
| ConstFloat32 | ConstScalar                                           |
| ConstFloat64 | ConstScalar                                           |
| Int8         | ConstScalar, Scalar                                   |
| Int16        | ConstScalar, Scalar                                   |
| Int32        | ConstScalar, Scalar                                   |
| Int64        | ConstScalar, Scalar                                   |
| Int          | ConstScalar, Scalar                                   |
| Float32      | ConstScalar, Scalar                                   |
| Float64      | ConstScalar, Scalar                                   |
| Real32       | ConstScalar, Scalar, MagicScalar                      |
| Real64       | ConstScalar, Scalar, MagicScalar                      |

The *ConstScalar*, *Scalar* and *MagicScalar* interfaces define the following operations:

| Function     | Description                                           |
|--------------|------------------------------------------------------ |
| GetInt8      | Get value as int8                                     |
| GetInt16     | Get value as int16                                    |
| GetInt32     | Get value as int32                                    |
| GetInt64     | Get value as int64                                    |
| GetInt       | Get value as int                                      |
| GetFloat32   | Get value as float32                                  |
| GetFloat64   | Get value as float64                                  |
| Equals       | Check if two constants are equal                      |
| Greater      | True if first constant is greater                     |
| Smaller      | True if first constant is smaller                     |
| Sign         | Returns the sign of the scalar                        |

The *Scalar* and *MagicScalar* interfaces define the following operations:

| Function     | Description                                           |
|--------------|------------------------------------------------------ |
| SetInt8      | Set value by passing an int8 variable                 |
| SetInt16     | Set value by passing an int16 variable                |
| SetInt32     | Set value by passing an int32 variable                |
| SetInt64     | Set value by passing an int64 variable                |
| SetInt       | Set value by passing an int variable                  |
| SetFloat32   | Set value by passing an float32 variable              |
| SetFloat64   | Set value by passing an float64 variable              |

The *Scalar* and *MagicScalar* interfaces define the following mathematical operations:

| Function     | Description                                           |
| ------------ | ----------------------------------------------------- |
| Min          | Minimum                                               |
| Max          | Maximum                                               |
| Abs          | Absolute value                                        |
| Sign         | Sign                                                  |
| Neg          | Negation                                              |
| Add          | Addition                                              |
| Sub          | Substraction                                          |
| Mul          | Multiplication                                        |
| Div          | Division                                              |
| Pow          | Power                                                 |
| Sqrt         | Square root                                           |
| Exp          | Exponential function                                  |
| Log          | Logarithm                                             |
| Log1p        | Logarithm of 1+x                                      |
| Log1pExp     | Logarithm of 1+Exp(x)                                 |
| Logistic     | Standard logistic function                            |
| Erf          | Error function                                        |
| Erfc         | Complementary error function                          |
| LogErfc      | Log complementary error function                      |
| Sigmoid      | Numerically stable sigmoid function                   |
| Sin          | Sine                                                  |
| Sinh         | Hyperbolic sine                                       |
| Cos          | Cosine                                                |
| Cosh         | Hyperbolic cosine                                     |
| Tan          | Tangent                                               |
| Tanh         | Hyperbolic tangent                                    |
| LogAdd       | Addition on log scale                                 |
| LogSub       | Substraction on log scale                             |
| SmoothMax    | Differentiable maximum                                |
| LogSmoothMax | Differentiable maximum on log scale                   |
| Gamma        | Gamma function                                        |
| Lgamma       | Log gamma function                                    |
| Mlgamma      | Multivariate log gamma function                       |
| GammaP       | Lower incomplete gamma function                       |
| BesselI      | Modified Bessel function of the first kind            |
| LogBesselI   | Log of the Modified Bessel function of the first kind |

## Vectors and Matrices

Autodiff supports vectors and matrices including basic linear algebra operations. Vectors support the following linear algebra operations:

| Function | Description                      |
| -------- | -------------------------------- |
| VaddV    | Element-wise addition            |
| VsubV    | Element-wise substraction        |
| VmulV    | Element-wise multiplication      |
| VdivV    | Element-wise division            |
| VaddS    | Addition of a scalar             |
| VsubS    | Substraction of a scalar         |
| VmulS    | Multiplication with a scalar     |
| VdivS    | Division by a scalar             |
| VdotV    | Dot product                      |

Matrices support the following linear algebra operations:

| Function | Description                      |
| -------- | -------------------------------- |
| MaddM    | Element-wise addition            |
| MsubM    | Element-wise substraction        |
| MmulM    | Element-wise multiplication      |
| MdivM    | Element-wise division            |
| MaddS    | Addition of a scalar             |
| MsubS    | Substraction of a scalar         |
| MmulS    | Multiplication with a scalar     |
| MdivS    | Division by a scalar             |
| MdotM    | Matrix product                   |
| Outer    | Outer product                    |

## Algorithms

The algorithms package contains more complex linear algebra and optimization routines:

| Package             | Description                                             |
| ------------------- | ------------------------------------------------------- |
| bfgs                | Broyden-Fletcher-Goldfarb-Shanno (BFGS) algorithm       |
| blahut              | Blahut algorithm (channel capacity)                     |
| cholesky            | Cholesky and LDL factorization                          |
| determinant         | Matrix determinants                                     |
| eigensystem         | Compute Eigenvalues and Eigenvectors                    |
| gaussJordan         | Gauss-Jordan algorithm                                  |
| gradientDescent     | Vanilla gradient desent algorithm                       |
| gramSchmidt         | Gram-Schmidt algorithm                                  |
| hessenbergReduction | Matrix Hessenberg reduction                             |
| lineSearch          | Line-search (satisfying the Wolfe conditions)           |
| matrixInverse       | Matrix inverse                                          |
| msqrt               | Matrix square root                                      |
| msqrtInv            | Inverse matrix square root                              |
| newton              | Newton's method (root finding and optimization)         |
| qrAlgorithm         | QR-Algorithm for computing Schur decompositions         |
| rprop               | Resilient backpropagation                               |
| svd                 | Singular Value Decomposition (SVD)                      |
| saga                | SAGA stochastic average gradient descent method         |

## Basic usage

Import the autodiff library with
```go
  import . "github.com/pbenner/autodiff"
```
A scalar holding the value *1.0* can be defined in several ways, i.e.
```go
  a := NullScalar(Real64Type)
  a.SetFloat64(1.0)
  b := NewReal64(1.0)
  c := NewFloat64(1.0)
```
*a* and *b* are both *MagicScalar*s, however *a* has type *Scalar* whereas *b* has type **Real64* which implements the *Scalar* interface. Variable *c* is of type *Float64* which cannot carry any derivatives. Basic operations such as additions are defined on all Scalars, i.e.
```go
  a.Add(a, b)
```
which stores the result of adding *a* and *b* in *a*. If *autodiff/simple* is imported, one may also use
```go
  d := Add(a, b)
```
where the result is stored in a new variable *d*. The *ConstFloat64* type allows to define float64 constants without allocation of additional memory. For instance
```go
  a.Add(a, ConstFloat64(1.0))
```
adds a constant value to *a* where a type cast is used to define the constant *1.0*.

To differentiate a function
```go
  f := func(x, y ConstScalar) MagicScalar {
    // compute f(x,y) = x*y^3 + 4
    r := NewReal64()
    r.Pow(y, ConstFloat64(3.0))
    r.Mul(r, x)
    r.Add(r, ConstFloat64(4.0))
    return r
  }
```
first two reals are defined
```go
  x := NewReal(2)
  y := NewReal(4)
```
that store the value at which the derivative of *f* should be evaluated. Afterwards, *x* and *y* must be defined as variables with
```go
  Variables(2, x, y)
```
where the first argument says that derivatives up to second order should be computed. After evaluating *f*, i.e.
```go
  z := f(x, y)
```
the function value at *(x,y) = (2, 4)* can be retrieved with *z.GetFloat64()*. The first and second partial derivatives can be accessed with *z.GetDerivative(i)* and *z.GetHessian(i, j)*, where the arguments specify the index of the variable. For instance, the derivative of *f* with respect to *x* is returned by *z.GetDerivative(0)*, whereas the derivative with respect to *y* by *z.GetDerivative(1)*.

## Basic linear algebra

Vectors and matrices can be created with
```go
  v := NewDenseFloat64Vector([]float64{1,2})
  m := NewDenseFloat64Matrix([]float64{1,2,3,4}, 2, 2)

  v_ := NewDenseReal64Vector([]float64{1,2})
  m_ := NewDenseReal64Matrix([]float64{1,2,3,4}, 2, 2)
```
where *v* has length 2 and *m* is a 2x2 matrix. With
```go
  v := NullDenseFloat64Vector(2)
  m := NullDenseFloat64Matrix(2, 2)
```
all values are initially set to zero. Vector and matrix elements can be accessed with the *At*, *MagicAt* or *ConstAt* methods, which return a reference to the scalar implementing either a *Scalar*, *MagicScalar* or *ConstScalar*, i.e.
```go
  m.At(1,1).Add(v.ConstAt(0), v.ConstAt(1))
```
adds the first two values in *v* and stores the result in the lower right element of the matrix *m*. Autodiff supports basic linear algebra operations, for instance, the vector matrix product can be computed with
```go
  w := NullDenseFloat64Vector(2)
  w.MdotV(m, v)
```
where the result is stored in w. Other operations, such as computing the eigenvalues and eigenvectors of a matrix, require importing the respective package from the algorithm library, i.e.
```go
  import "github.com/pbenner/autodiff/algorithm/eigensystem"

  lambda, _, _ := eigensystem.Run(m)
```

## Examples

### Gradient descent

Compare vanilla gradient descent with resilient backpropagation
```go
  import . "github.com/pbenner/autodiff"
  import   "github.com/pbenner/autodiff/algorithm/gradientDescent"
  import   "github.com/pbenner/autodiff/algorithm/rprop"
  import . "github.com/pbenner/autodiff/simple"

  f := func(x_ ConstVector) MagicScalar {
    x := x_.ConstAt(0)
    // x^4 - 3x^3 + 2
    r := NewReal64()
    s := NewReal64()
    r.Pow(x.ConstAt(0), ConstFloat64(4.0)
    s.Mul(ConstFloat64(3.0), s.Pow(x, ConstFloat64(3.0)))
    r.Add(ConstFloat64(2.0), r.Add(r, s))
    return r
  }
  x0 := NewDenseFloat64Vector([]float64{8})
  // vanilla gradient descent
  xn1, _ := gradientDescent.Run(f, x0, 0.0001, gradientDescent.Epsilon{1e-8})
  // resilient backpropagation
  xn2, _ := rprop.Run(f, x0, 0.0001, 0.4, rprop.Epsilon{1e-8})
```
![Gradient descent](demo/example1/example1.png)


### Matrix inversion

Compute the inverse *r* of a matrix *m* by minimizing the Frobenius norm *||mb - I||*
```go
  import . "github.com/pbenner/autodiff"
  import   "github.com/pbenner/autodiff/algorithm/rprop"
  import . "github.com/pbenner/autodiff/simple"

  // define matrix r
  m := NewDenseFloat64Matrix([]float64{1,2,3,4}, 2, 2)
  // create identity matrix I
  I := NullDenseFloat64Matrix(2, 2)
  I.SetIdentity()

  // magic variables for computing the Frobenius norm and its derivative
  t := NewDenseReal64Matrix(2, 2)
  s := NewReal64()
  // objective function
  f := func(x ConstVector) MagicScalar {
    t.Set(x)
    s.Mnorm(t.MsubM(t.MmulM(m, t), I))
    return s
  }
  r, _ := rprop.Run(f, r.GetValues(), 0.01, 0.1, rprop.Epsilon{1e-12})
```

### Newton's method

Find the root of a function *f* with initial value *x0 = (1,1)*

```go
  import . "github.com/pbenner/autodiff"
  import   "github.com/pbenner/autodiff/algorithm/newton"
  import . "github.com/pbenner/autodiff/simple"

  t := NullReal64()

  f := func(x ConstVector) MagicVector {
    x1 := x.ConstAt(0)
    x2 := x.ConstAt(1)
    y  := NullDenseReal64Vector(2)
    y1 := y.At(0)
    y2 := y.At(1)
    // y1 = x1^2 + x2^2 - 6
    t .Pow(x1, ConstFloat64(2.0))
    y1.Add(y1, t)
    t .Pow(x2, ConstFloat64(2.0))
    y1.Add(y1, t)
    y1.Sub(y1, ConstFloat64(6.0))
    // y2 = x1^3 - x2^2
    t .Pow(x1, ConstFloat64(3.0))
    y2.Add(y2, t)
    t .Pow(x2, ConstFloat64(2.0))
    y2.Sub(y2, t)

    return y
  }

  x0    := NewDenseFloat64Vector([]float64{1,1})
  xn, _ := newton.RunRoot(f, x0, newton.Epsilon{1e-8})
```

### Minimize Rosenbrock's function
Compare Newton's method, BFGS and Rprop for minimizing Rosenbrock's function

```go
  import   "fmt"

  import . "github.com/pbenner/autodiff"
  import   "github.com/pbenner/autodiff/algorithm/rprop"
  import   "github.com/pbenner/autodiff/algorithm/bfgs"
  import   "github.com/pbenner/autodiff/algorithm/newton"

  f := func(x ConstVector) (MagicScalar, error) {
    // f(x1, x2) = (a - x1)^2 + b(x2 - x1^2)^2
    // a = 1
    // b = 100
    // minimum: (x1,x2) = (a, a^2)
    a := ConstFloat64(  1.0)
    b := ConstFloat64(100.0)
    c := ConstFloat64(  2.0)
    s := NullReal64()
    t := NullReal64()
    s.Pow(s.Sub(a, x.ConstAt(0)), c)
    t.Mul(b, t.Pow(t.Sub(x.ConstAt(1), t.Mul(x.ConstAt(0), x.ConstAt(0))), c))
    s.Add(s, t)
    return s, nil
  }
  hook_rprop := func(gradient, step []float64, x ConstVector, y ConstScalar) bool {
    fmt.Fprintf(fp1, "%s\n", x.Table())
    fmt.Println("x       :", x)
    fmt.Println("gradient:", gradient)
    fmt.Println("y       :", y)
    fmt.Println()
    return false
  }
  hook_bfgs := func(x, gradient ConstVector, y ConstScalar) bool {
    fmt.Fprintf(fp2, "%s\n", x.Table())
    fmt.Println("x       :", x)
    fmt.Println("gradient:", gradient)
    fmt.Println("y       :", y)
    fmt.Println()
    return false
  }
  hook_newton := func(x, gradient ConstVector, hessian ConstMatrix, y ConstScalar) bool {
    fmt.Fprintf(fp3, "%s\n", x.Table())
    fmt.Println("x       :", x)
    fmt.Println("gradient:", gradient)
    fmt.Println("y       :", y)
    fmt.Println()
    return false
  }

  x0 := NewDenseFloat64Vector([]float64{-0.5, 2})

  rprop.Run(f, x0, 0.05, []float64{1.2, 0.8},
    rprop.Hook{hook_rprop},
    rprop.Epsilon{1e-10})

  bfgs.Run(f, x0,
    bfgs.Hook{hook_bfgs},
    bfgs.Epsilon{1e-10})

  newton.RunMin(f, x0,
    newton.HookMin{hook_newton},
    newton.Epsilon{1e-8},
    newton.HessianModification{"LDL"})
```
![Gradient descent](demo/rosenbrock/rosenbrock.png)

### Constrained optimization

Maximize the function *f(x, y) = x + y* subject to *x^2 + y^2 = 1* by finding the critical point of the corresponding Lagrangian

```go
  import . "github.com/pbenner/autodiff"
  import   "github.com/pbenner/autodiff/algorithm/newton"
  import . "github.com/pbenner/autodiff/simple"

  z := NullReal64()
  t := NullReal64()
  // define the Lagrangian
  f := func(x_ ConstVector) (MagicScalar, error) {
    // z = x + y + lambda(x^2 + y^2 - 1)
    x      := x_.ConstAt(0)
    y      := x_.ConstAt(1)
    lambda := x_.ConstAt(2)
    z.Reset()
    t.Pow(x, ConstFloat64(2.0))
    z.Add(z, t)
    t.Pow(y, ConstFloat64(2.0))
    z.Add(z, t)
    z.Sub(z, ConstFloat64(1.0))
    z.Mul(z, lambda)
    z.Add(z, y)
    z.Add(z, x)

    return z, nil
  }
  // initial value
  x0    := NewDenseFloat64Vector([]float64{3,  5, 1})
  // run Newton's method
  xn, _ := newton.RunCrit(
      f, x0,
      newton.Epsilon{1e-8})
```
