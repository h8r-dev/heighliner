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
	DaggerConstraint = "=0.2.10"
	// DaggerDefault is the default version of dagger to download.
	DaggerDefault = "0.2.10"
)

const (
	// NhctlConstraint defines all acceptable versions of nhctl.
	NhctlConstraint = "=0.6.16"
	// NhctlDefault is the default nhctl version to download.
	NhctlDefault = "0.6.16"
)

const (
	// TerraformConstraint defines all acceptable versions of terrafform.
	TerraformConstraint = "=1.1.9"
	// TerraformDefault is the default terraform version to use.
	TerraformDefault = "1.1.9"
)
