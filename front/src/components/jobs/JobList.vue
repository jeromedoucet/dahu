<template>
  <div>
    <div class="jobs">
      <job-item v-for="job in jobs" :job="job" :key="job.id"/>
        <router-link class="new-job" :to="{path: '/jobs/creation', exact: true}">
          <new-item/>
        </router-link>
    </div>
    <b-alert 
      id="fetch-jobs-failure-msg" 
      class="test-msg" 
      :show="errorMsg !== ''" 
      variant="danger"
    >
      {{errorMsg}}
    </b-alert>
  </div>
</template>
<script>
import JobItem from '@/components/jobs/JobItem.vue'
import NewItem from '@/components/controls/NewItem.vue'
import { fetchJobs } from '@/requests/jobs';
export default {
  // todo type validation of jobs ?
  components: {
    JobItem,
    NewItem,
  },
  data () {
    return {
      jobs: null,
      errorMsg: '',
    }
  },
  methods: {
    loadJobs: async function () {
      try {
        this.jobs = await fetchJobs();
      } catch (err) {
        this.errorMsg = `An error has happend when fetching the jobs : ${err.message}`;
      } 
    }
  },
  beforeMount: async function() {
    await this.loadJobs();
  }
}
</script>
<style scoped>
.jobs {
  display: flex;
  justify-content: safe center;
}
.jobs .new-job {
  margin-left: 50px;
  align-self: center;
}
</style>
