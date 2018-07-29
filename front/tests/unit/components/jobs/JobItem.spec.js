import { mount, createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import VueRouter from 'vue-router';
import flushPromises from 'flush-promises'
import JobItem from '@/components/jobs/JobItem.vue'

const localVue = createLocalVue();
localVue.use(BootstrapVue);
localVue.use(VueRouter);

const router = new VueRouter({
  routes: []
});

const createJobItem = propsData => mount(JobItem, { propsData, localVue, router });

describe('JobItem', () => {

  it('should redirect to job pipelline when click on element', async () => {
    // given
    const job = {id: 1, name: "job 1"};
    const cmp = createJobItem({job});
    cmp.vm.$router.push = jest.fn();

    // when
    cmp.find('.job-card').trigger('click');
    await flushPromises();

    // then
    expect(cmp.vm.$router.push.mock.calls.length).toBe(1);
    expect(cmp.vm.$router.push.mock.calls[0][0]).toBe('/jobs/1/pipeline');
  });
});
