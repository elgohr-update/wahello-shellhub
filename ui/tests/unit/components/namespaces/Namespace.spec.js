import Vuex from 'vuex';
import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuetify from 'vuetify';
import Namespace from '@/components/namespace/Namespace';

describe('Namespace', () => {
  const localVue = createLocalVue();
  const vuetify = new Vuetify();
  localVue.use(Vuex);

  let wrapper;

  const inANamespace = true;

  const namespace = {
    name: 'namespace3',
    owner: 'user1',
    member_names: ['user6', 'user7', 'user8'],
    tenant_id: 'e359bf484715',
  };

  const store = new Vuex.Store({
    namespaced: true,
    state: {
      namespace,
    },
    getters: {
      'namespaces/get': (state) => state.namespace,
    },
    actions: {
      'namespaces/fetch': () => {},
      'namespaces/get': () => {},
      'namespaces/switchNamespace': () => {},
      'namespaces/setOwnerStatus': () => {},
      'snackbar/showSnackbarErrorLoading': () => {},
      'snackbar/showSnackbarErrorAssociation': () => {},
    },
  });

  ///////
  // In this case, check owner fields rendering in enterprise version.
  ///////

  describe('Enterprise version', () => {
    beforeEach(() => {
      wrapper = shallowMount(Namespace, {
        store,
        localVue,
        stubs: ['fragment', 'router-link'],
        propsData: { inANamespace: !inANamespace },
        mocks: {
          $env: {
            isEnterprise: true,
          },
        },
        vuetify,
      });

      jest.spyOn(Storage.prototype, 'getItem').mockReturnValue(namespace.tenant_id);
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue Instance', () => {
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Receives data in props', () => {
      expect(wrapper.vm.inANamespace).toEqual(!inANamespace);
    });
    it('Compare data with default value', () => {
      expect(wrapper.vm.inANamespace).toEqual(false);
    });
    it('Process data in the computed', () => {
      expect(wrapper.vm.namespace).toEqual(namespace);
      expect(wrapper.vm.hasNamespace).toEqual(true);
      expect(wrapper.vm.tenant).toEqual(namespace.tenant_id);
    });

    //////
    // HTML validation
    //////

    it('Renders the template with components', () => {
      expect(wrapper.find('[data-test="namespaceList-component"]').exists()).toEqual(true);
      expect(wrapper.find('[data-test="namespaceAdd-component"]').exists()).toEqual(true);
      expect(wrapper.find('[data-test="namespaceAddNoNamespace-component"]').exists()).toEqual(false);
    });
  });

  ///////
  // In this case, check owner fields rendering in open version
  // of the template.
  ///////

  describe('Open version', () => {
    beforeEach(() => {
      wrapper = shallowMount(Namespace, {
        store,
        localVue,
        stubs: ['fragment', 'router-link'],
        propsData: { inANamespace: !inANamespace },
        mocks: {
          $env: {
            isEnterprise: false,
          },
        },
        vuetify,
      });

      jest.spyOn(Storage.prototype, 'getItem').mockReturnValue(namespace.tenant_id);
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue Instance', () => {
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Receives data in props', () => {
      expect(wrapper.vm.inANamespace).toEqual(!inANamespace);
    });
    it('Compare data with default value', () => {
      expect(wrapper.vm.inANamespace).toEqual(false);
    });
    it('Process data in the computed', () => {
      expect(wrapper.vm.namespace).toEqual(namespace);
      expect(wrapper.vm.hasNamespace).toEqual(true);
      expect(wrapper.vm.tenant).toEqual(namespace.tenant_id);
    });

    //////
    // HTML validation
    //////

    it('Renders the template with components', () => {
      expect(wrapper.find('[data-test="namespaceList-component"]').exists()).toEqual(true);
      expect(wrapper.find('[data-test="namespaceAdd-component"]').exists()).toEqual(false);
      expect(wrapper.find('[data-test="namespaceAddNoNamespace-component"]').exists()).toEqual(false);
    });
  });
});
