import React from 'react';
import { Provider } from 'react-redux';
import { render, unmountComponentAtNode } from 'react-dom';
import { MemoryRouter } from 'react-router-dom';
import { mount, shallow } from 'enzyme';

import App from './App';
import configureStore from '../../store/configureStore';
//next import will load envVars from local override of app/public/conf.js
import '../../../public/conf';

const div = document.createElement('div');
const store = configureStore();

it('shallow renders app', () => {
  shallow(<App />);
});

it('deep renders home view', () => {
  render(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    </Provider>, div);
  unmountComponentAtNode(div);
});

it('deep renders results view', () => {
  render(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/results']}>
        <App />
      </MemoryRouter>
    </Provider>, div);
  unmountComponentAtNode(div);
});

it('deep renders graph view', () => {
  render(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/graph/0x0']}>
        <App />
      </MemoryRouter>
    </Provider>, div);
  unmountComponentAtNode(div);
});

xit('ensures that search bar input text is equal between menubar and home view', () => {
  const SEARCH_STR_A = 'foobar';
  const SEARCH_STR_B = 'bazqux';
  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App />
      </MemoryRouter>
    </Provider>);

  //1 input element in menu and 1 in home view
  expect(wrapper.find('input')).toHaveLength(2);

  //change search text A in menu bar input
  let inputMB = wrapper.find('input').at(0);
  inputMB.simulate('change', { target: { value: SEARCH_STR_A}});
  let inputH = wrapper.find('input').at(1);
  //confirm that it now appears in home input element
  expect(inputH.props().value).toEqual(SEARCH_STR_A);
  //and in the store
  expect(store.getState().query.current).toEqual(SEARCH_STR_A);
  //change in home input
  inputH.simulate('change', { target: { value: SEARCH_STR_B}});
  //re-capture menu bar input and confirm change there
  inputMB = wrapper.find('input').at(0);
  expect(inputMB.props().value).toEqual(SEARCH_STR_B);
  //and store
  expect(store.getState().query.current).toEqual(SEARCH_STR_B);
});