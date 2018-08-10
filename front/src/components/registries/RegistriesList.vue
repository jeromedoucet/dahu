<template>
  <b-container fluid>
    <b-row>
      <b-col cols="4" align-self="center">
        <b-list-group>
          <b-list-group-item 
            v-for="registry in registries" 
            button 
            v-on:click="() => selectRegistry(registry)"
            :key="registry.id"
            :active="!!selectedRegistry && selectedRegistry.id === registry.id"
          >
            {{registry.name}}
          </b-list-group-item>
          <b-list-group-item 
            button 
            v-on:click="() => selectRegistry(null)"
            key="-1"
            :active="!selectedRegistry"
          >
            Create a new registry
          </b-list-group-item>
        </b-list-group>
			</b-col>
			<b-col cols="8" align-self="center">
				<registry-details :registry="selectedRegistry" v-on:registry-saved="fetchRegistries()">
				</registry-details>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
// todo test
import RegistryDetails from '@/components/registries/RegistryDetails.vue'
import { fetchDockerRegistries } from '@/requests/dockerRegistries';
export default {
  components: {
		RegistryDetails,
	},
  data () {
    return {
      selectedRegistry: null,
      registries: [],
    }
  },
  methods: {
    selectRegistry: function(registry) {
      this.selectedRegistry = registry;
    },
    fetchRegistries: async function() {
      try {
        this.registries = await fetchDockerRegistries();
				if (this.registries.length > 0) {
					this.selectedRegistry = this.registries[0];
				}
      } catch (err) {
        // todo error msg
      }
    }
  },
	beforeMount: function() {
		this.fetchRegistries();
	},
}
</script>

<style scoped>
</style>
