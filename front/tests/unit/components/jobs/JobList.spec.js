import { mount, createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import VueRouter from 'vue-router';
import * as jobsRequests from '@/requests/jobs';
import flushPromises from 'flush-promises'
import { FetchError } from '@/requests/utils'
import JobList from '@/components/jobs/JobList.vue'
import JobItem from '@/components/jobs/JobItem.vue'

const localVue = createLocalVue();
localVue.use(BootstrapVue);
localVue.use(VueRouter);

const router = new VueRouter({
  routes: []
});

const createJobList = propsData => mount(JobList, { propsData, localVue, router });

describe('JobList', () => {
  beforeAll(() => {
    jobsRequests.fetchJobs = jest.fn();
  });

  beforeEach(() => {
    jobsRequests.fetchJobs.mockReset();
  });

  afterAll(() => {
    jobsRequests.fetchJobs.mockRestore();
  });

  it('should load all jobs when mount', async () => {
    // given
    const jobs = [{name: "job 1"}, {name: 'job 2'}]
    jobsRequests.fetchJobs.mockResolvedValue(jobs);
    const cmp = createJobList();

    // when
    await flushPromises();

    // then
    expect(cmp.find('#fetch-jobs-failure-msg').exists()).toBeFalsy();
    expect(jobsRequests.fetchJobs.mock.calls.length).toBe(1);
    expect(cmp.findAll(JobItem)).toHaveLength(2);
  });

  it('should show an error message when fetch jobs fail', async () => {
    // given
    const jobs = [{name: "job 1"}, {name: 'job 2'}]
    jobsRequests.fetchJobs.mockRejectedValue(new FetchError('some error', 400));
    const cmp = createJobList();

    // when
    await flushPromises();

    // then
    expect(jobsRequests.fetchJobs.mock.calls.length).toBe(1);
    expect(cmp.findAll(JobItem)).toHaveLength(0);
    expect(cmp.find('#fetch-jobs-failure-msg').exists()).toBeTruthy();
    expect(cmp.find('#fetch-jobs-failure-msg').text()).toBe('An error has happend when fetching the jobs : some error');
  });
});
