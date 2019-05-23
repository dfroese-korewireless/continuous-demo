import jetbrains.buildServer.configs.kotlin.v2018_2.*
import jetbrains.buildServer.configs.kotlin.v2018_2.buildSteps.ExecBuildStep
import jetbrains.buildServer.configs.kotlin.v2018_2.buildSteps.dockerCommand
import jetbrains.buildServer.configs.kotlin.v2018_2.buildSteps.exec
import jetbrains.buildServer.configs.kotlin.v2018_2.triggers.vcs

version = "2018.2"

project {
    description = "A sample project for experimenting with TeamCity Kotlin DSL"

    buildType(Default)
}

object Default : BuildType({
    name = "Default"
    buildNumberPattern = "1.0.%build.counter%"

    vcs {
        root(DslContext.settingsRoot)
    }

    steps {
				exec {
					name = "Run Build Script"
					path = "./scripts/build.sh"
					dockerImage = "golang:1.12.4"
					dockerImagePlatform = ExecBuildStep.ImagePlatform.Linux
					dockerRunParameters = "--rm"
				}
				exec {
					name = "Run Test Script"
					path = "./scripts/test.sh"
					dockerImage = "golang:1.12.4"
					dockerImagePlatform = ExecBuildStep.ImagePlatform.Linux
					dockerRunParameters = "--rm"
				}
        dockerCommand {
            name = "Build docker image"
            commandType = build {
                source = path {
                    path = "dockerfile"
                }
                namesAndTags = "continuous-demo"
            }
        }
				dockerCommand {
					name = "Create image tar"
					commandType = other {
						subCommand = "save"
						commandArgs = "-o demo-docker-image-%env.BUILD_NUMBER%.tar.gz continuous-demo"
					}
				}
				step {
					name = "TeamCity Veracode Test"
					type = "teamcity-veracode-plugin"
					param("include", "app-demo-%env.BUILD_NUMBER%.tar.gz, demo-docker-image-%env.BUILD_NUMBER%.tar.gz")
					param("appName", "TeamCity VeraCode Demo")
					param("createProfile", "true")
					param("criticality", "VeryLow")
					param("credentialsType", "Username/Password")
					param("waitForScan", "false")
					param("useGlobalCredentials", "true")
					param("createSandbox", "false")
					param("version", "Teamcity VeraCode Demo")
				}
    }

    triggers {
        vcs {
					branchFilter = """
						-:*
						+:<default>
					""".trimIndent()
					groupCheckinsByCommitter = true
        }
    }
})
