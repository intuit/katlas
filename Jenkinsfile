podTemplate(label: 'jnlp', containers: [
        containerTemplate(name: 'ci', image: 'docker.artifactory.a.intuit.com/dev/containers/jnlp-slave-for-cicd/service/jnlp-slave-ci-docker-maven:a87edc2', ttyEnabled: true, command: 'cat', args: ''),
	containerTemplate(name: 'golang', image: 'golang:1.11.2', ttyEnabled: true, command: 'cat'),
    ],
  volumes: [hostPathVolume(hostPath: '/var/run/dind/docker.sock', mountPath: '/var/run/docker.sock')]) {

    def repoArtifactory = "docker.artifactory.a.intuit.com"
    def imageArtifactory = "dev/devx/k8scmdb/service"
    def appNameUi = "katlas-browser"
    def appNameApi = "katlas-service"
    def cdProjectUrl = "github.intuit.com/dev-devx/cutlass-api-deploy.git"

    node('jnlp') {
        container('ci') {
        artifactory_cred_id = 'k8scmdb-token' //token can push to artifactory
        github_cred_id = 'github_div_cred' //token can push to -deploy project, IBP 2 token does not have permission

        // Initialize
        stage('Initialization') {
            checkout scm
            commitId = sh(returnStdout: true, script: "git rev-parse HEAD").trim()
            shortCommit = sh(returnStdout: true, script: "git rev-parse --short HEAD").trim()
            // output all the environment variables, testing purpose
            echo sh(script: 'env', returnStdout: true)
	    // pick the latest jira commit
    	    commitMsg = sh(returnStdout: true, script: "git log -1 --pretty=%B").trim()
    	    echo "The Commit message: ${commitMsg}"
        }

        // Test Stage
        stage('Test') {
            container('golang') {
		sh "hostname;ls -al"
		sh '''cd service 
		    find db apis resources -name "*_test.go" ! -name "*mock_test.go" -exec rm {} \\;
		    export GOPATH=`pwd`/.gopath~ && export GOBIN=`pwd`/.gopath~/bin && make -d test-coverage
		'''
	    }
        }

        stage('docker build') {
            echo "Building image ${appNameApi}:${shortCommit}"

            sh "docker images"
            sh "cd service && docker build --no-cache -f Dockerfile -t ${appNameApi}:${shortCommit} ."

            echo "Building image ${appNameUi}:${shortCommit}"
            sh "cd app; docker build --no-cache -f Dockerfile -t ${appNameUi}:${shortCommit} ."
        }

        stage('Push docker image to Artifactory') {
            docker.withRegistry("https://${repoArtifactory}", artifactory_cred_id) {
                // Pushing multiple tags is cheap, as all the layers are reused.
                sh "docker tag ${appNameApi}:${shortCommit} ${repoArtifactory}/${imageArtifactory}/${appNameApi}:${shortCommit}"
                sh "docker tag ${appNameUi}:${shortCommit} ${repoArtifactory}/${imageArtifactory}/${appNameUi}:${shortCommit}"

                // mnau - testing with jfrog
                sh "docker version"
                sh "docker info"
                sh "docker images"
                sh "docker push ${repoArtifactory}/${imageArtifactory}/${appNameApi}:${shortCommit}"
                sh "docker push ${repoArtifactory}/${imageArtifactory}/${appNameUi}:${shortCommit}"
            } // docker.withRegistry
        } // stage

        // Trigger CD
        stage('Trigger Deployment') {
            cdProjectFolder = 'prj-deploy'
            withCredentials([string(credentialsId: github_cred_id, variable: 'GITHUB_TOKEN')]) {
            echo "****inside CD Git repoi****"
            sh "git clone https://${GITHUB_TOKEN}@${cdProjectUrl} ${cdProjectFolder}"
            }
            dir(cdProjectFolder){
                // update the phase and image tag
                params = readYaml file: "params.yaml"
                //for API
                params['image']['repository'] = "${repoArtifactory}/${imageArtifactory}/${appNameApi}"
                echo   "API repo: ${params['image']['repository']}"
                params['image']['tag'] = shortCommit
                echo   "API tag: ${params['image']['tag']}"
                //and UI (notice 'imageui' in place of 'image' attr)
                params['imageui']['repository'] = "${repoArtifactory}/${imageArtifactory}/${appNameUi}"
                echo   "UI repo: ${params['imageui']['repository']}"
                params['imageui']['tag'] = shortCommit
                echo   "UI tag: ${params['imageui']['tag']}"
                if (env.BRANCH_NAME == 'master') {
                    params['phase'] = 'release'
                	params['branch'] = env.BRANCH_NAME
                	params['commitId'] = commitId
                	params['ask'] = "false"
                	sh 'rm params.yaml'
                	sh 'chmod -R a+rw .'
                	writeYaml file: 'params.yaml', data: params
                	sh "ls -ltra ${pwd()}"
                	sh "cat params.yaml"
                	// config the user/email
                	sh 'git config user.email "ci@intuit.com"'
                	sh 'git config user.name "ci"'
    
                	// commit to repo
                	sh 'git add params.yaml'
    
                	sh "git commit -m \"${commitMsg} ,pipeline from $shortCommit, build: ${env.BUILD_URL}\""
                	sh "git push origin master"
                }
            }
        }  // stage
      } // ci
    } // node
  } //pipeline
