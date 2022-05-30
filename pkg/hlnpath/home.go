package hlnpath

// This helper builds paths to Helm's configuration, cache and data paths.
const lp = lazypath("hln")

// ConfigPath returns the path where Helm stores configuration.
func ConfigPath(elem ...string) string { return lp.configPath(elem...) }

// CachePath returns the path where Helm stores cached objects.
func CachePath(elem ...string) string { return lp.cachePath(elem...) }

// DataPath returns the path where Helm stores data.
func DataPath(elem ...string) string { return lp.dataPath(elem...) }
