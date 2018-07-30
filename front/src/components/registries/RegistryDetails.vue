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
				description="Url under wich is the registry accessible (registry.gitlab.com for gitlab), default is dockerhub"
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
			<b-row class="my-3" align-h="end">
				<b-col cols="auto">
					<button-spin 
						@click.native="() => {}" 
						id="registry-test-button"
						type="button"
						variant="secondary"
						:spinning="false"
						:disabled="true"
						label="Test it!"
					>
					</button-spin>
				</b-col>
				<b-col cols="auto">
					<button-spin 
						@click.native="() => {}" 
						id="registry-save-button"
						type="button"
						variant="primary"
						:spinning="false"
						:disabled="true"
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
// todo mark if it's a new or not (id ?)
// todo copy props to data
// todo when update, save disable is no change
// todo test button
// todo think of feeback messages
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
				name: '',
				url: '',
				user: '',
				password: '',
			}
		}
	},
	methods: {
		initForm: function(givenRegistry) {
			if (givenRegistry) {
				// update of an existing registry
				this.form.name = givenRegistry.name;
			} else {
				// creation of a new registry
				this.form.name = '';
			}
		}
	},
	// todo test me !!!
	beforeMount: function() {
		this.initForm(this.registry);
	},
	// todo test me !!!
	watch: { 
		registry: function(newRegistry) {
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
