package version

const (
	// DevelopmentVersion may contain bugs.
	DevelopmentVersion = "dev"
)

var (
	// Version holds the complete version number. Filled in at linking time.
	Version = DevelopmentVersion

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""
)

const (
	// DaggerConstraint defines all acceptable versions of dagger.
	DaggerConstraint = ">=0.2.5, <=0.2.6"
	// DaggerDefault is the default version of dagger to download.
	DaggerDefault = "0.2.6"
)
