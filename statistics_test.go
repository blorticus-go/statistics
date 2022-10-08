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
