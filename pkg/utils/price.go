package utils

func PriceDiffPercent(priceA string, priceB string) float32 {
	diffPercent := (StringToFloat32(priceA) - StringToFloat32(priceB)) / StringToFloat32(priceB) * 100
	return diffPercent
}

// if compareA is larger than compareB, then return true
func Compare(compareA string, compareB string, weightA float32, weightB float32) bool {
	// string to float32
	compareAf := StringToFloat32(compareA)
	compareBf := StringToFloat32(compareB)

	if weightA != 0 {
		compareAf = compareAf * weightA
	}

	if weightB != 0 {
		compareBf = compareBf * weightB
	}

	return compareAf > compareBf
}
