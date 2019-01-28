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

  it('should create an action to receive metadaa', () => {
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
  it('should get fetch query with 1 result', done => {
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
      expect(store.getActions()[1].metadata).toEqual(APP_METADATA);
      done();
    }, 3);
  });

  it('should not fetch metadata already cached', done => {
    const store = mockStore({ query: { metadata: { application: APP_METADATA} } });

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
  })
});

//MOCK_RESP contains 2 duplicate objects, this will exercise code to filter out dupes and should leave us with resulting length of 1
const MOCK_RESP_EMPTY = { objects: [] };
const MOCK_RESP_DUPE_LEN = 1;
const MOCK_RESP_DUPE = {
  objects: [
    {
      cluster: [
        {
          k8sobj: 'K8sObj',
          name: 'preprod-west2.cluster.k8s.local',
          objtype: 'Cluster',
          resourceid: 'preprod-west2.cluster.k8s.local',
          resourceversion: '0',
          uid: '0x13d7'
        }
      ],
      clusterip: '100.68.144.37',
      k8sobj: 'K8sObj',
      labels:
        '{"app":"helm-chart","chart":"helm-chart-0.1.0","heritage":"Tiller","release":"profile-testing"}',
      name: 'profile-testing-helm-chart',
      namespace: [
        {
          k8sobj: 'K8sObj',
          labels: 'null',
          name: 'profile-testing',
          objtype: 'Namespace',
          resourceid: 'preprod-west2.cluster.k8s.local:profile-testing',
          resourceversion: '0',
          uid: '0x1e84'
        }
      ],
      objtype: 'Service',
      ports:
        '[{"nodePort":32313,"port":80,"protocol":"TCP","targetPort":8080}]',
      resourceid:
        'preprod-west2.cluster.k8s.local:profile-testing:profile-testing-helm-chart',
      resourceversion: '3574724',
      selector: '{"app":"mysqlserver","release":"profile-testing"}',
      servicetype: 'NodePort',
      uid: '0x4c05'
    },
    {
      cluster: [
        {
          k8sobj: 'K8sObj',
          name: 'preprod-west2.cluster.k8s.local',
          objtype: 'Cluster',
          resourceid: 'preprod-west2.cluster.k8s.local',
          resourceversion: '0',
          uid: '0x13d7'
        }
      ],
      clusterip: '100.68.144.37',
      k8sobj: 'K8sObj',
      labels:
        '{"app":"helm-chart","chart":"helm-chart-0.1.0","heritage":"Tiller","release":"profile-testing"}',
      name: 'profile-testing-helm-chart',
      namespace: [
        {
          k8sobj: 'K8sObj',
          labels: 'null',
          name: 'profile-testing',
          objtype: 'Namespace',
          resourceid: 'preprod-west2.cluster.k8s.local:profile-testing',
          resourceversion: '0',
          uid: '0x1e84'
        }
      ],
      objtype: 'Service',
      ports:
        '[{"nodePort":32313,"port":80,"protocol":"TCP","targetPort":8080}]',
      resourceid:
        'preprod-west2.cluster.k8s.local:profile-testing:profile-testing-helm-chart',
      resourceversion: '3574724',
      selector: '{"app":"mysqlserver","release":"profile-testing"}',
      servicetype: 'NodePort',
      uid: '0x4c05'
    }
  ]
};

const APP_METADATA = {
  uid: '0xfb7eb',
  name: 'application',
  fields: [
    {
      fieldname: 'k8sobj',
      fieldtype: 'string',
      mandatory: true,
      cardinality: 'one'
    },
    {
      fieldname: 'objtype',
      fieldtype: 'string',
      mandatory: true,
      cardinality: 'one'
    },
    {
      fieldname: 'name',
      fieldtype: 'string',
      mandatory: true,
      cardinality: 'one'
    },
    {
      fieldname: 'resourceid',
      fieldtype: 'string',
      mandatory: false,
      cardinality: 'one'
    },
    {
      fieldname: 'labels',
      fieldtype: 'json',
      mandatory: false,
      cardinality: 'one'
    },
    {
      fieldname: 'resourceversion',
      fieldtype: 'string',
      mandatory: true,
      cardinality: 'one'
    },
    {
      fieldname: 'creationtime',
      fieldtype: 'string',
      mandatory: true,
      cardinality: 'one'
    }
  ]
};
