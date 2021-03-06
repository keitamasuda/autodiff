/* Copyright (C) 2015-2020 Philipp Benner
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
import "encoding/json"
import "math"
import "io/ioutil"
import "os"
import "testing"

/* -------------------------------------------------------------------------- */

func TestReal(t *testing.T) {

  a := NewReal64(1.0)

  if a.GetFloat64() != 1.0 {
    t.Error("test failed")
  }

  a.Add(a, NewReal64(2.0))

  if a.GetFloat64() != 3.0 {
    t.Error("test failed")
  }
}

func TestDiff1(t *testing.T) {

  t1 := NullReal64()

  f := func(x Scalar) ConstScalar {
    return t1.Add(t1.Mul(NewReal64(2), t1.Pow(x, ConstFloat64(3))), NewReal64(4))
  }
  x := NewReal64(9)

  Variables(2, x)

  y := f(x)

  if y.GetDerivative(0) != 486 {
    t.Error("test failed")
  }
  if y.GetHessian(0, 0) != 108 {
    t.Error("test failed")
  }
}

func TestDiff2(t *testing.T) {

  f := func(x Scalar) Scalar {
    y := x.CloneScalar()
    y.Pow(y, ConstFloat64(3))
    y.Mul(y, NewReal64(2))
    y.Add(y, NewReal64(4))
    return y
  }
  x := NewReal64(9)

  Variables(2, x)

  y := f(x)

  if y.GetDerivative(0) != 486 {
    t.Error("test failed")
  }
  if y.GetHessian(0, 0) != 108 {
    t.Error("test failed")
  }
}

func TestMul(t *testing.T) {

  a := NewReal64(13.123)
  b := NewReal64( 4.321)

  Variables(2, a)

  a.Mul(a, a) // a^2
  a.Mul(a, a) // a^4
  a.Mul(a, b) // a^4 b

  if math.Abs(a.GetFloat64() - 128149.4603376) > 1e-4 {
    t.Error("test failed")
  }
  if math.Abs(a.GetDerivative(0) - 39061.025783) > 1e-4 {
    t.Error("test failed")
  }
  if math.Abs(a.GetHessian(0, 0) - 8929.5951649) > 1e-4 {
    t.Error("test failed")
  }
}

func TestPow1(t *testing.T) {
  x := NewReal64(3.4)
  k := NewReal64(4.1)

  Variables(2, x, k)

  r := NullReal64()
  r.Pow(x, k)

  if math.Abs(r.GetDerivative(0) - 182.124553) > 1e-4  ||
    (math.Abs(r.GetDerivative(1) - 184.826947) > 1e-4) {
    t.Error("test failed")
  }
  if math.Abs(r.GetHessian(0, 0) - 166.054739) > 1e-4  ||
    (math.Abs(r.GetHessian(1, 1) - 226.186676) > 1e-4) {
    t.Error("test failed")
  }
}

func TestPow2(t *testing.T) {
  x := NewReal64(-3.4)
  k := NewReal64( 4.0)

  Variables(2, x, k)

  r := NullReal64()
  r.Pow(x, k)

  if math.Abs(r.GetDerivative(0) - -157.216) > 1e-4  ||
    (math.Abs(r.GetHessian(0, 0) -  138.720) > 1e-4) {
    t.Error("test failed")
  }
  if !math.IsNaN(r.GetDerivative(1))  ||
    (!math.IsNaN(r.GetDerivative(1))) {
    t.Error("test failed")
  }
}

func TestTan(t *testing.T) {

  a := NewReal64(4.321)
  Variables(1, a)

  s := NullReal64()
  s.Tan(a)

  if math.Abs(s.GetDerivative(0) - 6.87184) > 0.0001 {
    t.Error("Incorrect derivative for Tan()!", s.GetDerivative(0))
  }
}

func TestTanh1(t *testing.T) {

  a := NewReal64(4.321)
  Variables(2, a)

  s := NullReal64()
  s.Tanh(a)

  if math.Abs(s.GetDerivative(0) -  0.00070588) > 0.0000001 {
    t.Error("test failed")
  }
  if math.Abs(s.GetHessian(0, 0) - -0.00141127) > 0.0000001 {
    t.Error("test failed")
  }
}

func TestTanh2(t *testing.T) {

  a := NewReal64(4.321)
  Variables(2, a)

  a.Tanh(a)

  if math.Abs(a.GetDerivative(0) -  0.00070588) > 0.0000001 {
    t.Error("test failed")
  }
  if math.Abs(a.GetHessian(0, 0) - -0.00141127) > 0.0000001 {
    t.Error("test failed")
  }
}

func TestErf(t *testing.T) {

  a := NewReal64(0.23)
  Variables(2, a)

  s := NullReal64()
  s.Erf(a)

  if math.Abs(s.GetDerivative(0) -  1.07023926) > 1e-6 ||
    (math.Abs(s.GetHessian(0, 0) - -0.49231006) > 1e-6) {
    t.Error("test failed")
  }
}

func TestErfc(t *testing.T) {

  a := NewReal64(0.23)
  Variables(2, a)

  s := NullReal64()
  s.Erfc(a)

  if math.Abs(s.GetDerivative(0) - -1.07023926) > 1e-6 ||
    (math.Abs(s.GetHessian(0, 0) -  0.49231006) > 1e-6) {
    t.Error("test failed")
  }
}

func TestLogErfc1(t *testing.T) {

  a := NewReal64(0.23)
  Variables(2, a)

  s := NullReal64()
  s.LogErfc(a)

  if math.Abs(s.GetDerivative(0) - -1.436606354) > 1e-6 {
    t.Error("test failed")
  }
  if math.Abs(s.GetHessian(0, 0) - -1.402998894) > 1e-6 {
    t.Error("test failed")
  }
}

func TestLogErfc2(t *testing.T) {

  a := NewReal64(0.23)
  Variables(2, a)

  a.LogErfc(a)

  if math.Abs(a.GetDerivative(0) - -1.436606354) > 1e-6 {
    t.Error("test failed")
  }
  if math.Abs(a.GetHessian(0, 0) - -1.402998894) > 1e-6 {
    t.Error("test failed")
  }
}

func TestGamma(t *testing.T) {

  a := NewReal64(4.321)
  Variables(2, a)

  s := NullReal64()
  s.Gamma(a)

  if math.Abs(s.GetDerivative(0) - 12.2353264) > 1e-6 ||
    (math.Abs(s.GetHessian(0, 0) - 18.8065398) > 1e-6) {
    t.Error("test failed")
  }
}

func TestGammaP(t *testing.T) {

  x := NewReal64(4.321)
  Variables(2, x)

  s := NullReal64()
  s.GammaP(9.125, x)

  if math.Abs(s.GetFloat64() - 0.029234) > 1e-6        ||
    (math.Abs(s.GetDerivative(0) - 0.036763) > 1e-6) ||
    (math.Abs(s.GetHessian(0, 0) - 0.032364) > 1e-6) {
    t.Error("test failed")
  }
}

func TestLogBessel(t *testing.T) {

  v := NewReal64(10.0)
  x := NewReal64(20.0)
  r := NewReal64(0.0)

  Variables(2, x)

  r.LogBesselI(v.GetFloat64(), x)

  if math.Abs(r.GetFloat64() - 15.0797) > 1e-4 {
    t.Error("test failed")
  }
  if math.Abs(r.GetDerivative(0) - 1.09804) > 1e-4 {
    t.Error("test failed")
  }
  if math.Abs(r.GetHessian(0, 0) - -0.0106002) > 1e-4 {
    t.Error("test failed")
  }

}

func TestHessian(t *testing.T) {
  x := NewReal64(1.5)
  y := NewReal64(2.5)
  k := NewReal64(3.0)

  Variables(2, x, y)

  t1 := NullReal64()
  t2 := NullReal64()
  // y = x^3 + y^3 - 3xy
  z := NullReal64()
  z.Sub(t1.Add(t1.Pow(x, k), t2.Pow(y, k)), t2.Mul(NewReal64(3.0), t2.Mul(x, y)))

  if math.Abs(z.GetHessian(0, 0) -  9) > 1e-6  ||
    (math.Abs(z.GetHessian(0, 1) - -3) > 1e-6) ||
    (math.Abs(z.GetHessian(1, 0) - -3) > 1e-6) ||
    (math.Abs(z.GetHessian(1, 1) - 15) > 1e-6) {
    t.Error("test failed")
  }
}

func TestRealJson(t *testing.T) {

  writeJson := func(filename string, obj interface{}) error {
    if f, err := os.Create(filename); err != nil {
      return err
    } else {
      b, err := json.MarshalIndent(obj, "", "  ")
      if err != nil {
        return err
      }
      if _, err := f.Write(b); err != nil {
        return err
      }
    }
    return nil
  }
  readJson := func(filename string, obj interface{}) error {
    if f, err := os.Open(filename); err != nil {
      return err
    } else {
      buffer, err := ioutil.ReadAll(f)
      if err != nil {
        return err
      }
      if err := json.Unmarshal(buffer, obj); err != nil {
        return err
      }
    }
    return nil
  }
  {
    filename := "real_test.1.json"

    r1 := NewReal64(1.5)
    r1.Alloc(1,2)
    r2 := NewReal64(0.0)

    if err := writeJson(filename, r1); err != nil {
      t.Error(err); return
    }
    if err := readJson(filename, r2); err != nil {
      t.Error(err); return
    }
    if r1.GetFloat64() != r2.GetFloat64() {
      t.Error("test failed")
    }
    os.Remove(filename)
  }
  {
    filename := "real_test.2.json"

    r1 := NewReal64(1.5)
    r1.Alloc(1,2)
    r1.SetDerivative(0, 2.3)
    r2 := NewReal64(0.0)

    if err := writeJson(filename, r1); err != nil {
      t.Error(err); return
    }
    if err := readJson(filename, r2); err != nil {
      t.Error(err); return
    }
    if r1.GetFloat64() != r2.GetFloat64() {
      t.Error("test failed")
    }
    if r1.GetDerivative(0) != r2.GetDerivative(0) {
      t.Error("test failed")
    }
    os.Remove(filename)
  }
  {
    filename := "real_test.3.json"

    r1 := NewReal64(1.5)
    r1.Alloc(1,2)
    r1.SetDerivative(0, 2.3)
    r1.SetHessian(0, 0, 3.4)
    r2 := NewReal64(0.0)

    if err := writeJson(filename, r1); err != nil {
      t.Error(err); return
    }
    if err := readJson(filename, r2); err != nil {
      t.Error(err); return
    }
    if r1.GetFloat64() != r2.GetFloat64() {
      t.Error("test failed")
    }
    if r1.GetDerivative(0) != r2.GetDerivative(0) {
      t.Error("test failed")
    }
    if r1.GetHessian(0, 0) != r2.GetHessian(0, 0) {
      t.Error("test failed")
    }
    os.Remove(filename)
  }
}

func TestSmoothMax(t *testing.T) {
  x  := NewDenseReal64Vector([]float64{-1,0,2,3,4,5})
  r  := NewReal64(0.0)
  t1 := NewReal64(0.0)
  t2 := NewReal64(0.0)

  r.SmoothMax(x, ConstFloat64(10), [2]Scalar{t1, t2})

  if math.Abs(r.GetFloat64() - 5) > 1e-4 {
    t.Error("test failed")
  }
}

func TestLogSmoothMax(t *testing.T) {
  x  := NewDenseReal64Vector([]float64{0,1,10203,3,4,30,6,7,1000,8,9,10})
  r  := NewReal64(0.0)
  t1 := NewReal64(0.0)
  t2 := NewReal64(0.0)
  t3 := NewReal64(0.0)

  r.LogSmoothMax(x, ConstFloat64(10), [3]Scalar{t1, t2, t3})

  if math.Abs(r.GetFloat64() - 10203) > 1e-4 {
    t.Error("test failed")
  }
}
