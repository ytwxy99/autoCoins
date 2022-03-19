package utils

func PriceDiffPercent(priceA string, priceB string) float32 {
	diffPercent := (StringToFloat32(priceA) - StringToFloat32(priceB))/StringToFloat32(priceB) * 100
	return diffPercent
}
