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
	return medianOfAFloatSet(set.valuesSortedInAscendingOrder).computedMedian
}

func (set *StatisticalSampleSet) Mode() (modeFrequencyCount uint, valuesSeenThatManyTimes []float64) {
	return set.modeTracker.Modes()
}

func (set *StatisticalSampleSet) Range() float64 {
	lastElementIndex := len(set.valuesSortedInAscendingOrder) - 1

	return set.valuesSortedInAscendingOrder[lastElementIndex] - set.valuesSortedInAscendingOrder[0]
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

func (set *StatisticalSampleSet) InterQuartileRange() (q1 float64, q3 float64, iqr float64) {
	switch len(set.valuesSortedInAscendingOrder) {
	case 1:
		return set.valuesSortedInAscendingOrder[0], set.valuesSortedInAscendingOrder[0], 0.0

	case 2:
		return set.valuesSortedInAscendingOrder[0], set.valuesSortedInAscendingOrder[1], set.valuesSortedInAscendingOrder[1] - set.valuesSortedInAscendingOrder[0]
	}

	q2MedianInfo := medianOfAFloatSet(set.valuesSortedInAscendingOrder)

	var q1MedianInfo, q3MedianInfo *medianValueAndBracketInfo

	if q2MedianInfo.medianIsBetweenTwoValues {
		q1MedianInfo = medianOfAFloatSet(set.valuesSortedInAscendingOrder[0:q2MedianInfo.indexOfMedianRightBracketInSet])
		q3MedianInfo = medianOfAFloatSet(set.valuesSortedInAscendingOrder[q2MedianInfo.indexOfMedianRightBracketInSet:])
	} else {
		q1MedianInfo = medianOfAFloatSet(set.valuesSortedInAscendingOrder[0:q2MedianInfo.indexOfMedianInSet])
		q3MedianInfo = medianOfAFloatSet(set.valuesSortedInAscendingOrder[q2MedianInfo.indexOfMedianInSet+1:])
	}

	iqr = q3MedianInfo.computedMedian - q1MedianInfo.computedMedian

	return q1MedianInfo.computedMedian, q3MedianInfo.computedMedian, iqr
}

func thereIsAnOddNumberOfSamplesInTheSet(set []float64) bool {
	return len(set)&1 != 0
}

type medianValueAndBracketInfo struct {
	computedMedian                 float64
	medianIsBetweenTwoValues       bool
	indexOfMedianInSet             int
	indexOfMedianLeftBracketInSet  int
	indexOfMedianRightBracketInSet int
}

func medianOfAFloatSet(set []float64) *medianValueAndBracketInfo {
	midPoint := len(set) / 2

	if thereIsAnOddNumberOfSamplesInTheSet(set) {
		return &medianValueAndBracketInfo{
			computedMedian:           set[midPoint],
			medianIsBetweenTwoValues: false,
			indexOfMedianInSet:       midPoint,
		}
	}

	leftMidpointValue := set[midPoint-1]
	rightMidpointValue := set[midPoint]

	return &medianValueAndBracketInfo{
		computedMedian:                 (rightMidpointValue + leftMidpointValue) / 2.0,
		medianIsBetweenTwoValues:       true,
		indexOfMedianLeftBracketInSet:  midPoint - 1,
		indexOfMedianRightBracketInSet: midPoint,
	}
}
