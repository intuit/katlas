import React from 'react';
import { Provider } from 'react-redux';
import { render, unmountComponentAtNode } from 'react-dom';
import { MemoryRouter } from 'react-router-dom';
import { mount } from 'enzyme';

import GraphContainer from './GraphContainer';
import configureStore from '../../store/configureStore';
//next import will load envVars from local override of app/public/conf.js
import '../../../public/conf';

jest.useFakeTimers();

const div = document.createElement('div');
const store = configureStore();

it('deep renders graph container', () => {
  render(
    <Provider store={store}>
      <MemoryRouter>
        <GraphContainer/>
      </MemoryRouter>
    </Provider>, div);
  unmountComponentAtNode(div);
});

xit('deep renders graph container while making async request', (done) => {
  //mock out the global fetch, don't actually trigger XHR
  window.fetch = jest.fn().mockImplementation(() =>
    Promise.resolve(new Response(JSON.stringify(MOCK_RESP), {status:200})));
  //render component with path which will cause data fetch
  render(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/graph/0x0']}>
        <GraphContainer/>
      </MemoryRouter>
    </Provider>, div);
  //force timers to complete, so as to trigger request
  jest.runOnlyPendingTimers();
  //ensure that fetch was called at least once
  expect(window.fetch.mock.calls.length).toBeGreaterThan(0);
  unmountComponentAtNode(div);
  done();
});

xit('shows a spinner during outstanding request', () => {
  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <GraphContainer/>
      </MemoryRouter>
    </Provider>);
  wrapper.setState({waitingOnReq: true});
  expect(wrapper.find('CircularProgress').length).toBe(1);
});


const MOCK_RESP = {
  "objects": [{
    "cluster": [{
      "k8sobj": "K8sObj",
      "name": "preprod-west2.cluster.k8s.local",
      "objtype": "Cluster",
      "resourceid": "preprod-west2.cluster.k8s.local",
      "resourceversion": "0",
      "uid": "0x13d7"
    }],
    "clusterip": "100.68.144.37",
    "k8sobj": "K8sObj",
    "labels": "{\"app\":\"helm-chart\",\"chart\":\"helm-chart-0.1.0\",\"heritage\":\"Tiller\",\"release\":\"profile-testing\"}",
    "name": "profile-testing-helm-chart",
    "namespace": [{
      "k8sobj": "K8sObj",
      "labels": "null",
      "name": "profile-testing",
      "objtype": "Namespace",
      "resourceid": "preprod-west2.cluster.k8s.local:profile-testing",
      "resourceversion": "0",
      "uid": "0x1e84"
    }],
    "objtype": "Service",
    "ports": "[{\"nodePort\":32313,\"port\":80,\"protocol\":\"TCP\",\"targetPort\":8080}]",
    "resourceid": "preprod-west2.cluster.k8s.local:profile-testing:profile-testing-helm-chart",
    "resourceversion": "3574724",
    "selector": "{\"app\":\"mysqlserver\",\"release\":\"profile-testing\"}",
    "servicetype": "NodePort",
    "uid": "0x4c05"
  }]
};