pipeline {
    agent any

    parameters {
        string(name: 'REPO', defaultValue: 'https://github.com/ipostnikov/kbot', description: 'GitHub repo to clone')
        string(name: 'BRANCH', defaultValue: 'main', description: 'Repo branch for building')
        choice(
            name: 'OS',
            choices: ['linux', 'darwin', 'windows'],
            description: 'Target operating system'
        )
        choice(
            name: 'ARCH',
            choices: ['amd64', 'arm64'],
            description: 'Target architecture'
        )
        booleanParam(
            name: 'SKIP_TESTS',
            defaultValue: false,
            description: 'Skip running tests'
        )
        booleanParam(
            name: 'SKIP_LINT',
            defaultValue: false,
            description: 'Skip running linter'
        )
    }

    environment {

        CURRENT_REPO = params.REPO
        CURRENT_BRANCH = params.BRANCH
    }

    stages {
        stage("clone") {
            steps {
                echo "CLONING REPOSITORY: ${CURRENT_REPO} BRANCH: ${CURRENT_BRANCH}"
                git branch: "${CURRENT_BRANCH}", url: "${CURRENT_REPO}"
            }
        }

        stage("test") {
            // This stage will be skipped if SKIP_TESTS parameter is true
            when { expression { return !params.SKIP_TESTS } }
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
                    echo 'IMAGE BUILD EXECUTION STARTED'
                    sh 'make image'
                }
            }
        }

        stage("push") {
            steps {
                script {
                    docker.withRegistry('ghcr.io') {
                        sh 'make push'
                    }
                }
            }
        }
    }
}
