pipeline {
    agent any

    environment {
        DOCKERHUB_USER = "YOUR_DOCKERHUB_USER"
        IMAGE_NAME = "kerjadekat-backend"
        GITOPS_REPO = "git@github.com:YOUR_ORG/kerjadekat.git"
    }

    stages {
        stage('Lint & Test') {
            steps {
                dir('backend') {
                    sh 'go vet ./...'
                    sh 'go test ./...'
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                dir('backend') {
                    sh "docker build -t ${DOCKERHUB_USER}/${IMAGE_NAME}:${GIT_COMMIT} -t ${DOCKERHUB_USER}/${IMAGE_NAME}:latest ."
                }
            }
        }

        stage('Push to Docker Hub') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub-creds', passwordVariable: 'DOCKER_PASS', usernameVariable: 'DOCKER_USER')]) {
                    sh "echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin"
                    sh "docker push ${DOCKERHUB_USER}/${IMAGE_NAME}:${GIT_COMMIT}"
                    sh "docker push ${DOCKERHUB_USER}/${IMAGE_NAME}:latest"
                }
            }
        }

        stage('Update GitOps Manifest') {
            steps {
                sshagent(['github-ssh-key']) {
                    sh '''
                        # Configure Git
                        git config --global user.email "ci@kerjadekat.id"
                        git config --global user.name "Jenkins CI"

                        # Clone GitOps repo
                        git clone ${GITOPS_REPO} gitops-repo
                        cd gitops-repo

                        # Update deployment YAML using yq
                        yq e ".spec.template.spec.containers[0].image = \\"${DOCKERHUB_USER}/${IMAGE_NAME}:${GIT_COMMIT}\\"" -i gitops/base/backend/deployment.yaml

                        # Commit and push
                        git add gitops/base/backend/deployment.yaml
                        git diff-index --quiet HEAD || git commit -m "ci: update backend image to ${GIT_COMMIT} [skip ci]"
                        git push origin main
                    '''
                }
            }
        }
    }
}
