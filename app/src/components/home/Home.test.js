import React from 'react';
import { Provider } from 'react-redux';
import { render, unmountComponentAtNode } from 'react-dom';
import { MemoryRouter } from 'react-router-dom';
import { mount, shallow } from 'enzyme';

import configureStore from '../../store/configureStore';
import Home from './Home';
//next import will load envVars from local override of app/public/conf.js
import '../../../public/conf';

const div = document.createElement('div');
let store;
beforeEach(() => {
  store = configureStore();
});

it('deep renders home component', () => {
  render(
    <Provider store={store}>
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    </Provider>, div);
  unmountComponentAtNode(div);
});

it('has one input element', () => {
  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    </Provider>);
  expect(wrapper.find('input')).toHaveLength(1);
});

xit('submits a valid query', () => {
  const SEARCH_STR = 'foobar';
  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    </Provider>);

  //change search text in menu bar input
  let input = wrapper.find('input').first();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and submit the query
  input.simulate('keypress', {key: 'Enter'});
  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.query.current).toEqual(SEARCH_STR);
  expect(nowStore.query.submitted).toEqual(true);
  expect(nowStore.query.isWaiting).toEqual(true);

});

it('tries to submit an empty query', () => {
  const SEARCH_STR = '';
  const NOTIFY_MSG = 'Minimum length of Search word must be 3 characters.';
  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    </Provider>);

  //change search text in menu bar input
  let input = wrapper.find('input').first();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});
  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.query.current).toEqual(SEARCH_STR);
  expect(nowStore.query.submitted).toEqual(false);
  expect(nowStore.query.isWaiting).toEqual(false);
  expect(nowStore.notify.msg).toEqual(NOTIFY_MSG);
});

it('tries to submit too short a query', () => {
  const SEARCH_STR = 'fo';
  const NOTIFY_MSG = 'Minimum length of Search word must be 3 characters.';
  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    </Provider>);

  //change search text in menu bar input
  let input = wrapper.find('input').first();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});
  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.query.current).toEqual(SEARCH_STR);
  expect(nowStore.query.submitted).toEqual(false);
  expect(nowStore.query.isWaiting).toEqual(false);
  expect(nowStore.notify.msg).toEqual(NOTIFY_MSG);
  //TODO:DM - also check for existience of notification here or split into separate test; search wrapper for a DOM node or CSS class that changes when notification is shown
});

it('tries to submit with a tab', () => {
  const SEARCH_STR = 'fo';
  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    </Provider>);

  //change search text in menu bar input
  let input = wrapper.find('input').first();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and hit tab key rather than submit to ensure NON enter key presses are
  //handled correctly
  input.simulate('keypress', {key: 'Tab'});
  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.query.current).toEqual(SEARCH_STR);
  expect(nowStore.query.submitted).toEqual(false);
  expect(nowStore.query.isWaiting).toEqual(false);
});
