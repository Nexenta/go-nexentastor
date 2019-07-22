pipeline {
    parameters {
        string(
            name: 'TEST_NS_SINGLE',
            defaultValue: 'https://10.3.199.247:8443',
            description: 'Single NS API address',
            trim: true
        )
        string(
            name: 'TEST_NS_HA_1',
            defaultValue: '',
            description: 'NS HA cluster 1 node API address',
            trim: true
        )
        string(
            name: 'TEST_NS_HA_2',
            defaultValue: '',
            description: 'NS HA cluster 2 node API address',
            trim: true
        )
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
                sh "make test-unit-container"
            }
        }
        stage('Tests [e2e-ns-single]') {
            steps {
                sh """
                    TEST_NS_SINGLE=${params.TEST_NS_SINGLE} \
                    make test-e2e-ns-single-container;
                """
            }
        }
        stage('Tests [e2e-ns-cluster]') {
            when {
                allOf {
                    expression { params.TEST_NS_HA_1 != '' }
                    expression { params.TEST_NS_HA_1 != null }
                    expression { params.TEST_NS_HA_2 != '' }
                    expression { params.TEST_NS_HA_2 != null }
                }
            }
            steps {
                sh """
                    TEST_NS_HA_1=${params.TEST_NS_HA_1} \
                    TEST_NS_HA_2=${params.TEST_NS_HA_2} \
                    make test-e2e-ns-cluster-container;
                """
            }
        }
    }
}
