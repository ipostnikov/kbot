pipeline {
    agent any

    parameters {
        string(name: 'REPO', defaultValue: 'https://github.com/ipostnikov/kbot', description: 'GitHub repository to clone')
        string(name: 'BRANCH', defaultValue: 'main', description: 'Repository branch for building')
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

    stages {
        stage("clone") {
            steps {
                echo "CLONING REPOSITORY: ${params.REPO} BRANCH: ${params.BRANCH}"
                git branch: "${params.BRANCH}", url: "${params.REPO}"
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