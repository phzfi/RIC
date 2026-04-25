#!/usr/bin/env groovy
//def BUILD_ENV
//def VERSION // Note! in regular pipeline it seems that introducing VERSION as def causes VERSION to be set as null.
def CHANGELOG
def SLACK_CHANNEL
def BUILD_ENV

//Note! This project uses declarative pipeline syntax on Jenkins
pipeline {
  agent { label 'slave && vagrant' }

  environment {
    PROJECT_NAME = "ric"
    //the idea is that the team is full stack (meaning able to fix also CI issues), values high quality,
    //and is able and willing to fix errors on CI immediately, receiving notifications on project channel
    SLACK_CHANNEL = "#marketing"

    BRANCH = "${env.GIT_BRANCH.contains('/') ? env.GIT_BRANCH.split('/')[1] : env.GIT_BRANCH}"

    VERSION = "${currentBuild.number}"
  }

  options {
    //colors
    ansiColor('xterm')
    //set default pipeline timeout to 3hours if there is a jam, it will abort automatically
    timeout(time: 180, unit: 'MINUTES')
    buildDiscarder(logRotator(numToKeepStr: '50'))
  }

  triggers {
    //gitlab(triggerOnPush: true, triggerOnMergeRequest: true, branchFilterType: 'All')
    pollSCM('H/5 8-20 1-6 * *')
  }

  stages {
    stage("Prepare") {
      steps {
        //prevent Jenkins wrong branch checkout failure
        //see https://stackoverflow.com/questions/44928459/is-it-possible-to-rename-default-declarative-checkout-scm-step
        checkout scm

        //If building custom branch, the BUILD_ENV setting above returns null, revert to dev
        script {
          echo "Branch: ${BRANCH}"
          BUILD_ENV = 'dev'
          if (BRANCH == 'master') {
            BUILD_ENV = 'prod'
          } else if (BRANCH == 'develop') {
            BUILD_ENV = 'stg'
          } else {
            BUILD_ENV = 'dev'
          }
          env.BUILD_ENV = BUILD_ENV
          echo "Build Env: ${BUILD_ENV}"
        }

        //parse CHANGELOG
        script {
          def changeLogSets = currentBuild.rawBuild.changeSets
          CHANGELOG = ""
          for (int i = 0; i < changeLogSets.size(); i++) {
            def entries = changeLogSets[i].items
            for (int j = 0; j < entries.length; j++) {
              def entry = entries[j]
              CHANGELOG = CHANGELOG + "${entry.author}: ${entry.msg}\n"
            }
          }
          //prevent double builds, check if changelog is empty, skip
          if (CHANGELOG && CHANGELOG.trim().length() == 0) {
            currentBuild.result = 'SUCCESS'
            return
          }
        }
        echo "Release Notes:\n${CHANGELOG}"
      }
    }

    stage("Clean") {
      steps {
        script {
          sh "./down.sh || true"
          sh "./clean.sh"
        }
      }
    }

    stage("Provision") {
      steps {
        echo "Build -stage does all the necessary steps"
        //timeout(time: 30, unit: 'MINUTES') {
          //sh script: "./up.sh", returnStatus: true
        //}
      }
    }

    stage("Quality") {
      steps {
        echo "Running static quality analysis"
        echo "TODO: Please add a task to implement CI-6 https://wiki.phz.fi/NonFunctionalRequirements#CI"
      }
    }

    stage("Test") {
      steps {
        echo "Running tests with coverage"
        script {
          sh "./coverage.sh"
        }
        script {
          def cloverExists = sh(script: "test -f ${WORKSPACE}/reports/coverage/clover.xml && echo 'yes' || echo 'no'", returnStdout: true).trim() == "yes"
          if (cloverExists) {
            step([
                $class: 'CloverPublisher',
                cloverReportDir: "${WORKSPACE}/reports/coverage",
                cloverReportFileName: 'clover.xml',
                healthyTarget: [methodCoverage: 70, conditionalCoverage: 80, statementCoverage: 80],
                unhealthyTarget: [methodCoverage: 50, conditionalCoverage: 50, statementCoverage: 50],
                failingTarget: [methodCoverage: 0, conditionalCoverage: 0, statementCoverage: 0]
            ])
          } else {
            echo "No Clover report found - skipping CloverPublisher"
          }
        }
        publishHTML(target: [
            reportDir: "${WORKSPACE}/reports/coverage",
            reportFiles: 'coverage.html',
            reportName: 'Coverage Report'
        ])
      }
    }

    stage("Performance") {
      steps {
        echo "Running performance tests"
        echo "TODO: Please add a task to implement CI-8 https://wiki.phz.fi/NonFunctionalRequirements#CI"
      }
    }

    stage("Build") {
      steps {
        withCredentials([
          [$class: 'UsernamePasswordMultiBinding', credentialsId: 'DOCKER_HUB', usernameVariable: 'DOCKER_HUB_USERNAME', passwordVariable: 'DOCKER_HUB_PASSWORD']
        ]) {
          script {
            currentBuild.result = hudson.model.Result.SUCCESS.toString()
            if (currentBuild.result=='SUCCESS') {
              sh "docker/build.sh ${BUILD_ENV} ${VERSION}"
            } else {
              echo "FAIL: Not deploying because currentBuild.result = ${currentBuild.result}"
            }
          }
        }
      }
    }

    stage("Deploy") {
      steps {
        withCredentials([
          [$class: 'UsernamePasswordMultiBinding', credentialsId: 'DOCKER_HUB', usernameVariable: 'DOCKER_HUB_USERNAME', passwordVariable: 'DOCKER_HUB_PASSWORD']
        ]) {
          script {
            currentBuild.result = hudson.model.Result.SUCCESS.toString()
            if (currentBuild.result=='SUCCESS') {
              sh "docker/deploy.sh ${BUILD_ENV} ${VERSION}"
            } else {
              echo "FAIL: Not deploying because currentBuild.result = ${currentBuild.result}"
            }
          }
        }
      }
    }

    stage("Tag") {
      steps {
          script {
            if (env.BUILD_ENV != 'dev') {
              sshagent(credentials: ['github']) {
                withEnv(["SSH_AUTH_SOCK=${env.SSH_AUTH_SOCK}"]) { // explicitly propagate ssh auth sock

                sh('TAG_NAME="' + env.BUILD_ENV + '-' + env.VERSION + '" && '
                + 'git tag -d $TAG_NAME || true && ' // delete 'existing' tag from local git repository. (if previous push failed)
                + 'git tag -a $TAG_NAME -m Jenkins && ' // create new tag
                + 'sleep 2 && ' // wait ssh-agent to initialize
                + 'git push -v git@github.com:phzfi/ric refs/tags/$TAG_NAME --no-verify' // push the new tag via SSH
                )
                }
              }
            } else {
              echo "Skipping Git Tag and Push for git development branches..."
            }
          }
      }
    }
  }

  post {
    always {
      script {
        currentBuild.result = hudson.model.Result.SUCCESS.toString()
        if (currentBuild.result!='SUCCESS') {
          echo "FAIL: At Post beginning currentBuild.result = ${currentBuild.result}"
        }
        //See https://docs.cloudbees.com/docs/cloudbees-ci-kb/latest/troubleshooting-guides/how-to-troubleshoot-hudson-filepath-is-missing-in-pipeline-run
        if (getContext(hudson.FilePath)) {
          sh "./clean.sh || true"
        }
        // Workaround to the clean issue, can't delete folder as folder is owned by docker user 'root'.
        sh "sudo chown -R jenkins:jenkins $WORKSPACE || true"

        currentBuild.result = hudson.model.Result.SUCCESS.toString()
        if (currentBuild.result!='SUCCESS') {
          echo "FAIL: Post action end currentBuild.result = ${currentBuild.result}"
        }
      }
    }

    success {
      slackSend channel: "${env.SLACK_CHANNEL}", color: "good", message: "Deployed ${env.JOB_NAME}#${env.BUILD_NUMBER} successfully to ${env.BUILD_ENV}, please Smoke Tests (see README.md #4.1). Add Reaction thumbsup or thumbsdown to indicate Smoke Test cases pass or not.\n${CHANGELOG}"

      emailext (
        subject: "Deployed ${env.JOB_NAME} to ${env.BUILD_ENV} [${env.BUILD_NUMBER}]",
        body: """<p>New build completed and deployed successfully: '${env.JOB_NAME} [${env.BUILD_NUMBER}]':</p>
          <p>Please Smoke Tests (see README.md #4.1). Add Reaction thumbsup or thumbsdown on Slack to indicate Smoke Test cases pass or not.</p>
          <p>${CHANGELOG}""",
        recipientProviders: [[$class: 'DevelopersRecipientProvider']]
      )
      script {
        sh './down.sh > /dev/null 2>&1 || true'
      }
    }

    unstable {
      slackSend channel: "${env.SLACK_CHANNEL}", color: "warning", message: "Unstable build ${env.JOB_NAME}#${env.BUILD_NUMBER} to ${env.BUILD_ENV}, please fix: ${env.BUILD_URL}console#footer\n${CHANGELOG}"

      emailext (
        subject: "Unstable build ${env.JOB_NAME} to ${env.BUILD_ENV} [${env.BUILD_NUMBER}]",
        body: """<p>Unstable build: '${env.JOB_NAME} [${env.BUILD_NUMBER}]':</p>
          <p>Check console output at &QUOT;<a href='${env.BUILD_URL}console#footer'>${env.JOB_NAME} [${env.BUILD_NUMBER}]</a>&QUOT;</p>
          <p>${CHANGELOG}""",
        recipientProviders: [[$class: 'DevelopersRecipientProvider']]
      )
    }

    failure {
      slackSend channel: "${env.SLACK_CHANNEL}", color: "danger", message: "FAIL ${env.JOB_NAME}#${env.BUILD_NUMBER} to ${env.BUILD_ENV}, please fix: ${env.BUILD_URL}console#footer\n${CHANGELOG}"

      emailext (
        subject: "Failed to build ${env.JOB_NAME} to ${env.BUILD_ENV} [${env.BUILD_NUMBER}]",
        body: """<p>Build failed: '${env.JOB_NAME} [${env.BUILD_NUMBER}]':</p>
          <p>Check console output at &QUOT;<a href='${env.BUILD_URL}console#footer'>${env.JOB_NAME} [${env.BUILD_NUMBER}]</a>&QUOT;</p>
          <p>${CHANGELOG}""",
        recipientProviders: [[$class: 'DevelopersRecipientProvider']]
      )
    }
  }
}
