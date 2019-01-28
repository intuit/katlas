import React from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter } from 'react-router-dom';
import { mount } from 'enzyme';
import fetchMock from 'fetch-mock';

import {AutoHideDuration} from './Notifier';
import App from '../app/App';
import ApiService from '../../services/ApiService';
import {QUERY_LEN_ERR} from "../../utils/errors";
import '../../../public/conf.js'; // to import the configuration

//Use the real store rather than mock Store to keep consistent with the same that is used by HttpService.
import store from '../../store.js';

function sleep (time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

it('Notifier does not open snackbar for valid query', () => {

  const SEARCH_STR = 'foobar';

  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App/>
      </MemoryRouter>
    </Provider>);

  //Check for id of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('#snackbar-message-id').length).toEqual(0);

  let input = wrapper.find('input').last();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});

  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.notify.msg).toEqual('');

  //Check for id of Notifier Snack Bar. Present now.
  expect(wrapper.find('#snackbar-message-id').length).toEqual(0);
});

it('Notifier can open snackbar for empty query', () => {

  const SEARCH_STR = '';

  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App/>
      </MemoryRouter>
    </Provider>);

  //Check for id of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('#snackbar-message-id').length).toEqual(0);

  let input = wrapper.find('input').last();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});

  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.notify.msg).toEqual(QUERY_LEN_ERR);

  //Check for id of Notifier Snack Bar. Present now.
  expect(wrapper.find('#snackbar-message-id').length).toEqual(1);
});

it('Notifier can open snackbar for short query', () => {

    const SEARCH_STR = 'f';

    const wrapper = mount(
      <Provider store={store}>
        <MemoryRouter>
          <App/>
        </MemoryRouter>
      </Provider>);

    //Check for id of Notifier Snack Bar. Initially not present.
    expect(wrapper.find('#snackbar-message-id').length).toEqual(0);

    let input = wrapper.find('input').last();
    input.simulate('change', { target: { value: SEARCH_STR}});
    //and try to submit the query
    input.simulate('keypress', {key: 'Enter'});

    //check for expected state of the store
    const nowStore = store.getState();
    expect(nowStore.notify.msg).toEqual(QUERY_LEN_ERR);

    //Check for id of Notifier Snack Bar. Present now.
    expect(wrapper.find('#snackbar-message-id').length).toEqual(1);
});

it('Notifier can open snackbar for errors returned from httpService', () => {

  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App/>
      </MemoryRouter>
    </Provider>);

  //Check for id of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('#snackbar-message-id').length).toEqual(0);

  const response = {
    body: {},
    status: 400
  };

  fetchMock.get('*', response, { overwriteRoutes: true });

  return ApiService.getQSLResult("Cluster").then((response) => {
    expect(fetchMock.called).toBeTruthy();
    expect(response).toEqual(null);

    //check for expected state of the store
    const nowStore = store.getState();
    expect(nowStore.notify.msg).not.toEqual('');
    expect(nowStore.notify.timestamp).not.toEqual(0);
    let timeDiff = +new Date() - nowStore.notify.timestamp;
    //Check timeDiff close to current time.
    expect(timeDiff).toBeLessThan(AutoHideDuration);
  });

});


it('Notifier can close snackbar after AutoHideDuration', () => {

  let input;
  const SEARCH_STR = '';

  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App/>
      </MemoryRouter>
    </Provider>);

  //Check for id of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('#snackbar-message-id').length).toEqual(0);

  input = wrapper.find('input').last();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});

  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.notify.msg).toEqual(QUERY_LEN_ERR);

  //Check for id of Notifier Snack Bar. Present now.
  expect(wrapper.find('#snackbar-message-id').length).toEqual(1);

  sleep(AutoHideDuration).then(() => {
    //Check for id of Notifier Snack Bar. Not present now as we are over the AutoHideDuration.
    //TODO:SS - this expectation isn't actually being executed during the lifespan of the test, if I introduce the "done()" async feature of the test, it exceeds jest timeout. overall, I think you'll need to fix the test to simulate the time rather than actually sleep for it
    expect(wrapper.find('#snackbar-message-id').length).toEqual(0);
  });

});
