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
        echo "Build steps here"
      }
    }
    stage('Test') {
      steps {
        echo 'TODO: add tests'
      }
    }
    stage('Build .deb') {
      steps {
        echo "build .deb here"
      }
    }
  }
  post {
    always {
      sh '''
       echo "Clean up here!"
      '''
    }
    failure {
      echo 'TODO: add failure'
    }
  }
}
