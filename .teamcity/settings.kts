import jetbrains.buildServer.configs.kotlin.v2018_1.*
import jetbrains.buildServer.configs.kotlin.v2018_1.vcs.GitVcsRoot
import jetbrains.buildServer.configs.kotlin.v2018_1.triggers.vcs
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.script
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.exec
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.dockerCommand
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.DockerBuildStep.Source

version = "2018.1"

val gitVcs = GitVcsRoot({
    id("ContinuousDemo")
    name = "Continuous-Demo"
    url = "https://github.com/dfroese-korewireless/continuous-demo.git"
})

project {
    description = "A sample project for experimenting with TeamCity Kotlin DSL"

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
            name = "Start build container"
            scriptContent = "docker run --rm -d --name go-build-container  golang:1.10 tail -f /dev/null"
        }

        script {
            name = "Copy source code into container"
            scriptContent = "docker cp . go-build-container:/src"
        }

        script {
            name = "Setup scripts"
            scriptContent = "docker exec go-build-container chmod +x /src/scripts/build.sh"
        }

        script {
            name = "Run build script"
            scriptContent = "docker exec go-build-container /src/scripts/build.sh"
        }

        script {
            name = "Copy artiface archive out of container"
            scriptContent = "docker cp go-build-container:/artifacts/app.tar.gz ."
        }

        script {
            name = "Stop build container"
            scriptContent = "docker stop go-build-container"
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
