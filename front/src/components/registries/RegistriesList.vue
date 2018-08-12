<template>
  <b-container fluid>
    <b-row>
      <b-col cols="4" align-self="center">
        <b-list-group>
          <b-list-group-item 
            v-for="registry in registries" 
            v-on:click="() => selectRegistry(registry)"
            :key="registry.id"
            :active="!!selectedRegistry && selectedRegistry.id === registry.id"
          >
            {{registry.name}}
            <delete-item 
              :on-delete="() => deleteRegistry(registry.id)" 
              v-on:item-deleted="fetchRegistries()" 
              :item-label="registry.name" 
            />
          </b-list-group-item>
        </b-list-group>
        <new-item class="new-registry" @click.native="() => selectRegistry(null)" />
			</b-col>
			<b-col cols="8" align-self="center">
				<registry-details :registry="selectedRegistry" v-on:registry-saved="fetchRegistries()">
				</registry-details>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
// todo tests
import RegistryDetails from '@/components/registries/RegistryDetails.vue'
import { fetchDockerRegistries, deleteDockerRegistry } from '@/requests/dockerRegistries';
import DeleteItem from '@/components/controls/DeleteItem';
import NewItem from '@/components/controls/NewItem.vue'
export default {
  components: {
		RegistryDetails,
    NewItem,
    DeleteItem,
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
    },
    deleteRegistry: async function(id) {
      await deleteDockerRegistry(id)
    },
  },
	beforeMount: function() {
		this.fetchRegistries();
	},
}
</script>

<style scoped>
.new-registry {
  margin-top: 5px;
}
</style>
