import jetbrains.buildServer.configs.kotlin.v2018_1.*
import jetbrains.buildServer.configs.kotlin.v2018_1.triggers.vcs
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.script
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.exec
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.dockerCommand
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.DockerBuildStep.Source

version = "2018.1"

project {
    buildType(Default)

    params {
        param("env.BuildNumber", "1.0.%build.counter%")
    }
}

object Default : BuildType({
    name = "Default"

    vcs {
        root(DslContext.settingsRoot)
    }

    triggers {
        vcs {
        }
    }

    steps {
        exec {
            name = "Diagnostic check"
            path = "./scripts/diagnostic.sh"
        }

        script {
            name = "Build config file"
            scriptContent = "sed -e 's/##version##/%env.BuildNumber%/g' appsettings.dev.json > appsettings.json"
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

        exec {
            name = "Run images"
            path = "./scripts/deploy.sh"
        }
    }
})
