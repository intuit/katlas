import React from 'react';
import {render, unmountComponentAtNode} from 'react-dom';
import {MemoryRouter} from 'react-router-dom';
import {mount, shallow} from 'enzyme';

import GraphContainer from './GraphContainer';
//next import will load envVars from local override of app/public/conf.js
import '../../../public/conf';

jest.useFakeTimers();

const div = document.createElement('div');

it('shallow renders graph container', () => {
  shallow(<GraphContainer/>);
});

it('deep renders graph container', () => {
  render(
    <MemoryRouter>
      <GraphContainer/>
    </MemoryRouter>, div);
  unmountComponentAtNode(div);
});

// it('deep renders graph container while making async request', () => {
//   //mock out the global fetch, don't actually trigger XHR
//   window.fetch = jest.fn().mockImplementation(() =>
//     Promise.resolve(new Response(MOCK_RESP, {status:200})));
//   //render component with path which will cause data fetch
//   render(
//     <MemoryRouter initialEntries={['/graph/0x0']}>
//       <GraphContainer/>
//     </MemoryRouter>, div);
//   //force timers to complete, so as to trigger request
//   jest.runOnlyPendingTimers();
//   unmountComponentAtNode(div);
//   //ensure that fetch was called at least once
//   expect(window.fetch.mock.calls.length).toBeGreaterThan(0);
// });

it('shows a spinner during outstanding request', () => {
  const wrapper = mount(
    <MemoryRouter>
      <GraphContainer/>
    </MemoryRouter>);
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