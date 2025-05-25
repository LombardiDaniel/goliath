package common

// Batch collects items from a channel into a slice of the specified size.
// Stops if the channel is closed or empty. Does NOT block if the channel is empty.
func Batch[T any](ch <-chan T, size uint32) []T {
	batch := make([]T, 0, size)

	for range size {
		select {
		case item, ok := <-ch:
			if !ok {
				// Channel is closed, return the batch collected so far
				return batch
			}
			batch = append(batch, item)
		default:
			// Channel has no data ready — return immediately
			return batch
		}
	}

	return batch
}

// RemoveFrom removes all occurrences of a specified item from a slice.
func RemoveFrom[T comparable](slice []T, item T) []T {
	var newSlice []T
	for _, v := range slice {
		if v != item {
			newSlice = append(newSlice, v)
		}
	}

	return newSlice
}

// IsSubset checks if all elements of the subset are present in the superset.
func IsSubset[T comparable](subset []T, superset []T) bool {
	checkMap := make(map[T]bool)
	for _, element := range superset {
		checkMap[element] = true
	}
	for _, value := range subset {
		if !checkMap[value] {
			return false // Early return if not found
		}
	}
	return true
}
