pipeline {
  agent {label 'docker'}

  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }
    stage('Build') {
      steps {
        sh '''docker-compose -f docker-compose.build.yml up -d --build --force-recreate'''
        sh '''docker exec -i  ric_build /bin/bash /ric/scripts/build.sh'''

      }
    }
    stage('Test') {
      steps {
        sh '''docker exec -i  ric_build  /bin/bash /ric/scripts/run_tests.sh'''
        junit "server/junit.xml"
        cobertura autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: "server/coverage.xml", conditionalCoverageTargets: '70, 0, 0', failUnhealthy: false, failUnstable: false, lineCoverageTargets: '80, 0, 0', maxNumberOfBuilds: 0, methodCoverageTargets: '80, 0, 0', onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false
      }
    }
    stage('Build .deb') {
      steps {
        echo "build .deb here"
        sh '''docker exec -i  ric_build /bin/bash /ric/scripts/build_deb.sh'''
        archiveArtifacts artifacts: 'build/phz-ric.deb', onlyIfSuccessful: true
      }
    }
  }
  post {
    always {
      sh '''
        docker-compose -f docker-compose.build.yml down --rmi all
      '''
    }
    failure {
      echo 'TODO: add failure'
    }
  }
}
