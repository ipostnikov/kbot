pipeline {
    agent any

    environment {
        REPO = 'https://github.com/ipostnikov/kbot'
        BRANCH = 'main'
    }

    stages {
        stage("clone") {
            steps {
                echo 'CLONE REPOSITORY'
                git branch: "${BRANCH}", url: "${REPO}"
            }
        }

        stage("test") {
            steps {
                echo 'TEST EXECUTION STARTED'
                sh 'make test'
            }
        }

        stage("build") {
            steps {
                echo 'BUILD EXECUTION STARTED'
                sh 'make build'
            }
        }

        stage("image") {
            steps {
                script {
                    echo 'BUILD EXECUTION STARTED'
                    sh 'make image'
                }
            }
        }

        stage("push") {
            steps {
                script {
                    // Changed Docker registry to ghcr.io
                    docker.withRegistry('ghcr.io') { 
                        sh 'make push'
                    }
                }
            }
        }
    }
}
