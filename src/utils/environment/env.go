package environment

type CompilationSettings struct {
	SDKPath string
}

var CompilerEnvironment *CompilationSettings

func SetCompilationSettings(sdkPath string) {
	CompilerEnvironment = &CompilationSettings{
		SDKPath: sdkPath,
	}
}
