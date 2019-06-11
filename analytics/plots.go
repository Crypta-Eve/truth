package analytics

import (
	"fmt"

	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func GenerateSingleValuePlot(name string, values PairList) error {

	p, err := plot.New()
	if err != nil {
		return errors.Wrap(err, "Failed to create new plot")
	}

	p.Title.Text = name

	w := vg.Points(10)

	lossNums := make(plotter.Values, 30)
	names := []string{}

	for i := 0; i < 30; i++ {
		lossNums[i] = float64(values[i].Value)
		names = append(names, values[i].Key)
	}

	bc, err := plotter.NewBarChart(lossNums, w)

	p.Add(bc)

	p.NominalX(names...)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, fmt.Sprintf("./%v.png", name)); err != nil {
		return errors.Wrap(err, "Failed to save bar chart - "+name)
	}

	return nil
}
