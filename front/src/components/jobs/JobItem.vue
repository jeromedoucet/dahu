<template>
  <b-card 
		v-on:click="goToPipeline()"
    :title="job.name"
    bg-variant="secondary"
    class="mb-2 card job-card"
  >
    <p class="card-text">
    </p>
  </b-card>
</template>
<script>
export default {
	props: {
		job: {
			type: Object,
			required: true
		}
	},
	methods: {
		goToPipeline: function() {
			this.$router.push(`/jobs/${this.job.id}/pipeline`);
		}
	}
}
</script>
<style scoped>
.card {
  color: #cccccc;
  margin: 0 1rem 0 1rem;
  width: 200px;
  height: 150px;
}
.card:hover {
  color: #ffffff;
  cursor: pointer;
}
</style>
