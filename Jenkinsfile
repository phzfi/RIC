def BUILD_ENV
//def VERSION // Note! in regular pipeline it seems that introducing VERSION as def causes VERSION to be set as null.
def CHANGELOG
def SLACK_CHANNEL

//Declarative pipeline example
pipeline {
  agent { label 'slave && vagrant' }

  environment {
    PROJECT_NAME = "ric"
    //the idea is that the team is full stack (meaning able to fix also CI issues), values high quality,
    //and is able and willing to fix errors on CI immediately, receiving notifications on project channel
    SLACK_CHANNEL = "#marketing"

    //Multibranch pipeline
    BUILD_ENV = [master: 'prod', develop: 'stg'].get(env.BRANCH_NAME, 'dev')
    VERSION = "${currentBuild.id}"
  }

  options {
    //colors
    ansiColor('xterm')
    //set default pipeline timeout to 3hours if there is a jam, it will abort automatically
    timeout(time: 180, unit: 'MINUTES')
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
          if (BUILD_ENV == null) {
            BUILD_ENV = 'dev'
          }
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
        echo "TODO: Please add a task to implement CI-7 https://wiki.phz.fi/NonFunctionalRequirements#CI"
        //sh "docker-compose run app yarn test-ci"
        //junit 'results/*.xml'
        step([
            $class: 'CloverPublisher',
            cloverReportDir: 'reports/coverage',
            cloverReportFileName: 'clover.xml',
            healthyTarget: [methodCoverage: 70, conditionalCoverage: 80, statementCoverage: 80],
            unhealthyTarget: [methodCoverage: 50, conditionalCoverage: 50, statementCoverage: 50],
            failingTarget: [methodCoverage: 0, conditionalCoverage: 0, statementCoverage: 0]
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
  }

  post {
    always {
      script {
        //See https://docs.cloudbees.com/docs/cloudbees-ci-kb/latest/troubleshooting-guides/how-to-troubleshoot-hudson-filepath-is-missing-in-pipeline-run
        if (getContext(hudson.FilePath)) {
          sh "./clean.sh || true"
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
        sh "./down.sh"
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
