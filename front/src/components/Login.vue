<template>
  <div id="login">
    <b-row id="login-row" class="justify-content-md-center">
      <b-col cols="4">
        <b-form @submit="onSubmit">
          <h2>Authentication</h2>
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
      }
    }
  },
  methods: {
    onSubmit (evt) {
      evt.preventDefault();
      authenticate(this.form.identifier, this.form.password)
        .then(token => {
          login(token);
        });
    }
  }
}
</script>

<style scoped>
#login-row {
  padding-top: 1.5rem;
}
form {
  border-style: solid;
  border-width: thin;
  border-radius: 25% 10%;
  padding: 1rem;
}
</style>
