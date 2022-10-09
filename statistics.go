package stats

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type StatisticalSampleSet struct {
	valuesSortedInAscendingOrder []float64
	sumOfAllValuesInTheSet       float64
	distributionTracker          *valueDistributionTracker
	modeTracker                  *modalTracker
	varianceTracker              *varianceTracker
}

var ErrorFloat64Overflow = errors.New("float64 overflow")
var ErrorFloat64Underflow = errors.New("float64 underflow")

func MakeStatisticalSampleSetFrom(samples []float64) (*StatisticalSampleSet, error) {
	if len(samples) == 0 {
		return nil, fmt.Errorf("there must be at least one sample in the set")
	}

	copyOfSamples := make([]float64, len(samples))
	copy(copyOfSamples, samples)
	sort.Float64s(copyOfSamples)

	sum := float64(0)
	for _, v := range samples {
		sum = sum + v
	}

	if sum == math.Inf(1) {
		return nil, ErrorFloat64Overflow
	}

	if sum == math.Inf(-1) {
		return nil, ErrorFloat64Underflow
	}

	distributionTracker := newValueDistributionTracker(copyOfSamples)
	modeTracker := newModalTracker(distributionTracker)

	set := &StatisticalSampleSet{
		valuesSortedInAscendingOrder: copyOfSamples,
		sumOfAllValuesInTheSet:       sum,
		distributionTracker:          distributionTracker,
		modeTracker:                  modeTracker,
	}

	set.varianceTracker = NewVarianceTracker(copyOfSamples, set)

	return set, nil
}

func (set *StatisticalSampleSet) Minimum() float64 {
	return set.valuesSortedInAscendingOrder[0]
}

func (set *StatisticalSampleSet) Maximum() float64 {
	return set.valuesSortedInAscendingOrder[len(set.valuesSortedInAscendingOrder)-1]
}

func (set *StatisticalSampleSet) Mean() float64 {
	return set.sumOfAllValuesInTheSet / float64(len(set.valuesSortedInAscendingOrder))
}

func (set *StatisticalSampleSet) Median() float64 {
	midPoint := len(set.valuesSortedInAscendingOrder) / 2

	if set.thereAreAnOddNumberOfSamples() {
		return set.valuesSortedInAscendingOrder[midPoint]
	}

	leftMidpointValue := set.valuesSortedInAscendingOrder[midPoint-1]
	rightMidpointValue := set.valuesSortedInAscendingOrder[midPoint]

	return (rightMidpointValue + leftMidpointValue) / 2.0
}

func (set *StatisticalSampleSet) Mode() (modeFrequencyCount uint, valuesSeenThatManyTimes []float64) {
	return set.modeTracker.Modes()
}

func (set *StatisticalSampleSet) Range() float64 {
	lastElementIndex := len(set.valuesSortedInAscendingOrder) - 1

	return set.valuesSortedInAscendingOrder[lastElementIndex] - set.valuesSortedInAscendingOrder[0]
}

func (set *StatisticalSampleSet) thereAreAnOddNumberOfSamples() bool {
	return len(set.valuesSortedInAscendingOrder)&1 != 0
}

func (set *StatisticalSampleSet) SampleVariance() float64 {
	return set.varianceTracker.Variance() / (float64(len(set.valuesSortedInAscendingOrder)) - 1)
}

func (set *StatisticalSampleSet) PopulationVariance() float64 {
	return set.varianceTracker.Variance() / float64(len(set.valuesSortedInAscendingOrder))
}

func (set *StatisticalSampleSet) SampleStdev() float64 {
	return math.Sqrt(set.SampleVariance())
}

func (set *StatisticalSampleSet) PopulationStdev() float64 {
	return math.Sqrt(set.PopulationVariance())
}
