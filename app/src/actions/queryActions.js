import * as types from './actionTypes';
import {ApiService} from "../services/apiService";

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
    return ApiService.getKeyword(query)
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