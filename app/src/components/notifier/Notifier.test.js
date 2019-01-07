import React from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter } from 'react-router-dom';
import { mount } from 'enzyme';
import fetchMock from 'fetch-mock';

import configureStore from '../../store/configureStore';
import {AutoHideDuration} from './Notifier';
import App from '../app/App';
import {HttpService} from '../../services/httpService';
import {QUERY_LEN_ERR} from "../../utils/errors";

//Use the real store rather than mock Store to keep consistent with the same that is used by HttpService.
import store from '../../store.js';

const div = document.createElement('div');

function sleep (time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

//Using hostNodes hack instead of wrapper.find('.Notifier-root-120').length - refer https://github.com/airbnb/enzyme/issues/1253
//expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);

it('Notifier does not open snackbar for valid query', () => {

  const SEARCH_STR = 'foobar';

  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App/>
      </MemoryRouter>
    </Provider>);

  //Check for class name of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);

  let input = wrapper.find('input').last();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});

  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.notify.msg).toEqual('');

  //Check for class name of Notifier Snack Bar. Present now.
  expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);
});

it('Notifier can open snackbar for empty query', () => {

  const SEARCH_STR = '';

  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App/>
      </MemoryRouter>
    </Provider>);

  //Check for class name of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);

  let input = wrapper.find('input').last();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});

  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.notify.msg).toEqual(QUERY_LEN_ERR);

  //Check for class name of Notifier Snack Bar. Present now.
  expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(1);
});

it('Notifier can open snackbar for short query', () => {

    const SEARCH_STR = 'f';

    const wrapper = mount(
      <Provider store={store}>
        <MemoryRouter>
          <App/>
        </MemoryRouter>
      </Provider>);

    //Check for class name of Notifier Snack Bar. Initially not present.
    expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);

    let input = wrapper.find('input').last();
    input.simulate('change', { target: { value: SEARCH_STR}});
    //and try to submit the query
    input.simulate('keypress', {key: 'Enter'});

    //check for expected state of the store
    const nowStore = store.getState();
    expect(nowStore.notify.msg).toEqual(QUERY_LEN_ERR);

    //Check for class name of Notifier Snack Bar. Present now.
    expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(1);
});

it('Notifier can open snackbar for errors returned from httpService', () => {

  const wrapper = mount(
    <Provider store={store}>
      <MemoryRouter>
        <App/>
      </MemoryRouter>
    </Provider>);

  //Check for class name of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);

  let dummyUrl = "http://katlas.com/v1/qsl";
  let dummyParams = {qslstring: 'Cluster'};
  const response = {
    body: {},
    status: 400
  };

  fetchMock.get('*', response, { overwriteRoutes: true });

  return HttpService.get({
    url: dummyUrl,
    params: dummyParams
  }).then((response) => {
    expect(fetchMock.called).toBeTruthy();
    expect(response).toEqual(null);

    //check for expected state of the store
    const nowStore = store.getState();
    console.log('nowstore.notify.type=' + nowStore.notify.type);
    console.log('nowstore.notify.msg=' + nowStore.notify.msg);
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

  //Check for class name of Notifier Snack Bar. Initially not present.
  expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);

  input = wrapper.find('input').last();
  input.simulate('change', { target: { value: SEARCH_STR}});
  //and try to submit the query
  input.simulate('keypress', {key: 'Enter'});

  //check for expected state of the store
  const nowStore = store.getState();
  expect(nowStore.notify.msg).toEqual(QUERY_LEN_ERR);

  //Check for class name of Notifier Snack Bar. Present now.
  expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(1);

  sleep(AutoHideDuration).then(() => {
    //Check for class name of Notifier Snack Bar. Not present now as we are over the AutoHideDuration.
    expect(wrapper.find('.Notifier-root-120').hostNodes().length).toEqual(0);
  });

});
