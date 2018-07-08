<template>
  <div id="home">
    <div class="jobs">
      <b-card title="Job Name"
              bg-variant="success"
              tag="article"
              class="mb-2 card">
        <p class="card-text">
          This is a job on success state.
        </p>
      </b-card>
      <b-card title="Job Name"
              bg-variant="danger"
              tag="article"
              class="mb-2 card">
        <p class="card-text">
          This is a job on failure state.
        </p>
      </b-card>
      <b-card title="Job Name"
              bg-variant="secondary"
              tag="article"
              class="mb-2 card">
        <p class="card-text">
          This is a job with no run.
        </p>
      </b-card>
      <router-link class="new-job" :to="{path: '/jobs/creation', exact: true}">
        <i class="fa fa-plus-circle" aria-hidden="true"></i>
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
}
</script>

<style scoped>
.jobs {
  display: flex;
  justify-content: safe center;
}
.jobs .card {
  color: #cccccc;
  margin: 0 1rem 0 1rem;
  width: 200px;
  height: 150px;
}
.jobs .card:hover {
  color: #ffffff;
  cursor: pointer;
}
.jobs .new-job {
  color: #cccccc;
  font-size: 75px;
  align-self: center;
  margin-left: 50px;
}
.jobs .new-job:hover {
  color: #ffffff;
  cursor: pointer;
}
</style>
