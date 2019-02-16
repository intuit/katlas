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

});
