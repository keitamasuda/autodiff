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

import   "fmt"

import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

type CategoricalDistribution struct {
  Theta Vector
  t     Scalar
}

/* -------------------------------------------------------------------------- */

func NewCategoricalDistribution(theta_ Vector) (*CategoricalDistribution, error) {
  if theta_.Dim() == 0 {
    return nil, fmt.Errorf("theta has invalid length")
  }
  t     := theta_.ElementType()
  theta := NullVector(t, theta_.Dim())

  for i := 0; i < theta.Dim(); i++ {
    if theta_.At(i).GetValue() < 0 {
      return nil, fmt.Errorf("invalid negative probability")
    }
    theta.At(i).Log(theta_.At(i))
  }
  result := CategoricalDistribution{
    Theta: theta,
    t    : theta.At(0).CloneScalar() }

  return &result, nil

}

/* -------------------------------------------------------------------------- */

func (dist *CategoricalDistribution) Clone() *CategoricalDistribution {
  return &CategoricalDistribution{
    Theta : dist.Theta.CloneVector(),
    t     : dist.t    .CloneScalar() }
}

/* -------------------------------------------------------------------------- */

func (dist *CategoricalDistribution) ScalarType() ScalarType {
  return dist.Theta.ElementType()
}

func (dist *CategoricalDistribution) LogPdf(r Scalar, x Scalar) error {
  r.Set(dist.Theta.At(int(x.GetValue())))
  return nil
}

func (dist *CategoricalDistribution) Pdf(r Scalar, x Scalar) error {
  if err := dist.LogPdf(r, x); err != nil {
    return err
  }
  r.Exp(r)
  return nil
}

func (dist *CategoricalDistribution) LogCdf(r Scalar, x Scalar) error {
  r.Reset()
  for i := 0; i <= int(x.GetValue()); i++ {
    r.LogAdd(r, dist.Theta.At(i), dist.t)
  }
  return nil
}

func (dist *CategoricalDistribution) Cdf(r Scalar, x Scalar) error {
  if err := dist.LogCdf(r, x); err != nil {
    return err
  }
  r.Exp(r)
  return nil
}

/* -------------------------------------------------------------------------- */

func (dist *CategoricalDistribution) GetParameters() Vector {
  return dist.Theta
}

func (dist *CategoricalDistribution) SetParameters(parameters Vector) error {
  dist.Theta.Set(parameters)
  return nil
}
