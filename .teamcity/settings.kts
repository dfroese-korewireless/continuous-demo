import jetbrains.buildServer.configs.kotlin.v2018_1.*
import jetbrains.buildServer.configs.kotlin.v2018_1.triggers.vcs

version = "2018.1"

project {
    buildType(Default)

    params {
        param("system.BuildNumber", "1.0.%build.counter%")
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
        dockerCommand {
            name = "Build docker image"
            commandType = build {
                source = path {
                    path = "dockerfile"
                }
                namesAndTags = "continuous-demo"
            }
        }

        script {
            name = "Run images"
            scriptContent = "./deploy"
        }
    }
})
