package math2

func Round64(val float64) int64 {
	if(val - float64(int(val)) >= 0.5) {
		return int64(val) + 1
	} else {
		return int64(val)
	}
	
	// TODO wtf????
	return 0
}

func Round32(val float32) int32 {
	if(val - float32(int(val)) >= 0.5) {
		return int32(val) + 1
	} else {
		return int32(val)
	}
	
	// TODO wtf????
	return 0
}