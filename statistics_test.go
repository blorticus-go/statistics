package stats_test

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"

	stats "github.com/blorticus-go/statistics"
)

type meanTestCase struct {
	floatSet     []float64
	expectedMean float64
}

func (testCase *meanTestCase) RunTest() error {
	s, err := stats.MakeStatisticalSampleSetFrom(testCase.floatSet)
	if err != nil {
		return fmt.Errorf("on MakeStatisticalSampleSetFrom() got error: %s", err.Error())
	}

	mean := s.Mean()

	if mean != testCase.expectedMean {
		return fmt.Errorf("expected mean = (%f), got mean = (%f)", testCase.expectedMean, mean)
	}

	return nil
}

func TestMean(t *testing.T) {
	for testIndex, testCase := range []*meanTestCase{
		{
			floatSet:     []float64{1.0},
			expectedMean: 1.0,
		},
		{
			floatSet:     []float64{-1.0},
			expectedMean: -1.0,
		},
		{
			floatSet:     []float64{0},
			expectedMean: 0.0,
		},
		{
			floatSet:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0},
			expectedMean: 4.0,
		},
		{
			floatSet:     []float64{-1.0, -2.0, -3.0, -4.0, -5.0, -6.0, -7.0},
			expectedMean: -4.0,
		},
		{
			floatSet:     []float64{1.0, 2.0, 3.0, 4.0, -1.0, -2.0, -3.0, -4.0},
			expectedMean: 0.0,
		},
		{
			floatSet:     []float64{1.0, 2.0, 3.0, 4.0, -1.0, -2.0, -3.0},
			expectedMean: 0.5714285714285714,
		},
	} {
		if err := testCase.RunTest(); err != nil {
			t.Errorf("on Mean() test index (%d): %s", testIndex, err.Error())
		}
	}
}

type medianTestCase struct {
	floatSet       []float64
	expectedMedian float64
}

func (testCase *medianTestCase) RunTest() error {
	s, err := stats.MakeStatisticalSampleSetFrom(testCase.floatSet)
	if err != nil {
		return fmt.Errorf("on MakeStatisticalSampleSetFrom, got error = (%s)", err.Error())
	}

	median := s.Median()

	if median != testCase.expectedMedian {
		return fmt.Errorf("expected median (%f), got (%f)", testCase.expectedMedian, median)
	}

	return nil
}

func TestMedian(t *testing.T) {
	for testIndex, testCase := range []*medianTestCase{
		{
			floatSet:       []float64{0},
			expectedMedian: 0.0,
		},
		{
			floatSet:       []float64{1, 3},
			expectedMedian: 2.0,
		},
		{
			floatSet:       []float64{1, 3, 5},
			expectedMedian: 3.0,
		},
		{
			floatSet:       []float64{3.45, -0.22, 0, 2.5, 1000.5, -30.9875646},
			expectedMedian: 1.25,
		},
	} {
		if err := testCase.RunTest(); err != nil {
			t.Errorf("on test with index (%d): %s", testIndex, err.Error())
		}
	}
}

type modeTestCase struct {
	floatSet                   []float64
	largestDistributionCount   uint
	valuesWithThatDistribution []float64
}

func (testCase *modeTestCase) RunTest() error {
	s, err := stats.MakeStatisticalSampleSetFrom(testCase.floatSet)
	if err != nil {
		return fmt.Errorf("on MakeStatisticalSampleSetFrom, got error = (%s)", err.Error())
	}

	frequencyCount, valuesWithThatCount := s.Mode()

	if frequencyCount != testCase.largestDistributionCount {
		return fmt.Errorf("expected that the largest distribution count = (%d), got (%d)", testCase.largestDistributionCount, frequencyCount)
	}

	if len(valuesWithThatCount) != len(testCase.valuesWithThatDistribution) {
		return fmt.Errorf("expected that (%d) values matched most frequent count but (%d) did", len(testCase.valuesWithThatDistribution), len(valuesWithThatCount))
	}

	sort.Float64s(valuesWithThatCount)
	sort.Float64s(testCase.valuesWithThatDistribution)

	for i, expectedValue := range testCase.valuesWithThatDistribution {
		gotValue := valuesWithThatCount[i]

		if expectedValue != gotValue {
			return fmt.Errorf("in mode set, expected (%f), got (%f)", expectedValue, gotValue)
		}
	}

	return nil
}

func TestMode(t *testing.T) {
	for testIndex, testCase := range []*modeTestCase{
		{
			floatSet:                   []float64{0},
			largestDistributionCount:   1,
			valuesWithThatDistribution: []float64{0},
		},
		{
			floatSet:                   []float64{0, 1, -1, 5, 3, 1, 15, 3, 5, 1},
			largestDistributionCount:   3,
			valuesWithThatDistribution: []float64{1},
		},
		{
			floatSet:                   []float64{0, 1, 2, 3, 4, 5},
			largestDistributionCount:   1,
			valuesWithThatDistribution: []float64{0, 1, 2, 3, 4, 5},
		},
		{
			floatSet:                   []float64{0, 1, 2, 3, 4, 5, 6, 2, 4, 6, 2, 6},
			largestDistributionCount:   3,
			valuesWithThatDistribution: []float64{2, 6},
		},
	} {
		if err := testCase.RunTest(); err != nil {
			t.Errorf("on test with index (%d): %s", testIndex, err.Error())
		}
	}
}

func TestMakeStatisticalSampleSetFromErrors(t *testing.T) {
	floatSet := []float64{}

	_, err := stats.MakeStatisticalSampleSetFrom(floatSet)
	if err == nil {
		t.Errorf("on MakeStatisticalSampleSetFrom() with empty set, expected error, got none")
	}

	h := math.MaxFloat64
	floatSet = []float64{h, h}

	_, err = stats.MakeStatisticalSampleSetFrom(floatSet)
	if err == nil {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding float64max to itself should generate error, but does not")
	} else if err != stats.ErrorFloat64Overflow {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding float64max to itself, expected ErrorFloat64Overflow, got err = (%s)", reflect.TypeOf(err).String())
	}

	floatSet = []float64{h, h, h, h, h, h, h}
	_, err = stats.MakeStatisticalSampleSetFrom(floatSet)
	if err == nil {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding float64max to itself seven times should generate error, but does not")
	} else if err != stats.ErrorFloat64Overflow {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding float64max to itself seven times, expected ErrorFloat64Overflow, got err = (%s)", reflect.TypeOf(err).String())
	}

	floatSet = []float64{-h, -h}

	_, err = stats.MakeStatisticalSampleSetFrom(floatSet)
	if err == nil {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding -1 * float64max to itself should generate error, but does not")
	} else if err != stats.ErrorFloat64Underflow {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding -1 * float64max to itself, expected ErrorFloat64Underflow, got err = (%s)", reflect.TypeOf(err).String())
	}

	floatSet = []float64{-h, -h, -h, -h, -h, -h, -h}

	_, err = stats.MakeStatisticalSampleSetFrom(floatSet)
	if err == nil {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding -1 * float64max to itself seven times should generate error, but does not")
	} else if err != stats.ErrorFloat64Underflow {
		t.Errorf("on MakeStatisticalSampleSetFrom, adding -1 * float64max to itself seven times, expected ErrorFloat64Underflow, got err = (%s)", reflect.TypeOf(err).String())
	}

}

type rangeTestCase struct {
	floatSet      []float64
	expectedRange float64
}

func (testCase *rangeTestCase) RunTest() error {
	s, err := stats.MakeStatisticalSampleSetFrom(testCase.floatSet)
	if err != nil {
		return fmt.Errorf("on MakeStatisticalSampleSetFrom, got error = (%s)", err.Error())
	}

	gotRange := s.Range()

	if gotRange != testCase.expectedRange {
		return fmt.Errorf("expected Range (%f), got (%f)", testCase.expectedRange, gotRange)
	}

	return nil
}

func TestRange(t *testing.T) {
	for testIndex, testCase := range []*rangeTestCase{
		{
			floatSet:      []float64{0},
			expectedRange: 0,
		},
		{
			floatSet:      []float64{1},
			expectedRange: 0,
		},
		{
			floatSet:      []float64{-1},
			expectedRange: 0,
		},
		{
			floatSet:      []float64{0, 1, 2, 3},
			expectedRange: 3,
		},
		{
			floatSet:      []float64{2, 3, 1, 0},
			expectedRange: 3,
		},
		{
			floatSet:      []float64{-6, -6, 7, 0, -3, 10, -12},
			expectedRange: 22,
		},
	} {
		if err := testCase.RunTest(); err != nil {
			t.Errorf("on test with index (%d): %s", testIndex, err.Error())
		}
	}
}

type varianceAndStdevTestCase struct {
	floatSet                   []float64
	expectedSampleVariance     float64
	expectedPopulationVariance float64
	expectedSampleStdev        float64
	expectedPopulationStdev    float64
}

func (testCase *varianceAndStdevTestCase) RunTest() error {
	s, err := stats.MakeStatisticalSampleSetFrom(testCase.floatSet)
	if err != nil {
		return fmt.Errorf("on MakeStatisticalSampleSetFrom, got error = (%s)", err.Error())
	}

	sampleVariance := s.SampleVariance()
	sampleStdev := s.SampleStdev()
	populationVariance := s.PopulationVariance()
	populationStdev := s.PopulationStdev()

	if math.IsNaN(testCase.expectedSampleVariance) {
		if !math.IsNaN(sampleVariance) {
			return fmt.Errorf("expected ample variance (NaN), got (%f)", sampleVariance)
		}
	} else {
		if sampleVariance != testCase.expectedSampleVariance {
			return fmt.Errorf("expected sample variance (%f), got (%f)", testCase.expectedSampleVariance, sampleVariance)
		}
	}

	if math.IsNaN(testCase.expectedSampleStdev) {
		if !math.IsNaN(sampleStdev) {
			return fmt.Errorf("expected sample stdev (NaN), got (%f)", sampleStdev)
		}
	} else {
		if sampleStdev != testCase.expectedSampleStdev {
			return fmt.Errorf("expected sample stdev (%f), got (%f)", testCase.expectedSampleStdev, sampleStdev)
		}
	}

	if populationVariance != testCase.expectedPopulationVariance {
		return fmt.Errorf("expected population variance (%f), got (%f)", testCase.expectedPopulationVariance, populationVariance)
	}

	if populationStdev != testCase.expectedPopulationStdev {
		return fmt.Errorf("expected population stdev (%f), got (%f)", testCase.expectedPopulationStdev, populationStdev)
	}

	return nil
}

func TestVarianceAndStdev(t *testing.T) {
	for testIndex, testCase := range []*varianceAndStdevTestCase{
		{
			floatSet:                   []float64{0},
			expectedSampleVariance:     math.NaN(),
			expectedPopulationVariance: 0,
			expectedSampleStdev:        math.NaN(),
			expectedPopulationStdev:    0,
		},
		{
			floatSet:                   []float64{10},
			expectedSampleVariance:     math.NaN(),
			expectedPopulationVariance: 0,
			expectedSampleStdev:        math.NaN(),
			expectedPopulationStdev:    0,
		},
		{
			floatSet:                   []float64{-10},
			expectedSampleVariance:     math.NaN(),
			expectedPopulationVariance: 0,
			expectedSampleStdev:        math.NaN(),
			expectedPopulationStdev:    0,
		},
		{
			floatSet:                   []float64{100, 100, 100, 100},
			expectedSampleVariance:     0,
			expectedPopulationVariance: 0,
			expectedSampleStdev:        0,
			expectedPopulationStdev:    0,
		},
		{
			floatSet:                   []float64{46, 37, 40, 33, 42, 36, 40, 47, 34, 45},
			expectedSampleVariance:     float64(224) / float64(9),
			expectedPopulationVariance: 22.4,
			expectedSampleStdev:        math.Sqrt(float64(224) / float64(9)),
			expectedPopulationStdev:    math.Sqrt(22.4),
		},
		{
			floatSet:                   []float64{1.90, 3.00, 2.53, 3.71, 2.12, 1.76, 2.71, 1.39, 4.00, 3.33},
			expectedSampleVariance:     float64(6.77185) / float64(9),
			expectedPopulationVariance: float64(6.77185) / float64(10),
			expectedSampleStdev:        math.Sqrt(float64(6.77185) / float64(9)),
			expectedPopulationStdev:    math.Sqrt(float64(6.77185) / float64(10)),
		},
	} {
		if err := testCase.RunTest(); err != nil {
			t.Errorf("on test with index (%d): %s", testIndex, err.Error())
		}
	}
}
