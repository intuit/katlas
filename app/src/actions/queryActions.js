import * as types from './actionTypes';
import {ApiService} from "../services/ApiService";
import * as notifyActions from '../actions/notifyActions';

export const changeQuery = str => ({
  type: types.CHANGE_QUERY,
  query: str
});

export const submitQuery = () => ({type: types.SUBMIT_QUERY});

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
      if(query.length < 3) {
        const msg = 'Minimum length of Search word must be 3 characters.';
        notifyActions.notify(msg);
        return;
      }
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