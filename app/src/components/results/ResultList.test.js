import React from 'react';
import { Provider } from 'react-redux';
import { render, unmountComponentAtNode } from 'react-dom';
import { MemoryRouter } from 'react-router-dom';
import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';

import initialState from '../../reducers/initialState';
import ResultList from './ResultList';

const div = document.createElement('div');
const middlewares = [thunk];
const mockStore = configureMockStore(middlewares);

describe('query result', () => {
  it('can load initial state', () => {
    const store = mockStore(initialState);
    render(
      <Provider store={store}>
        <MemoryRouter>
          <ResultList query={store.getState().query} />
        </MemoryRouter>
      </Provider>,
      div
    );
    unmountComponentAtNode(div);
  });

  it('can compose the result table layout', () => {
    const queryStr = 'deployment{*}.replicaset[@count(pod)<3]{*}.pod{*}';
    const metadata = {
      deployment: {
        uid: '0xfb80d',
        name: 'deployment',
        fields: [
          {
            fieldname: 'labels',
            fieldtype: 'json',
            mandatory: false,
            cardinality: 'one'
          },
          {
            fieldname: 'name',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'numreplicas',
            fieldtype: 'int',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'resourceversion',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'strategy',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'creationtime',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'cluster',
            fieldtype: 'relationship',
            mandatory: true,
            refdatatype: 'cluster',
            cardinality: 'one'
          },
          {
            fieldname: 'availablereplicas',
            fieldtype: 'int',
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
            fieldname: 'namespace',
            fieldtype: 'relationship',
            mandatory: true,
            refdatatype: 'namespace',
            cardinality: 'one'
          },
          {
            fieldname: 'resourceid',
            fieldtype: 'string',
            mandatory: false,
            cardinality: 'one'
          },
          {
            fieldname: 'application',
            fieldtype: 'relationship',
            mandatory: false,
            refdatatype: 'application',
            cardinality: 'many'
          },
          {
            fieldname: 'k8sobj',
            fieldtype: 'true',
            mandatory: true,
            cardinality: 'one'
          }
        ]
      },
      replicaset: {
        uid: '0xfb821',
        name: 'replicaset',
        fields: [
          {
            fieldname: 'resourceid',
            fieldtype: 'string',
            mandatory: false,
            cardinality: 'one'
          },
          {
            fieldname: 'cluster',
            fieldtype: 'relationship',
            mandatory: true,
            refdatatype: 'cluster',
            cardinality: 'one'
          },
          {
            fieldname: 'k8sobj',
            fieldtype: 'true',
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
          },
          {
            fieldname: 'objtype',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'namespace',
            fieldtype: 'relationship',
            mandatory: true,
            refdatatype: 'namespace',
            cardinality: 'one'
          },
          {
            fieldname: 'owner',
            fieldtype: 'relationship',
            mandatory: false,
            refdatatype: 'deployment',
            cardinality: 'one'
          },
          {
            fieldname: 'numreplicas',
            fieldtype: 'int',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'podspec',
            fieldtype: 'json',
            mandatory: false,
            cardinality: 'one'
          }
        ]
      },
      pod: {
        uid: '0xfb858',
        name: 'pod',
        fields: [
          {
            fieldname: 'name',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'phase',
            fieldtype: 'string',
            mandatory: false,
            cardinality: 'one'
          },
          {
            fieldname: 'ip',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'containers',
            fieldtype: 'json',
            mandatory: false,
            cardinality: 'many'
          },
          {
            fieldname: 'starttime',
            fieldtype: 'string',
            mandatory: false,
            cardinality: 'one'
          },
          {
            fieldname: 'resourceid',
            fieldtype: 'string',
            mandatory: false,
            cardinality: 'one'
          },
          {
            fieldname: 'cluster',
            fieldtype: 'relationship',
            mandatory: true,
            refdatatype: 'cluster',
            cardinality: 'one'
          },
          {
            fieldname: 'namespace',
            fieldtype: 'relationship',
            mandatory: true,
            refdatatype: 'namespace',
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
            fieldname: 'nodename',
            fieldtype: 'relationship',
            mandatory: true,
            refdatatype: 'node',
            cardinality: 'one'
          },
          {
            fieldname: 'ownertype',
            fieldtype: 'string',
            mandatory: false,
            cardinality: 'one'
          },
          {
            fieldname: 'creationtime',
            fieldtype: 'string',
            mandatory: true,
            cardinality: 'one'
          },
          {
            fieldname: 'k8sobj',
            fieldtype: 'true',
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
            fieldname: 'volumes',
            fieldtype: 'json',
            mandatory: false,
            cardinality: 'many'
          },
          {
            fieldname: 'owner',
            fieldtype: 'relationship',
            mandatory: false,
            refdatatype: 'replicaset,daemonset,statefulset',
            cardinality: 'one'
          }
        ]
      }
    };
    const store = mockStore({
      ...initialState,
      query: {
        ...initialState.query,
        current: queryStr,
        metadata: metadata
      }
    });
    render(
      <Provider store={store}>
        <MemoryRouter>
          <ResultList query={store.getState().query} />
        </MemoryRouter>
      </Provider>,
      div
    );
    unmountComponentAtNode(div);
  });
});
