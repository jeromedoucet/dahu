<template>
	<div id="job-creation">
		<b-container fluid>
			<b-row class="my-3">
				<b-col sm="12">
					<h3>Job creation</h3>
				</b-col>
			</b-row>
			<b-row class="my-3">
				<b-col sm="2"><label for="newJobNameInput">Name* : </label></b-col>
				<b-col sm="10">
					<b-form-input id="newJobNameInput" type="text" required placeholder="Enter name">
					</b-form-input>
				</b-col>
			</b-row>
			<b-row class="my-3">
				<b-col sm="2"><label for="authSchemeInput">Authentication* : </label></b-col>
				<b-col sm="10">
					<b-form-select id="authSchemeInput" v-model="authSchemSelected" :options="authSchemOptions">
					</b-form-select>
				</b-col>
			</b-row>
			<div id="http-auth-conf" v-if="authSchemSelected === 'http'" class="job-authentication-conf">
				<b-row class="my-3">
					<b-col sm="12">
						<h5>http authentication configuration</h5>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="newJobUrlInput">Http url* : </label></b-col>
					<b-col sm="10">
						<b-form-input 
              id="newJobUrlInput" 
              type="url" 
              required 
              v-model="httpForm.url"
              placeholder="http url of your project (should start with 'http://' or 'https://')"
            >
						</b-form-input>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="repositoryLoginInput">Login : </label></b-col>
					<b-col sm="10">
						<b-form-input 
              id="repositoryLoginInput" 
              type="text" 
              v-model="httpForm.user"
              placeholder="User name that should be used for authentication"
            >
						</b-form-input>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="repositoryPasswordInput">Password : </label></b-col>
					<b-col sm="10">
						<b-form-input 
              id="repositoryPasswordInput" 
              type="password" 
              v-model="httpForm.password"
              placeholder="Password that should be used for authentication"
            >
						</b-form-input>
					</b-col>
				</b-row>
        <b-row class="my-3">
          <b-col sm="2">
            <b-button v-b-toggle.collapseinfouserpwd size="sm" variant="link">More info</b-button>
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
			<div id="ssh-auth-conf" v-else-if="authSchemSelected === 'ssh'" class="job-authentication-conf">
				<b-row class="my-3">
					<b-col sm="12">
						<h5>Ssh authentication configuration</h5>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="newJobUrlInput">Ssh url* : </label></b-col>
					<b-col sm="10">
						<b-form-input 
              id="newJobUrlInput" 
              type="text" 
              required 
              v-model="sshForm.url"
              placeholder="ssh url of your project (should start with 'xxx@')"
            >
						</b-form-input>
					</b-col>
				</b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="sshPrivateKeyInput">Ssh key* : </label></b-col>
					<b-col sm="10">
						<b-form-textarea 
              id="sshPrivateKeyInput" 
              placeholder="Enter your ssh private key" 
              required 
              v-model="sshForm.key"
              :rows="3"
            >
						</b-form-textarea>
					</b-col>
        </b-row>
				<b-row class="my-3">
					<b-col sm="2"><label for="sshPrivateKeyPasswordInput">Ssh key password : </label></b-col>
					<b-col sm="10">
						<b-form-textarea 
              id="sshPrivateKeyPasswordInput" 
              placeholder="Enter your ssh private key password" 
              v-model="sshForm.keyPassword"
              :rows="3"
            >
						</b-form-textarea>
					</b-col>
        </b-row>
        <b-row class="my-3">
          <b-col sm="2">
            <b-button v-b-toggle.collapseinfossh size="sm" variant="link">More info</b-button>
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
      <b-alert 
        id="repo-test-failure-msg" 
        class="test-msg" 
        :show="errorMsg !== ''" 
        variant="danger"
      >
        An error has happend during test : {{errorMsg}}
      </b-alert>
      <b-alert 
        id="repo-test-success-msg" 
        class="test-msg" 
        v-if="isSuccess"
        show 
        variant="success"
      >
        Repository configuration is correct !
      </b-alert>
      <b-row class="my-3" align-h="end">
				<b-col cols="auto">
					<button-spin 
            @click.native="onTest()" 
            id="job-test-button" 
            type="button" 
            variant="secondary"
            :spinning="testPending"
            label="Test it!"
          >
					</button-spin>
				</b-col>
				<b-col cols="auto">
					<button-spin 
            id="job-creation-button" 
            type="submit" 
            variant="primary"
            :disabled="testPending"
            label="Create"
          >
          </button-spin>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>
<script>
import { testRepo } from '@/requests/scm';
import ButtonSpin from '@/components/controls/ButtonSpin.vue';

export default {
  components: {
    ButtonSpin
  },
  data () {
    return {
      authSchemSelected: null,
      authSchemOptions: [
        { value: null, text: 'No authentication' },
        { value: 'ssh', text: 'Ssh private key' },
        { value: 'http', text: 'User password tuple' },
      ],
      sshForm: {
        url: '',
        key: '',
        keyPassword: ''
      },
      httpForm: {
        url: '',
        user: '',
        password: ''
      },
      errorMsg: '',
      isSuccess: false,
      testPending: false
    }
  },
  methods: {
    onTest: async function () {
      this.errorMsg = '';
      this.isSuccess = false;
      this.testPending = true;
      const scmConf = this.authSchemSelected === 'ssh' ? { sshAuth: this.sshForm } : { httpAuth: this.httpForm };
      try {
        await testRepo(scmConf);
        this.isSuccess = true;
      } catch (err) {
        this.errorMsg = err.message;
      } finally {
        this.testPending = false;
      }
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
.test-msg {
  margin-top: 1rem;
}
</style>
