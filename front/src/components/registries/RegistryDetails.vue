<template>
	<div id="registry-configuration">
		<b-form @submit="() => {}">
			<h3>Registry configuration</h3>
			<b-form-group 
				horizontal
				label="Name"
				label-for="registry-name-input"
				description="Meaningfull name for the registry"
			>
				<b-form-input 
					id="registry-name-input"
					type="text"
					v-model="form.name"
					required
					placeholder="Enter name">
				</b-form-input>
			</b-form-group>
			<b-form-group 
				horizontal
				label="Url"
				label-for="registry-url-input"
				description="Url under wich is the registry accessible (registry.gitlab.com for gitlab or index.docker.io for dockerhub)"
			>
				<b-form-input 
					id="registry-url-input"
					type="text"
					v-model="form.url"
					placeholder="Enter url">
				</b-form-input>
			</b-form-group>
			<b-form-group 
				horizontal
				label="User"
				label-for="registry-user-input"
				description="User that should be used for authentication (mandatory for private repositories !)"
			>
				<b-form-input 
					id="registry-user-input"
					type="text"
					v-model="form.user"
					placeholder="Enter user">
				</b-form-input>
			</b-form-group>
			<b-form-group 
				horizontal
				label="Password"
				label-for="registry-password-input"
				description="Password that should be used for authentication (mandatory for private repositories !)"
			>
				<b-form-input 
					id="registry-user-input"
					type="password"
					v-model="form.password"
					placeholder="Enter password">
				</b-form-input>
			</b-form-group>
      <b-alert 
        id="registry-success-msg" 
        class="test-msg" 
        :show="successMsg !== ''"
        variant="success"
      >
				{{successMsg}}
      </b-alert>
      <b-alert 
        id="registry-failure-msg" 
        class="test-msg" 
        :show="errorMsg !== ''" 
        variant="danger"
      >
        {{errorMsg}}
      </b-alert>
			<b-row class="my-3" align-h="end">
				<b-col cols="auto">
					<button-spin 
						@click.native="testRegistry" 
						id="registry-test-button"
						type="button"
						variant="secondary"
						:spinning="testPending"
						:disabled="savePending"
						label="Test it !"
					>
					</button-spin>
				</b-col>
				<b-col cols="auto">
					<button-spin 
						@click.native="saveRegistry" 
						id="registry-save-button"
						type="button"
						variant="primary"
						:spinning="savePending"
						:disabled="saveDisabled() || testPending"
						label="Save"
					>
					</button-spin>
				</b-col>
			</b-row>
		</b-form>
	</div>
</template>

<script>
import ButtonSpin from '@/components/controls/ButtonSpin.vue';
import { 
	testDockerRegistry,
	createDockerRegistry,
	updateDockerRegistry,
	} from '@/requests/dockerRegistries';
import { isStringFilled } from '@/misc/validation';
// details of a docker registry
// this component has two 'mode'
// docker registry creation OR docker registry update
export default {
  components: {
    ButtonSpin
  },
	props: {
		registry: {
			type: Object,
			required: false
		}
	},
  data () {
    return {
			form: {
				lastModificationTime: '',
				name: '',
				url: '',
				user: '',
				password: '',
			},
			isUpdate: false,
			passwordFake: '',
      errorMsg: '',
      successMsg: '',
			savePending: false,
			testPending: false,
		}
	},
	methods: {
    testRegistry: async function() {
      try {
        this.resetStatuses();
				this.testPending = true;
        await testDockerRegistry(this.form);
				this.successMsg = "The test is successful";
      } catch(err) {
        this.errorMsg = `An error has happend during test : ${err.message.msg}`;
      } finally {
				this.testPending = false;
			}
    },
		saveRegistry: async function() {
      this.resetStatuses();
			this.savePending = true;
			if (this.isUpdate) {
				await this.updateRegistry();
			} else {
				await this.createRegistry();
			}
		},
		createRegistry: async function() {
			try {
				await createDockerRegistry(this.form);
				this.$emit('registry-saved')
			} catch(err) {
        this.errorMsg = `An error has happend during creation : ${err.message.msg}`;
			} finally {
				this.savePending = false;
			}
		},
		updateRegistry: async function() {
			try {
				// first compute the deltas
				// between the original registry
				// and values in the form. Such mechanisme is
				// necessary because some fields are not 
				// fetched for security purposed. Then 
				// it is not possible to write the entire
				// content from the form in the db row. We must Know
				// exactly which fields are changed.
				const changedFields = [];
				const fields = Object.keys(this.form);
				fields.forEach((field) => {
					if (this.form[field] !== this.registry[field]) {
						changedFields.push(field);
					}	
				});
				await updateDockerRegistry(this.registry.id, {
					...this.form,
					changedFields,
					});
				this.$emit('registry-saved')
			} catch(err) {
				if (err.status === 409) {
					this.errorMsg = 'An error has happend during the saving : there is a conflict ! the registry has been reloaded';
					this.form = err.message;
				} else {
					this.errorMsg = `An error has happend during the saving : ${err.message.msg}`;
				}
			} finally {
				this.savePending = false;
			}
		},
		saveDisabled : function() {
			return !this.createValid() && !this.updateValid()
		},
		createValid: function() {
			return !this.isUpdate && isStringFilled(this.form.name) && isStringFilled(this.form.url);
		},
		updateValid: function() {
			return this.isUpdate && (
				this.form.name !== this.registry.name ||
				this.form.url !== this.registry.url ||
				this.form.user !== this.registry.user ||
				this.form.password !== this.registry.password
			);
		},
    resetStatuses: function() {
				this.successMsg = '';
				this.errorMsg = '';
    },
		initForm: function(givenRegistry) {
			if (givenRegistry) {
				// update of an existing registry
				this.isUpdate = true;
				this.form.lastModificationTime = givenRegistry.lastModificationTime;
				this.form.name = givenRegistry.name;
				this.form.url = givenRegistry.url;
				this.form.user = givenRegistry.user;
				this.form.password = givenRegistry.password;
			} else {
				// creation of a new registry
				this.form.name = '';
				this.form.url = '';
				this.form.user = '';
				this.form.password = '';
				this.form.lastModificationTime = '';
				this.isUpdate = false;
			}
		},
	},
	beforeMount: function() {
		this.initForm(this.registry);
	},
	watch: { 
		registry: function(newRegistry) {
      this.resetStatuses();
			this.initForm(newRegistry);
		}
	},
}
</script>
<style scoped>
#registry-configuration {
  border: 1.2px solid #4E5D6C;
	padding: 1rem 2rem 0.5rem 2rem;
}
h3 {
	margin-bottom: 1rem;
}
</style>
