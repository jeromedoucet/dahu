import { mount } from '@vue/test-utils';
import { createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import ButtonSpin from '@/components/controls/ButtonSpin.vue'

const localVue = createLocalVue();
localVue.use(BootstrapVue);


const createButtonSpin = propsData => mount(ButtonSpin, { propsData, localVue });

describe('ButtonSpin.vue', () => {

  describe('props validation', () => {
    const cmp = createButtonSpin({label: 'something'});

    describe('variant', () => {
      const variant = cmp.vm.$options.props.variant;

      it('should be string type', () => {
        expect(variant.type).toBe(String);
      });

      it('should accept primary value', () => {
        expect(variant.validator && variant.validator('primary')).toBeTruthy() ;
      });

      it('should accept secondary value', () => {
        expect(variant.validator && variant.validator('secondary')).toBeTruthy(); 
      });

      it('should accept success value', () => {
        expect(variant.validator && variant.validator('success')).toBeTruthy();
      });

      it('should accept warning value', () => {
        expect(variant.validator && variant.validator('warning')).toBeTruthy();
      });

      it('should accept danger value', () => {
        expect(variant.validator && variant.validator('danger')).toBeTruthy();
      });

      it('should accept link value', () => {
        expect(variant.validator && variant.validator('link')).toBeTruthy();
      });

      it('should reject another value', () => {
        expect(variant.validator && variant.validator('hello')).toBeFalsy();
      });
    });

    describe('spinning', () => {
      const spinning = cmp.vm.$options.props.spinning;

      it('should be Boolean type', () => {
        expect(spinning.type).toBe(Boolean);
      });

      it('should be false by default', () => {
        expect(spinning.default).toBeFalsy();
      });
    });

    describe('type', () => {
      const type = cmp.vm.$options.props.type;

      it('should be String type', () => {
        expect(type.type).toBe(String);
      });   

      it('should has button as default value', () => {
        expect(type.default).toBe('button');
      });

      it('should accept button as value', () => {
        expect(type.validator && type.validator('button')).toBeTruthy();
      });

      it('should accept submit as value', () => {
        expect(type.validator && type.validator('submit')).toBeTruthy();
      });

      it('should reject another value', () => {
        expect(type.validator && type.validator('hello')).toBeFalsy();
      });
    });
  });

  describe('rendering', () => {

    it('should render the label correctly', () => {
      // given
      const label = 'some-button-label';

      // when
      const cmp = createButtonSpin({label});

      // then
      expect(cmp.find('button').text()).toBe(label);
    });

    it('should have button as default value for type attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something'});

      // then
      expect(cmp.find('button').attributes().type).toBe('button');
    });

    it('should accept submit value for type attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something', type: 'submit'});

      // then
      expect(cmp.find('button').attributes().type).toBe('submit');
    });

    it('should have primary as default value for variant attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something'});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['btn-primary']));
    });

    it('should accept secondary as variant attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something', variant: 'secondary'});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['btn-secondary']));
    });

    it('should accept success as variant attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something', variant: 'success'});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['btn-success']));
    });

    it('should accept warning as variant attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something', variant: 'warning'});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['btn-warning']));
    });

    it('should accept danger as variant attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something', variant: 'danger'});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['btn-danger']));
    });

    it('should accept link as variant attribute', () => {
      // when
      const cmp = createButtonSpin({label: 'something', variant: 'link'});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['btn-link']));
    });

    it('should add a spinner and disable the button when spinning is true', () => {
      // when
      const cmp = createButtonSpin({label: 'something', spinning: true});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['disabled']));
      expect(cmp.find('i').classes()).toEqual(expect.arrayContaining(['fa', 'fa-spinner']));
    });

    it('should disable the button without rendering the button when disable is true', () => {
      // when
      const cmp = createButtonSpin({label: 'something', disabled: true});

      // then
      expect(cmp.find('button').classes()).toEqual(expect.arrayContaining(['disabled']));
      expect(cmp.find('i').exists()).toBeFalsy();
    });
  });
});
