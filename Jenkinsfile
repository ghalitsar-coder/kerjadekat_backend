pipeline {
        agent any
        environment {
            DOCKERHUB_USER = "ghalitsar"
            IMAGE_NAME = "kerjadekat-backend"
        }
        stages {
            stage('Build Docker Image') {
                steps {
                    script {
                        // Ambil 7 karakter pertama dari commit hash git untuk tag image
                        env.COMMIT_HASH = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                        sh "docker build -t ${DOCKERHUB_USER}/${IMAGE_NAME}:${COMMIT_HASH} -t ${DOCKERHUB_USER}/${IMAGE_NAME}:latest -f Dockerfile ."
                    }
                }
            }
            stage('Push to Docker Hub') {
                steps {
                    withCredentials([usernamePassword(credentialsId: 'dockerhub-creds', passwordVariable: 'DOCKER_PASS', usernameVariable: 'DOCKER_USER')]) {
                        sh 'echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin'
                        sh "docker push ${DOCKERHUB_USER}/${IMAGE_NAME}:${COMMIT_HASH}"
                        sh "docker push ${DOCKERHUB_USER}/${IMAGE_NAME}:latest"
                    }
                }
            }
            stage('Update GitOps Repo (ArgoCD Trigger)') {
                steps {
                    // Pastikan kamu punya credentials 'github-credentials' di Jenkins
                    withCredentials([usernamePassword(credentialsId: 'github-credentials', passwordVariable: 'GIT_PASS', usernameVariable: 'GIT_USER')]) {
                        sh """
                        git clone https://${GIT_USER}:${GIT_PASS}@github.com/ghalitsar-coder/kerjadekat-gitops.git gitops_repo
                        cd gitops_repo
                        
                        git config user.email "jenkins@kerjadekat.com"
                        git config user.name "Jenkins CI"
                        
                        # Ubah tag image di file deployment backend
                        sed -i "s|image: ghalitsar/kerjadekat-backend:.*|image: ghalitsar/kerjadekat-backend:${COMMIT_HASH}|" gitops/base/backend/deployment.yaml
                        
                        git add gitops/base/backend/deployment.yaml
                        git commit -m "ci(backend): update image tag to ${COMMIT_HASH}" || echo "No changes"
                        git push origin master
                        """
                    }
                }
            }
        }
    }
