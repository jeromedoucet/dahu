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
            :active="selectedRegistry.id === registry.id"
          >
            {{registry.name}}
          </b-list-group-item>
        </b-list-group>
			</b-col>
			<b-col cols="8" align-self="center">
				<registry-details :registry="selectedRegistry">
				</registry-details>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
// todo default registry ? => explain it
// todo new button
import RegistryDetails from '@/components/registries/RegistryDetails.vue'
export default {
  components: {
		RegistryDetails,
	},
  data () {
    return {
      selectedRegistry: {id: 1, name: 'Private Dockerhub'},
      registries: [
        {id: 1, name: 'Private Dockerhub'},
        {id: 2, name: 'Gitlab registry'},
        {id: 3, name: 'Private corporate registry'},
      ],
    }
  },
  methods: {
    // todo test me !!
    selectRegistry: function(registry) {
      this.selectedRegistry = registry;
    }
  },
}
</script>

<style scoped>
</style>
