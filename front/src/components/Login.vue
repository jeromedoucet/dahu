<template>
  <div id="login">
    <b-row id="login-row" class="justify-content-md-center">
      <b-col cols="4">
        <b-form @submit="onSubmit">
          <h2 class="text">Authentication</h2>
          <b-form-group>
            <b-form-input id="login-identifier-field"
                          type="text"
                          v-model="form.identifier"
                          size="lg"
                          required
                          placeholder="Enter identifier">
            </b-form-input>
          </b-form-group>
          <b-form-group>
            <b-form-input id="login-password-field"
                          type="password"
                          v-model="form.password"
                          size="lg"
                          required
                          placeholder="Enter password">
            </b-form-input>
          </b-form-group>
          <b-form-group>
            <b-alert variant="danger"
                     dismissible
                     :show="!!errorMessage"
                     @dismissed="errorMessage=null">
              {{errorMessage}}
            </b-alert>
          </b-form-group>
          <b-button id="login-submit-button" 
                    type="submit" 
                    size="lg"
                    variant="primary">
            Login
          </b-button>
        </b-form>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import { authenticate } from '@/requests/authentication'
import { login } from '@/services/user';

export default {
  name: 'Login',
  props: {
  },
  data () {
    return {
      form: {
        identifier: '',
        password: ''
      },
      errorMessage: null
    }
  },
  methods: {
    onSubmit: async function (evt) {
      evt.preventDefault();
      try {
        const token = await authenticate(this.form.identifier, this.form.password);
        login(token);
        this.$router.go('/');
      } catch (error) {
        if (error.status === 401 || error.status === 404) {
          this.errorMessage = 'Authentication error. Please check your credentials and try again.';
        } else {
          this.errorMessage = 'Unknown error. Please contact your administrator.';
        }
      }
    }
  }
}
</script>

<style scoped>
#login-row {
  padding-top: 1.5rem;
}
</style>
