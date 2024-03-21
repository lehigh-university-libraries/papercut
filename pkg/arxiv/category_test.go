package arxiv_test

import (
	"reflect"
	"testing"

	"github.com/lehigh-university-libraries/papercut/pkg/arxiv"
)

func TestTransformLabels(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected map[string]string
	}{
		{
			name: "Sample Input",
			input: []byte(`<h4>cs.AI <span>(Artificial Intelligence)</span></h4>
			<h4>econ.GN <span>(General Economics)</span></h4>
			<h4>eess.SP <span>(Signal Processing)</span></h4>
			<h4>math.ST <span>(Statistics Theory)</span></h4>
			<h4>astro-ph.SR <span>(Solar and Stellar Astrophysics)</span></h4>
			<h4>gr-qc <span>(General Relativity and Quantum Cosmology)</span></h4>
			<h4>hep-ex <span>(High Energy Physics - Experiment)</span></h4>
			<h4>hep-lat <span>(High Energy Physics - Lattice)</span></h4>
			<h4>hep-ph <span>(High Energy Physics - Phenomenology)</span></h4>
			<h4>hep-th <span>(High Energy Physics - Theory)</span></h4>
			<h4>math-ph <span>(Mathematical Physics)</span></h4>
			<h4>nlin.SI <span>(Exactly Solvable and Integrable Systems)</span></h4>
			<h4>nucl-ex <span>(Nuclear Experiment)</span></h4>
			<h4>nucl-th <span>(Nuclear Theory)</span></h4>
			<h4>quant-ph <span>(Quantum Physics)</span></h4>
			<h4>q-bio.TO <span>(Tissues and Organs)</span></h4>
			<h4>q-fin.TR <span>(Trading and Market Microstructure)</span></h4>
			<h4>stat.ML <span>(Machine Learning)</span></h4>`),
			expected: map[string]string{
				"cs.AI":       "Artificial Intelligence",
				"econ.GN":     "General Economics",
				"eess.SP":     "Signal Processing",
				"math.ST":     "Statistics Theory",
				"astro-ph.SR": "Solar and Stellar Astrophysics",
				"gr-qc":       "General Relativity and Quantum Cosmology",
				"hep-ex":      "High Energy Physics - Experiment",
				"hep-lat":     "High Energy Physics - Lattice",
				"hep-ph":      "High Energy Physics - Phenomenology",
				"hep-th":      "High Energy Physics - Theory",
				"math-ph":     "Mathematical Physics",
				"nlin.SI":     "Exactly Solvable and Integrable Systems",
				"nucl-ex":     "Nuclear Experiment",
				"nucl-th":     "Nuclear Theory",
				"quant-ph":    "Quantum Physics",
				"q-bio.TO":    "Tissues and Organs",
				"q-fin.TR":    "Trading and Market Microstructure",
				"stat.ML":     "Machine Learning",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := arxiv.TransformLabels(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Test %q: expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}
