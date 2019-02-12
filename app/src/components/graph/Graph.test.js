import React from 'react';
import { Provider } from 'react-redux';
import { render, unmountComponentAtNode } from 'react-dom';
import { MemoryRouter } from 'react-router-dom';

import Graph from './Graph';
import configureStore from '../../store/configureStore';

const div = document.createElement('div');
const store = configureStore();

//TODO:DM - extend this test suite, especially wrt legend construction given mock data
it('deep renders graph', () => {
  render(
    <Provider store={store}>
      <MemoryRouter>
        <Graph/>
      </MemoryRouter>
    </Provider>, div);
  unmountComponentAtNode(div);
});
