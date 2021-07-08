pipeline {
    agent {
        node {
            label 'Jenkins'
        }
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '50'))
    }

    stages {
        stage("Checkout required repositories") {
            steps {
                script {
                    checkout([$class: 'GitSCM', branches: [[name: "develop"]], extensions: [], userRemoteConfigs: [[url: "git@ssh.dev.azure.com:v3/ECTLCGK/DPLT/DPLT_DemoAuthorizationPolicies"]]])
                }
            }
        }

        stage("Clean") {
            environment {
                SERVICE_NAME = "auth-demo"
            }
            steps {
                script {
                    echo 'Clean Stage...'
                    sh(script: "make clean", returnStatus: false, returnStdout: false)
                }
            }
        }

        stage("Docker Login") {
            steps {
                script {
                    echo 'Docker Login Stage...'
                    withCredentials([[
                                             $class          : "UsernamePasswordMultiBinding",
                                             credentialsId   : "${DOCKER_CREDENTIALS_ID}",
                                             usernameVariable: 'DOCKER_LOGIN',
                                             passwordVariable: 'DOCKER_PASSWORD',
                                     ]]) {
                        sh(script: 'docker login --username $DOCKER_LOGIN -p $DOCKER_PASSWORD $DOCKER_REGISTRY_LOGIN_URL', returnStatus: false, returnStdout: false)
                    }
                }
            }
        }


        stage("Build Auth demo") {
            environment {
                SERVICE_NAME = "auth-demo"
            }
            steps {
                script {
                    echo 'Build Auth demo...'
                    sh(script: "make build-image", returnStatus: false, returnStdout: false)
                }
            }
        }

        stage("Publish Auth Demo") {
            environment {
                SERVICE_NAME = "auth-demo"
            }
            steps {
                script {
                    echo 'Publish auth-demo Stage...'
                    sh(script: "make publish", returnStatus: false, returnStdout: false)
                }
            }
        }



        stage("Cleanup Docker Images") {
            environment {
                SERVICE_NAME = "auth-demo"
            }
            steps {
                script {
                    echo 'Cleanup Docker Images Stage...'
                    sh(script: "make cleanup-docker-images", returnStatus: false, returnStdout: false)
                }
            }
        }
    }

    //post {
    //    always {
    //        deleteDir()
    //    }
    //}
}
