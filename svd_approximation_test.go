package image_compression

import (
	"gonum.org/v1/gonum/mat"
	"math"
	"testing"
)

const approximateError = 100000

func equalWithApproximateError(actual *mat.Dense, expected *mat.Dense) bool {
	expectedData := expected.RawMatrix().Data
	for index, value := range actual.RawMatrix().Data {
		if math.Floor(value*approximateError)/approximateError != expectedData[index] {
			return false
		}
	}
	return true
}

var testApproximationMatrix = []struct {
	rank           int
	inputMatrix    *mat.Dense
	expectedMatrix *mat.Dense
}{
	{
		2,
		mat.NewDense(4, 3, []float64{12, 23, 323, 432, 53, 63, 7232, 82, 91, 23, 2, 121}),
		mat.NewDense(4, 3, []float64{11.94797, 27.90595, 322.57751, 432.45255, 10.32815, 66.67474, 7231.97331, 84.51631, 90.78330, 22.90880, 10.59927, 120.25946}),
	},
	{
		3,
		mat.NewDense(4, 3, []float64{12, 23, 323, 432, 53, 63, 7232, 82, 91, 23, 2, 121}),
		mat.NewDense(4, 3, []float64{12.00000, 22.99999, 322.99999, 431.99999, 52.99999, 62.99999, 7231.99999, 81.99999, 90.99999, 22.99999, 1.99999, 121}),
	},
}

func TestApproximation(t *testing.T) {
	for _, tc := range testApproximationMatrix {
		actual := approximate(tc.inputMatrix, tc.rank)
		expected := tc.expectedMatrix

		if !equalWithApproximateError(&actual, expected) {
			t.Errorf("Matrix %v expected, got %v", &expected, actual)
		}
	}
}
