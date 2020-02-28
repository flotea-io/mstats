package programs

var CurrentWorking = make(map[string]bool)

func IsCurrentlyWorking() bool {

	for _, isWorking := range CurrentWorking {
		if isWorking {
			return true
		}
	}
	return false
}
