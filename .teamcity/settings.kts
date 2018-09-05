import jetbrains.buildServer.configs.kotlin.v2018_1.*
import jetbrains.buildServer.configs.kotlin.v2018_1.vcs.GitVcsRoot
import jetbrains.buildServer.configs.kotlin.v2018_1.triggers.vcs
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.script
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.exec
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.dockerCommand
import jetbrains.buildServer.configs.kotlin.v2018_1.buildSteps.DockerBuildStep.Source

version = "2018.1"

project {
    description = "A sample project for experimenting with TeamCity Kotlin DSL"

    val gitVcs = GitVcsRoot({
        id("ContinuousDemo")
        name = "Continuous-Demo"
        url = "https://github.com/dfroese-korewireless/continuous-demo.git"
    })
    vcsRoot(gitVcs)

    // val buildTemplate = Template({
    //     id("Build")
    //     name = "build"

    //     vcs {
    //         root(gitVcs)
    //     }

    //     steps {
    //         exec {
    //             name = "Diagnostic Check"
    //             path = "./scripts/diagnostic.sh"
    //         }
    //     }

    //     triggers {
    //         vcs {
    //             id = "Trigger_1"
    //             quietPeriodMode = USE_DEFAULT
    //             triggerRules = """
    //                 +:root=${DslContext.projectId.absoluteId}_ContinuousDemo;:**
    //             """.trimIndent()
    //         }
    //     }

    //     failureConditions {
    //         executionTimeoutMin = 10
    //     }
    // })


    buildType(Default)

    params {
        param("env.BuildNumber", "1.0.%build.counter%")
    }
}

object Default : BuildType({
    name = "Default"

    vcs {
        root(gitVcs)
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
            scriptContent = "echo ${DslContext.projectId.absoluteId}"
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
