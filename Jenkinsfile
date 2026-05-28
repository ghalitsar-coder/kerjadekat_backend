pipeline {
    agent any

    environment {
        // Karena Jenkins di Docker, dockerhub user diambil dari env
        DOCKERHUB_USER = "ghalitsar"
        IMAGE_NAME = "kerjadekat-backend"
    }

    stages {
        stage('Test CI') {
            steps {
                echo "Hello from Backend Pipeline!"
                sh "echo 'Testing backend changes...'"
            }
        }

        stage('Build Image') {
            steps {
                script {
                    // Ambil short commit hash
                    def gitCommit = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                    
                    // Build image
                    sh "docker build -t ${DOCKERHUB_USER}/${IMAGE_NAME}:${gitCommit} -t ${DOCKERHUB_USER}/${IMAGE_NAME}:latest ."
                }
            }
        }
        
        stage('Push Image') {
            steps {
                echo "Skipping push for this test run. The Build stage worked!"
            }
        }
    }
}
