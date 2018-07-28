<template>
  <b-card 
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
