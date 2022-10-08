package stats

import "sync"

type valueDistributionTracker struct {
	mutex                     sync.Mutex
	mapHasBeenGenerated       bool
	sourceValueSet            []float64
	conditionallyGeneratedMap map[float64]uint
}

func newValueDistributionTracker(valueSet []float64) *valueDistributionTracker {
	return &valueDistributionTracker{
		mapHasBeenGenerated:       false,
		sourceValueSet:            valueSet,
		conditionallyGeneratedMap: make(map[float64]uint),
	}
}

func (generator *valueDistributionTracker) Map() map[float64]uint {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	if !generator.mapHasBeenGenerated {
		for _, v := range generator.sourceValueSet {
			if currentValueCount, valueIsInMap := generator.conditionallyGeneratedMap[v]; valueIsInMap {
				generator.conditionallyGeneratedMap[v] = currentValueCount + uint(1)
			} else {
				generator.conditionallyGeneratedMap[v] = uint(1)
			}
		}

		generator.mapHasBeenGenerated = true
	}

	return generator.conditionallyGeneratedMap
}

// a modal map is the inverse of a value distribution map.  That is, it is keyed by the number
// of occurances of a value and points to a list of values that occurred that number of times.
type modalTracker struct {
	mutex                     sync.Mutex
	mapHasBeenGenerated       bool
	valueDistributionTracker  *valueDistributionTracker
	conditionallyGeneratedMap map[uint][]float64
	highestFrequencyCount     uint
}

func newModalTracker(using *valueDistributionTracker) *modalTracker {
	return &modalTracker{
		valueDistributionTracker:  using,
		conditionallyGeneratedMap: make(map[uint][]float64),
	}
}

func generateModalMapFromADistributionMap(distributionMap map[float64]uint) (modalMap map[uint][]float64, highestOccuranceCount uint) {
	modalMap = make(map[uint][]float64)
	highestOccuranceCount = uint(0)

	for value, countOfTimesValueWasSeen := range distributionMap {
		if listOfValuesSeenThisManyTimes, thisCountIsInMap := modalMap[countOfTimesValueWasSeen]; thisCountIsInMap {
			modalMap[countOfTimesValueWasSeen] = append(listOfValuesSeenThisManyTimes, value)
		} else {
			listOfValuesSeenThisManyTimes = make([]float64, 1)
			listOfValuesSeenThisManyTimes[0] = value
			modalMap[countOfTimesValueWasSeen] = listOfValuesSeenThisManyTimes
		}

		if countOfTimesValueWasSeen > highestOccuranceCount {
			highestOccuranceCount = countOfTimesValueWasSeen
		}
	}

	return modalMap, highestOccuranceCount
}

func (generator *modalTracker) Modes() (numberOfTimesValuesWereSeen uint, valuesSeenThatManyTime []float64) {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	if !generator.mapHasBeenGenerated {
		generator.conditionallyGeneratedMap, generator.highestFrequencyCount = generateModalMapFromADistributionMap(generator.valueDistributionTracker.Map())
	}

	return generator.highestFrequencyCount, generator.conditionallyGeneratedMap[generator.highestFrequencyCount]
}
