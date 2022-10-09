package stats_test

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"

	stats "github.com/blorticus-go/statistics"
)

type statisticalSampleSetTestCase struct {
	floatSet                               []float64
	shouldTestMean                         bool
	expectedMean                           float64
	shouldTestMedian                       bool
	expectedMedian                         float64
	shouldTestMode                         bool
	expectedModeLargestDistributionCount   uint
	expectedModeValuesWithThatDistribution []float64
	shouldTestRange                        bool
	expectedRange                          float64
	shouldTestVariance                     bool
	expectedSampleVariance                 float64
	expectedPopulationVariance             float64
	shouldTestStdev                        bool
	expectedSampleStdev                    float64
	expectedPopulationStdev                float64
}

func (testCase *statisticalSampleSetTestCase) RunTest() error {
	s, err := stats.MakeStatisticalSampleSetFrom(testCase.floatSet)
	if err != nil {
		return fmt.Errorf("on MakeStatisticalSampleSetFrom() got error: %s", err.Error())
	}

	if testCase.shouldTestMean {
		gotMean := s.Mean()
		if gotMean != testCase.expectedMean {
			return fmt.Errorf("expected Mean (%f), got (%f)", testCase.expectedMean, gotMean)
		}
	}

	if testCase.shouldTestMedian {
		gotMedian := s.Median()
		if gotMedian != testCase.expectedMedian {
			return fmt.Errorf("expected Median (%f), got (%f)", testCase.expectedMedian, gotMedian)
		}
	}

	if testCase.shouldTestMode {
		frequencyCount, valuesWithThatCount := s.Mode()

		if frequencyCount != testCase.expectedModeLargestDistributionCount {
			return fmt.Errorf("expected that the largest distribution count = (%d), got (%d)", testCase.expectedModeLargestDistributionCount, frequencyCount)
		}

		if len(valuesWithThatCount) != len(testCase.expectedModeValuesWithThatDistribution) {
			return fmt.Errorf("expected that (%d) values matched most frequent count but (%d) did", len(testCase.expectedModeValuesWithThatDistribution), len(valuesWithThatCount))
		}

		sort.Float64s(valuesWithThatCount)
		sort.Float64s(testCase.expectedModeValuesWithThatDistribution)

		for i, expectedValue := range testCase.expectedModeValuesWithThatDistribution {
			gotValue := valuesWithThatCount[i]

			if expectedValue != gotValue {
				return fmt.Errorf("in mode set, expected (%f), got (%f)", expectedValue, gotValue)
			}
		}

	}

	if testCase.shouldTestRange {
		gotRange := s.Range()
		if gotRange != testCase.expectedRange {
			return fmt.Errorf("expected Range (%f), got (%f)", testCase.expectedRange, gotRange)
		}
	}

	if testCase.shouldTestVariance {
		sampleVariance := s.SampleVariance()
		populationVariance := s.PopulationVariance()

		if math.IsNaN(testCase.expectedSampleVariance) {
			if !math.IsNaN(sampleVariance) {
				return fmt.Errorf("expected sample variance (NaN), got (%f)", sampleVariance)
			}
		} else {
			if sampleVariance != testCase.expectedSampleVariance {
				return fmt.Errorf("expected sample variance (%f), got (%f)", testCase.expectedSampleVariance, sampleVariance)
			}
		}

		if populationVariance != testCase.expectedPopulationVariance {
			return fmt.Errorf("expected population variance (%f), got (%f)", testCase.expectedPopulationVariance, populationVariance)
		}
	}

	if testCase.shouldTestStdev {
		sampleStdev := s.SampleStdev()
		populationStdev := s.PopulationStdev()

		if math.IsNaN(testCase.expectedSampleStdev) {
			if !math.IsNaN(sampleStdev) {
				return fmt.Errorf("expected sample stdev (NaN), got (%f)", sampleStdev)
			}
		} else {
			if sampleStdev != testCase.expectedSampleStdev {
				return fmt.Errorf("expected sample stdev (%f), got (%f)", testCase.expectedSampleStdev, sampleStdev)
			}
		}

		if populationStdev != testCase.expectedPopulationStdev {
			return fmt.Errorf("expected population stdev (%f), got (%f)", testCase.expectedPopulationStdev, populationStdev)
		}
	}

	return nil
}

func TestStatisticalSets(t *testing.T) {
	for testIndex, testCase := range []*statisticalSampleSetTestCase{
		{
			floatSet:                               []float64{1.0},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			shouldTestVariance:                     true,
			shouldTestStdev:                        true,
			expectedMean:                           1.0,
			expectedMedian:                         1.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{1.0},
			expectedRange:                          0.0,
			expectedSampleVariance:                 math.NaN(),
			expectedPopulationVariance:             0,
			expectedSampleStdev:                    math.NaN(),
			expectedPopulationStdev:                0,
		},
		{
			floatSet:                               []float64{-1.0},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			shouldTestVariance:                     true,
			shouldTestStdev:                        true,
			expectedMean:                           -1.0,
			expectedMedian:                         -1.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{-1.0},
			expectedRange:                          0.0,
			expectedSampleVariance:                 math.NaN(),
			expectedPopulationVariance:             0,
			expectedSampleStdev:                    math.NaN(),
			expectedPopulationStdev:                0,
		},
		{
			floatSet:                               []float64{0},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			shouldTestVariance:                     true,
			shouldTestStdev:                        true,
			expectedMean:                           0.0,
			expectedMedian:                         0.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{0},
			expectedRange:                          0.0,
			expectedSampleVariance:                 math.NaN(),
			expectedPopulationVariance:             0,
			expectedSampleStdev:                    math.NaN(),
			expectedPopulationStdev:                0,
		},
		{
			floatSet:                               []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           4.0,
			expectedMedian:                         4.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0},
			expectedRange:                          6.0,
		},
		{
			floatSet:                               []float64{-1.0, -2.0, -3.0, -4.0, -5.0, -6.0, -7.0},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           -4.0,
			expectedMedian:                         -4.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{-7.0, -6.0, -5.0, -4.0, -3.0, -2.0, -1.0},
			expectedRange:                          6.0,
		},
		{
			floatSet:                               []float64{1.0, 2.0, 3.0, 4.0, -1.0, -2.0, -3.0, -4.0},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           0.0,
			expectedMedian:                         0.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{-4.0, -3.0, -2.0, -1.0, 1.0, 2.0, 3.0, 4.0},
			expectedRange:                          8.0,
		},
		{
			floatSet:                               []float64{1.0, 2.0, 3.0, 4.0, -1.0, -2.0, -3.0},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           0.5714285714285714,
			expectedMedian:                         1.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{-3.0, -2.0, -1.0, 1.0, 2.0, 3.0, 4.0},
			expectedRange:                          7.0,
		},
		{
			floatSet:                               []float64{1, 3},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           2.0,
			expectedMedian:                         2.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{1, 3},
			expectedRange:                          2.0,
		},
		{
			floatSet:                               []float64{1, 3, 5},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           3.0,
			expectedMedian:                         3.0,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{1, 3, 5},
			expectedRange:                          4.0,
		},
		{
			floatSet:                               []float64{3.45, -0.22, 0, 2.5, 1000.5, -30.9875646},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           162.5404059,
			expectedMedian:                         1.25,
			expectedModeLargestDistributionCount:   1,
			expectedModeValuesWithThatDistribution: []float64{-30.9875646, -0.22, 0, 2.5, 3.45, 1000.5},
			expectedRange:                          1031.4875646,
		},
		{
			floatSet:                               []float64{0, 1, -1, 5, 3, 1, 15, 3, 5, 1},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           3.3,
			expectedMedian:                         2.0,
			expectedModeLargestDistributionCount:   3,
			expectedModeValuesWithThatDistribution: []float64{1},
			expectedRange:                          16.0,
		},
		{
			floatSet:                               []float64{0, 1, 2, 3, 4, 5, 6, 2, 4, 6, 2, 6},
			shouldTestMean:                         true,
			shouldTestMedian:                       true,
			shouldTestMode:                         true,
			shouldTestRange:                        true,
			expectedMean:                           3.4166666666666665,
			expectedMedian:                         3.5,
			expectedModeLargestDistributionCount:   3,
			expectedModeValuesWithThatDistribution: []float64{2, 6},
			expectedRange:                          6.0,
		},
		{
			floatSet:                   []float64{100, 100, 100, 100},
			shouldTestVariance:         true,
			shouldTestStdev:            true,
			expectedSampleVariance:     0,
			expectedPopulationVariance: 0,
			expectedSampleStdev:        0,
			expectedPopulationStdev:    0,
		},
		{
			floatSet:                   []float64{46, 37, 40, 33, 42, 36, 40, 47, 34, 45},
			shouldTestVariance:         true,
			shouldTestStdev:            true,
			expectedSampleVariance:     float64(224) / float64(9),
			expectedPopulationVariance: 22.4,
			expectedSampleStdev:        math.Sqrt(float64(224) / float64(9)),
			expectedPopulationStdev:    math.Sqrt(22.4),
		},
		{
			floatSet:                   []float64{1.90, 3.00, 2.53, 3.71, 2.12, 1.76, 2.71, 1.39, 4.00, 3.33},
			shouldTestVariance:         true,
			shouldTestStdev:            true,
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
