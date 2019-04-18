pipeline {
    environment {
        TEST_NS_SINGLE = 'https://10.3.199.252:8443'
        TEST_NS_HA_1 = 'https://10.3.199.252:8443'
        TEST_NS_HA_2 = 'https://10.3.199.253:8443'
    }
    options {
        disableConcurrentBuilds()
    }
    agent {
        node {
            label 'solutions-126'
        }
    }
    stages {
        stage('Tests [unit]') {
            steps {
                sh 'make test-unit-container'
            }
        }
        stage('Tests [e2e-ns-single]') {
            steps {
                sh 'TEST_NS_SINGLE=${TEST_NS_SINGLE} make test-e2e-ns-single-container'
            }
        }
        stage('Tests [e2e-ns-cluster]') {
            steps {
                sh 'TEST_NS_HA_1=${TEST_NS_HA_1} TEST_NS_HA_2=${TEST_NS_HA_2} make test-e2e-ns-cluster-container'
            }
        }
    }
}
