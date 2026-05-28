pipeline {
    agent any
    environment {
        DOCKERHUB_USER = "ghalitsar"
        IMAGE_NAME = "kerjadekat-backend"
    }
    stages {
        stage('Checkout') {
            steps {
                sh '''
                # Clean workspace
                rm -rf *
                git clone https://github.com/ghalitsar-coder/kerjadekat_backend.git .
                git rev-parse --short HEAD > git_commit.txt
                '''
            }
        }
        stage('Build Docker Image') {
            steps {
                script {
                    def commitHash = readFile('git_commit.txt').trim()
                    sh "docker build -t ${DOCKERHUB_USER}/${IMAGE_NAME}:${commitHash} -t ${DOCKERHUB_USER}/${IMAGE_NAME}:latest ."
                }
            }
        }
        stage('Push to Docker Hub') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub-creds', passwordVariable: 'DOCKER_PASS', usernameVariable: 'DOCKER_USER')]) {
                    sh "echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin"
                    script {
                        def commitHash = readFile('git_commit.txt').trim()
                        sh "docker push ${DOCKERHUB_USER}/${IMAGE_NAME}:${commitHash}"
                        sh "docker push ${DOCKERHUB_USER}/${IMAGE_NAME}:latest"
                    }
                }
            }
        }
    }
}