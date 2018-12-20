import * as types from './actionTypes';
import {ApiService} from "../services/ApiService";
import * as notifyActions from '../actions/notifyActions';
import history from '../history';

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

export const receiveQueryResp = json => ({
  type: types.RECEIVE_QUERY,
  results: json
});

export function fetchQuery(query) {
  return dispatch => {

    let requestPromise;

    if (query.includes('@')) {
      requestPromise = ApiService.getQueryResult('/qsl', 'qslstring', query);
    } else {
      requestPromise = ApiService.getQueryResult('/query', 'keyword', query);
    }

    return requestPromise
    .then(handleResponse)
    .then(json => dispatch(receiveQueryResp(json)));
  };
}

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