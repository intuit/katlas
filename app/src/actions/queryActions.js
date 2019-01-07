import * as types from './actionTypes';
import {apiService} from "../services/apiService";
import * as notifyActions from './notifyActions';
import history from '../history';
import {QUERY_KEYWORD_SERVICE_PATH, QUERY_QSL_SERVICE_PATH, QUERY_PARAM_NAME, QSL_PARAM_NAME} from "../services/apiService";
import {QUERY_LEN_ERR} from "../utils/errors";

//QSL uses this symbol, which can be used as a Query type differentiator.
const QSL_TAG = '@';

export const changeQuery = str => ({
  type: types.CHANGE_QUERY,
  query: str
});

const submitQueryAction = () => ({
  type: types.SUBMIT_QUERY,
})

export function submitQuery(query) {
  return dispatch => {
    if(query !== '' && query.length >= 3) {
      dispatch(submitQueryAction());
      history.push('/results?query=' + encodeURIComponent(query));
    } else {
      dispatch(notifyActions.showNotify(QUERY_LEN_ERR));
      return;
    }
  };
}

export function fetchQuery(query) {
  return dispatch => {

    let requestPromise;

    if (query.includes(QSL_TAG)) {
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
