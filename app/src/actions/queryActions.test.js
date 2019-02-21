import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';

import * as actions from './queryActions';
import * as types from './actionTypes';
//next import will load envVars from local override of app/public/conf.js
import '../../public/conf';

const middlewares = [thunk];
const mockStore = configureMockStore(middlewares);

describe('query actions', () => {
  xit('should create an action to submit query', () => {
    const expectedAction = {
      type: types.SUBMIT_QUERY
    };
    expect(actions.submitQuery()).toEqual(expectedAction);
  });

  it('should create an action to receieve query', () => {
    const results = { foo: 'bar' };
    const expectedAction = {
      type: types.RECEIVE_QUERY,
      results
    };
    expect(actions.receiveQuery(results)).toEqual(expectedAction);
  });

  it('should create an action to request metadata', () => {
    const objType = 'deployment';
    const expectedAction = {
      type: types.REQUEST_METADATA,
      objType
    };
    expect(actions.requestMetadata(objType)).toEqual(expectedAction);
  });

  it('should create an action to receive metadata', () => {
    const objType = 'application';
    const expectedAction = {
      type: types.RECEIVE_METADATA,
      objType,
      metadata: APP_METADATA
    };
    expect(actions.receiveMetadata(objType, APP_METADATA)).toEqual(
      expectedAction
    );
  });
});

describe('asynch query actions', () => {
  it('should get fetch query with results', done => {
    //mock fetch so that API call will return immediately with mock data
    window.fetch = jest
      .fn()
      .mockImplementation(() =>
        Promise.resolve(
          new Response(JSON.stringify(MOCK_RESP_DUPE), { status: 200 })
        )
      );

    const store = mockStore({ query: { metadata: {} } });

    const fn = actions.fetchQuery('doesntMatter');
    expect(fn).toBeInstanceOf(Function);
    fn(store.dispatch, store.getState);

    //now we're expecting the action to have been triggered with particular
    //checking the length because it's not otherwise so easy to compare raw
    //network response vs. parsed data, will test that transformation logic
    //in reducer tests
    setTimeout(() => {
      expect(store.getActions()[1].results).toHaveLength(MOCK_RESP_DUPE_LEN);
      done();
    }, 3);
  });

  it('should be able to handle fetch query with empty result', done => {
    //re-mock fetch with diff mock value, an empty response
    window.fetch = jest
      .fn()
      .mockImplementation(() =>
        Promise.resolve(
          new Response(JSON.stringify(MOCK_RESP_EMPTY), { status: 200 })
        )
      );

    const store = mockStore({ query: { metadata: {} } });
    const fn = actions.fetchQuery('doesntMatter');
    expect(fn).toBeInstanceOf(Function);
    fn(store.dispatch, store.getState);

    //each fetchQuery will trigger 2 actions. first is requestQuery and second is receiveQuery which has the results
    //now we're looking for receiveQuery (the 3nd as there was another fetchQuery earlier) action to have been triggered
    setTimeout(() => {
      expect(store.getActions()[1].results).toHaveLength(0);
      done();
    }, 3);
  });

  it('should handle fetch metadata', done => {
    const store = mockStore({ query: { metadata: {} } });

    window.fetch = jest
      .fn()
      .mockImplementation(() =>
        Promise.resolve(
          new Response(JSON.stringify(APP_METADATA), { status: 200 })
        )
      );

    const fn = actions.fetchMetadata('application');
    expect(fn).toBeInstanceOf(Function);
    fn(store.dispatch, store.getState);

    setTimeout(() => {
      expect(store.getActions()).toHaveLength(2);
      expect(store.getActions()[1].objType).toEqual('application');
      expect(store.getActions()[1].metadata).toEqual(APP_METADATA.objects[0]);
      done();
    }, 3);
  });

  it('should not fetch metadata already cached', done => {
    const store = mockStore({
      query: { metadata: { application: APP_METADATA } }
    });

    window.fetch = jest
      .fn()
      .mockImplementation(() =>
        Promise.resolve(
          new Response(JSON.stringify(APP_METADATA), { status: 200 })
        )
      );

    const fn = actions.fetchMetadata('application');
    expect(fn).toBeInstanceOf(Function);
    fn(store.dispatch, store.getState);

    setTimeout(() => {
      expect(store.getActions()).toHaveLength(0);
      done();
    }, 3);
  });
});

//MOCK_RESP contains 2 duplicate objects, this will exercise code to filter out dupes and should leave us with resulting length of 1
const MOCK_RESP_EMPTY = { objects: [] };
const MOCK_RESP_DUPE_LEN = 2;
const MOCK_RESP_DUPE = {
  count: 2,
  objects: [
    {
      cluster: [
        {
          k8sobj: 'k8sobj',
          name: 'a.cluster.k8s.local',
          objtype: 'cluster',
          resourceid: 'cluster:a.cluster.k8s.local',
          resourceversion: '0',
          uid: '0x15fa0c'
        }
      ],
      creationtime: '2019-01-24T01:06:47Z',
      ip: '100.113.249.235',
      k8sobj: 'k8sobj',
      labels: '{"app":"cmk-controller","pod-template-hash":"549037671"}',
      name: 'cmk-controller-prd-qal-98f47cbc5-wf78d',
      namespace: [
        {
          creationtime: '2019-01-17T22:52:00Z',
          k8sobj: 'k8sobj',
          labels:
            '{"foobar.com/owner":"iksm","foobar.com/prune-label":"foo-controller-use2","name":"foo-controller-use2"}',
          name: 'foo-controller-use2',
          objtype: 'namespace',
          resourceid: 'namespace:a.cluster.k8s.local:foo-controller-use2',
          resourceversion: '38812781',
          uid: '0x1e5e81'
        }
      ],
      nodename: [
        {
          k8sobj: 'k8sobj',
          name: 'ip-10-150-106-103.us-east-2.compute.internal',
          objtype: 'node',
          resourceid:
            'node:a.cluster.k8s.local:ip-10-150-106-103.us-east-2.compute.internal',
          resourceversion: '0',
          uid: '0x1e1062'
        }
      ],
      objtype: 'pod',
      owner: [
        {
          creationtime: '2019-01-24T01:06:47Z',
          k8sobj: 'k8sobj',
          labels: '{"app":"cmk-controller","pod-template-hash":"549037671"}',
          name: 'cmk-controller-prd-qal-98f47cbc5',
          numreplicas: '1',
          objtype: 'replicaset',
          resourceid:
            'replicaset:a.cluster.k8s.local:foo-controller-use2:cmk-controller-prd-qal-98f47cbc5',
          resourceversion: '41162225',
          uid: '0x19f1b1'
        }
      ],
      ownertype: 'replicaset',
      phase: 'Running',
      resourceid:
        'pod:a.cluster.k8s.local:foo-controller-use2:cmk-controller-prd-qal-98f47cbc5-wf78d',
      resourceversion: '41162224',
      starttime: '2019-01-24T01:06:47Z',
      volumes:
        '[{"name":"cmk-controller-token-pjsdf","secret":{"secretName":"cmk-controller-token-pjsdf","defaultMode":420}}]'
    },
    {
      cluster: [
        {
          k8sobj: 'k8sobj',
          name: 'a-west2.cluster.k8s.local',
          objtype: 'cluster',
          resourceid: 'cluster:a-west2.cluster.k8s.local',
          resourceversion: '0',
          uid: '0x175990'
        }
      ],
      creationtime: '2019-01-15T23:44:13Z',
      k8sobj: 'k8sobj',
      labels:
        '{"app":"cutlass-ui","pod-template-hash":"6277035","splunk-index":"k8s_paas"}',
      name: 'cutlass-ui-deployment-b6cc479',
      namespace: [
        {
          creationtime: '2018-09-18T22:01:31Z',
          k8sobj: 'k8sobj',
          labels:
            '{"foobar.com/owner":"iksm","foobar.com/prune-label":"foo-api-usw2-ns","name":"foo-api-usw2-ns"}',
          name: 'foo-api-usw2-ns',
          objtype: 'namespace',
          resourceid: 'namespace:a-west2.cluster.k8s.local:foo-api-usw2-ns',
          resourceversion: '29971588',
          uid: '0x116621'
        }
      ],
      numreplicas: '0',
      objtype: 'replicaset',
      owner: [
        {
          availablereplicas: '1',
          creationtime: '2018-12-20T20:21:48Z',
          k8sobj: 'k8sobj',
          labels: 'null',
          name: 'cutlass-ui-deployment',
          numreplicas: '1',
          objtype: 'deployment',
          resourceid:
            'deployment:a-west2.cluster.k8s.local:foo-api-usw2-ns:cutlass-ui-deployment',
          resourceversion: '93979697',
          strategy: 'RollingUpdate',
          uid: '0x147369'
        }
      ],
      resourceid:
        'replicaset:a-west2.cluster.k8s.local:foo-api-usw2-ns:cutlass-ui-deployment-b6cc479',
      resourceversion: '84921498'
    }
  ],
  status: 200
};

const APP_METADATA = {
  status: 200,
  objects: [
    {
      uid: '0x1f47d1',
      name: 'application',
      objtype: 'metadata',
      fields: [
        {
          uid: '0x2ed85c',
          fieldname: 'name',
          fieldtype: 'string',
          mandatory: true,
          cardinality: 'one'
        },
        {
          uid: '0x2ed85d',
          fieldname: 'resourceid',
          fieldtype: 'string',
          mandatory: false,
          cardinality: 'one'
        },
        {
          uid: '0x2ed85e',
          fieldname: 'labels',
          fieldtype: 'json',
          mandatory: false,
          cardinality: 'one'
        },
        {
          uid: '0x2ed85f',
          fieldname: 'resourceversion',
          fieldtype: 'string',
          mandatory: true,
          cardinality: 'one'
        },
        {
          uid: '0x2ed860',
          fieldname: 'creationtime',
          fieldtype: 'string',
          mandatory: true,
          cardinality: 'one'
        },
        {
          uid: '0x2ed861',
          fieldname: 'k8sobj',
          fieldtype: 'string',
          mandatory: true,
          cardinality: 'one'
        },
        {
          uid: '0x2ed862',
          fieldname: 'objtype',
          fieldtype: 'string',
          mandatory: true,
          cardinality: 'one'
        }
      ],
      resourceversion: '6'
    }
  ]
};
