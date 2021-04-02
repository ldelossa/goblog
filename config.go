package goblog

// Version is the goblog binary version.
//
// This will increment when a new GoBlog is released.
// This supports identifying upgrades.
var Version string = "dev"

type Config struct {
	// The paths your front-end web applications serves.
	// When GoBlog encounters these paths it will serve
	// your web application's index.html
	//
	// This is how deep linking is supported.
	AppPaths []string
}
