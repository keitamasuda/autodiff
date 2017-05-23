/* Copyright (C) 2016 Philipp Benner
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

package distribution

/* -------------------------------------------------------------------------- */

//import   "fmt"
import   "math"
import   "testing"

import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

func TestBeta1(t *testing.T) {
  d, _ := NewBetaDistribution(NewReal(3), NewReal(5), false)
  x := NewVector(RealType, []float64{0.8})
  r := NewReal(0.0)

  if err := d.LogPdf(r, x); err != nil {
    t.Error(err)
  }
  if math.Abs(r.GetValue() - -2.230078e+00) > 1e-4 {
    t.Error("test failed")
  }
}

func TestBeta2(t *testing.T) {
  d, _ := NewBetaDistribution(NewReal(3), NewReal(5), true)
  x := NewVector(RealType, []float64{math.Log(0.8)})
  r := NewReal(0.0)

  if err := d.LogPdf(r, x); err != nil {
    t.Error(err)
  }
  if math.Abs(r.GetValue() - -2.230078e+00) > 1e-4 {
    t.Error("test failed")
  }
}
