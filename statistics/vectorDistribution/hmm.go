/* Copyright (C) 2017-2020 Philipp Benner
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

package vectorDistribution

/* -------------------------------------------------------------------------- */

import   "fmt"
import   "bytes"

import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/generic"

import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

type Hmm struct {
  generic.Hmm
  Edist []ScalarPdf
}

/* -------------------------------------------------------------------------- */

func NewHmm(pi Vector, tr Matrix, stateMap []int, edist []ScalarPdf) (*Hmm, error) {
  p, err := generic.NewHmmProbabilityVector(pi, false); if err != nil {
    return nil, err
  }
  t, err := generic.NewHmmTransitionMatrix(tr, false); if err != nil {
    return nil, err
  }
  if hmm, err := generic.NewHmm(p, t, stateMap); err != nil {
    return nil, err
  } else {
    if len(edist) == 0 {
      edist = make([]ScalarPdf, hmm.NEDists())
    } else {
      if hmm.NEDists() != len(edist) {
        return nil, fmt.Errorf("invalid number of emission distributions")
      }
    }
    return &Hmm{*hmm, edist}, nil
  }
}

/* -------------------------------------------------------------------------- */

func (obj *Hmm) Clone() *Hmm {
  edist := make([]ScalarPdf, len(obj.Edist))
  for i := 0; i < len(obj.Edist); i++ {
    if obj.Edist[i] != nil {
      edist[i] = obj.Edist[i].CloneScalarPdf()
    }
  }
  return &Hmm{*obj.Hmm.Clone(), edist}
}

func (obj *Hmm) CloneVectorPdf() VectorPdf {
  return obj.Clone()
}

/* -------------------------------------------------------------------------- */

func (obj *Hmm) Dim() int {
  return -1
}

func (obj *Hmm) LogPdf(r Scalar, x ConstVector) error {
  return obj.Hmm.LogPdf(r, HmmDataRecord{obj.Edist, x})
}

func (obj *Hmm) Posterior(r Scalar, x ConstVector, states [][]int) error {
  return obj.Hmm.Posterior(r, HmmDataRecord{obj.Edist, x}, states)
}

func (obj *Hmm) PosteriorMarginals(x ConstVector) ([]Vector, error) {
  return obj.Hmm.PosteriorMarginals(HmmDataRecord{obj.Edist, x})
}

func (obj *Hmm) Viterbi(x ConstVector) ([]int, error) {
  return obj.Hmm.Viterbi(HmmDataRecord{obj.Edist, x})
}

/* -------------------------------------------------------------------------- */

func (obj *Hmm) GetParameters() Vector {
  p := obj.Hmm.GetParameters()
  for i := 0; i < obj.NEDists(); i++ {
    p = p.AppendVector(obj.Edist[i].GetParameters())
  }
  return p
}

func (obj *Hmm) SetParameters(parameters Vector) error {
  n := obj.Hmm.GetParameters().Dim()
  if err := obj.Hmm.SetParameters(parameters.Slice(0,n)); err != nil {
    return err
  }
  parameters = parameters.Slice(n,parameters.Dim())
  if parameters.Dim() > 0 {
    for i := 0; i < obj.NEDists(); i++ {
      n := obj.Edist[i].GetParameters().Dim()
      if err := obj.Edist[i].SetParameters(parameters.Slice(0,n)); err != nil {
        return err
      }
      parameters = parameters.Slice(n, parameters.Dim())
    }
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func (obj *Hmm) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, obj.Hmm.String())
  fmt.Fprintf(&buffer, "Emissions:\n")
  for i := 0; i < obj.NEDists(); i++ {
    fmt.Fprintf(&buffer, "-> %+v\n", obj.Edist[i].GetParameters())
  }
  return buffer.String()
}

/* -------------------------------------------------------------------------- */

func (obj *Hmm) ImportConfig(config ConfigDistribution, t ScalarType) error {

  if err := obj.Hmm.ImportConfig(config, t); err != nil {
    return err
  }

  distributions := make([]ScalarPdf, len(config.Distributions))
  for i := 0; i < len(config.Distributions); i++ {
    if tmp, err := ImportScalarPdfConfig(config.Distributions[i], t); err != nil {
      return err
    } else {
      distributions[i] = tmp
    }
  }
  obj.Edist = distributions

  return nil
}

func (obj *Hmm) ExportConfig() ConfigDistribution {

  distributions := make([]ConfigDistribution, len(obj.Edist))
  for i := 0; i < len(obj.Edist); i++ {
    distributions[i] = obj.Edist[i].ExportConfig()
  }
  config := obj.Hmm.ExportConfig()
  config.Name = "vector:hmm distribution"
  config.Distributions = distributions

  return config
}
