pipeline {
    agent any
    
    environment{
        GOROOT="/var/jenkins_home/go"
    }
    

    stages {
        stage('配置'){
            steps{
                 sh '''
                $GOROOT/bin/go version
                '''
                checkout scm

            }
        }
        stage('拉取') {
            steps {
                // Get some code from a GitHub repository
                git branch: 'main', credentialsId: 'Github——imailang', url: 'https://github.com/imailang/TestMachine.git'
            }
        }
        stage('构建'){
            steps{
                sh """
                    $GOROOT/bin/go build -o bin/cugtm src/main.go
                    
                """
            }
        }
        stage('代码质量评估'){
            steps{
                def scannerHome = tool 'SonarScanner';
                withSonarQubeEnv() {
                    sh "${scannerHome}/bin/sonar-scanner"
                }
            } 
        }  
    }
}
