import * as types from './actionTypes';
import {apiService} from "../services/apiService";
import * as notifyActions from './notifyActions';
import history from '../history';

const QUERY_KEYWORD_SERVICE_PATH = '/query';
const QUERY_QSL_SERVICE_PATH = '/qsl';
const QUERY_PARAM_NAME = 'keyword';
const QSL_PARAM_NAME = 'qslstring';

export const changeQuery = str => ({
  type: types.CHANGE_QUERY,
  query: str
});

export function submitQuery(query) {
  return dispatch => {
    if(query !== '' && query.length >= 3) {
      history.push('/results?query=' + encodeURIComponent(query));
    } else {
      const msg = 'Minimum length of Search word must be 3 characters.';
      dispatch(notifyActions.showNotify(msg));
      return;
    }
  };
}

export function fetchQuery(query) {
  return dispatch => {

    let requestPromise;

    if (query.includes('@')) {
      requestPromise = apiService.getQueryResult(QUERY_QSL_SERVICE_PATH, QSL_PARAM_NAME, query);
    } else {
      requestPromise = apiService.getQueryResult(QUERY_KEYWORD_SERVICE_PATH, QUERY_PARAM_NAME, query);
    }

    return requestPromise
    .then(handleResponse)
    .then(json => dispatch(receiveQuery(json)));
  };
}

export const receiveQuery = json => ({
  type: types.RECEIVE_QUERY,
  results: json
});

function handleResponse(json) {
    let results = [];
    let existingUids = {};
    for (let objKey in json){
      let objArr = json[objKey];
      if (objArr.length) {
        objArr.forEach(obj => {
          //screen out duplicate UID entries
          if(!existingUids[obj.uid]){
            results.push(obj);
            existingUids[obj.uid] = true;
          }
        });
      }
    }
    return results;
  }