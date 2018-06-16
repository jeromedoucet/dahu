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
      <i class="fa fa-plus-circle" aria-hidden="true"></i>
    </div>
  </div>
</template>

<script>
export default {
}
</script>

<style scoped>
.jobs {
  margin-top: 3rem;
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
.jobs .fa-plus-circle {
  color: #cccccc;
  font-size: 75px;
  align-self: center;
  margin-left: 50px;
}
.jobs .fa-plus-circle:hover {
  color: #ffffff;
  cursor: pointer;
}
</style>
