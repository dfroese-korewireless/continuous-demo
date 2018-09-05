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

        script {
            name = "Write VCS information"
            scriptContent = "echo $DslContext.settingsRoot && echo $VcsRoot.name"
        }

        // script {
        //     name = "Start build container"
        //     scriptContent = "docker run --rm -d --name dotnet-build-contanier -v /opt/buildagent/artifacts:/artifacts microsoft/dotnet:2.1-sdk /bin/bash"
        // }

        // dockerCommand {
        //     name = "Build docker image"
        //     commandType = build {
        //         source = path {
        //             path = "dockerfile"
        //         }
        //         namesAndTags = "continuous-demo"
        //     }
        // }

        // exec {
        //     name = "Run images"
        //     path = "./scripts/deploy.sh"
        // }

        // script {
        //     name = "Stop build container"
        //     scriptContent = "docker stop dotnet-build-container"
        // }
    }
})
