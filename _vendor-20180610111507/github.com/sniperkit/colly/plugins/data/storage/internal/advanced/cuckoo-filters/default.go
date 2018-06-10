package cuckoofilters

const (
	// default max keys 100.000
	defaultMaxKeys = 100000

	// 4 entries per bucket is suggested by the paper in section 5.1, "Optimal bucket size"
	defaultEntriesPerBucket = 4

	// With 4 entries per bucket, we can expect up to 95% load factor
	defaultLoadFactor = 0.95

	// Length of fingerprints in bits
	defaultFingerprintsBits = 16

	// Arbitrarily chosen value
	defaultMaxDisplacements = 500
)
