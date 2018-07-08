<template>
	<div id="job-creation">
		<b-container fluid>
			<b-row class="my-3">
				<b-col sm="12">
					<h3>Job creation</h3>
				</b-col>
			</b-row>
			<b-row class="my-3">
				<b-col sm="2"><label for="newJobNameInput">Name : </label></b-col>
				<b-col sm="10">
					<b-form-input id="newJobNameInput" type="text" required placeholder="Enter name">
					</b-form-input>
				</b-col>
			</b-row>
			<b-row class="my-3">
				<b-col sm="2"><label for="authSchemeInput">Authentication : </label></b-col>
				<b-col sm="10">
					<b-form-select id="authSchemeInput" v-model="authSchemSelected" :options="authSchemOptions">
					</b-form-select>
				</b-col>
			</b-row>
			<div v-if="authSchemSelected === 'user-pwd'" class="job-authentication-conf">
				<b-row class="my-3">
					<b-col sm="12">
						<h5>https authentication configuration</h5>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="newJobUrlInput">Https url : </label></b-col>
					<b-col sm="10">
						<b-form-input id="newJobUrlInput" type="url" required placeholder="https url of your project (should start with 'https://`')">
						</b-form-input>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="repositoryLoginInput">Login : </label></b-col>
					<b-col sm="10">
						<b-form-input id="repositoryLoginInput" type="url" required placeholder="User name that should be used for authentication">
						</b-form-input>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="repositoryPasswordInput">Password : </label></b-col>
					<b-col sm="10">
						<b-form-input id="repositoryPasswordInput" type="url" required placeholder="Password that should be used for authentication">
						</b-form-input>
					</b-col>
				</b-row>
        <b-row class="my-3">
          <b-col sm="2">
            <b-button v-b-toggle.collapseinfouserpwd :size="sm" variant="link">More info</b-button>
          </b-col>
          <b-col sm="10">
            <b-collapse id="collapseinfouserpwd" class="mt-2 authentication-info">
              <span>
                Simple user - password tuple to access your repository.
                If you wish use this authentication scheme, you better have to
                create a dedicated user to restreint the access for Dahu.
                If you can't/ don't want do such thing, you should prefer
                ssh scheme.
              </span>
            </b-collapse>
          </b-col>
        </b-row>
			</div>
			<div v-else-if="authSchemSelected === 'ssh'" class="job-authentication-conf">
				<b-row class="my-3">
					<b-col sm="12">
						<h5>Ssh authentication configuration</h5>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="newJobUrlInput">Ssh url : </label></b-col>
					<b-col sm="10">
						<b-form-input id="newJobUrlInput" type="url" required placeholder="ssh url of your project (should start with 'xxx@')">
						</b-form-input>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="sshPrivateKeyInput">Ssh key : </label></b-col>
					<b-col sm="10">
						<b-form-textarea id="sshPrivateKeyInput" placeholder="Enter your ssh private key" :rows="3">
						</b-form-textarea>
					</b-col>
        </b-row>
        <b-row class="my-3">
          <b-col sm="2">
            <b-button v-b-toggle.collapseinfossh :size="sm" variant="link">More info</b-button>
          </b-col>
          <b-col sm="10">
            <b-collapse id="collapseinfossh" class="mt-2 authentication-info">
              <span>
                The private key will be used by Dahu to access your repository.
                For more information on key generation, see 
                <a 
                  target="_blank" 
                  rel="noopener noreferrer" 
                  href="https://git-scm.com/book/en/v2/Git-on-the-Server-Generating-Your-SSH-Private-Key">
                    the following tutorial.
                </a>
                Then, you must register the corresponding public key on your git server. See documentation for
                <a
                  target="_blank"
                  rel="nooper noreferrer"
                  href="https://help.github.com/articles/adding-a-new-ssh-key-to-your-github-account/"
                >
                  github
                </a>
                  or 
                <a
                  target="_blank"
                  rel="nooper noreferrer"
                  href="https://docs.gitlab.com/ee/ssh/"
                >
                  gitlab
                </a>
                  .
              </span>
            </b-collapse>
          </b-col>
        </b-row>
			</div>
			<b-row class="my-3" align-h="end">
				<b-col cols="auto">
					<b-button id="job-test-button" type="button" size="lg" variant="secondary">
						Test it !
					</b-button>
				</b-col>
				<b-col cols="auto">
					<b-button id="job-creation-button" type="submit" size="lg" variant="primary">
						Create
          </b-button>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>
<script>
// todo : test if case and render
// todo : add many explanations
// todo : implements the test 'on the fly'
export default {
  data () {
    return {
      authSchemSelected: null,
      authSchemOptions: [
        { value: null, text: 'No authentication' },
        { value: 'ssh', text: 'Ssh private key' },
        { value: 'user-pwd', text: 'User password tuple' }
      ]
    }
  }
}
</script>
<style scoped>
#job-creation {
  width: 50%;
  margin: auto;
  border: 1.2px solid #4E5D6C;
	padding: 0 1rem 0 1rem;
}
.job-authentication-conf {
	padding: 0.5rem 0.5rem 0.5rem 2rem;
  border: 0.8px solid #4E5D6C;
}
.authentication-info {
  text-align: left;
}
</style>
