pipeline {
  agent any

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
      }
    }
    stage('Build .deb') {
      steps {
        echo "build .deb here"
        sh '''docker exec -i  ric_build /bin/bash /ric/scripts/build_deb.sh'''
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
