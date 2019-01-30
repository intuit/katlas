import * as types from './actionTypes';
import * as notifyActions from './notifyActions';
//Use of app history here so that route navigations can be occur in actions,
//which are otherwise not wrappable withRouter
import history from '../history';
import ApiService from "../services/ApiService";
import { QUERY_LEN_ERR } from '../utils/errors';
import { validateQslQuery } from '../utils/validate';

//TODO:DM - is there a better place to define router related consts?
const APP_RESULTS_ROUTE = '/results?query=';
const GRAPH_ROUTE = '/graph/';

//QSL requests will always include this telltale character
//TODO:SS - check with @kianjones4 to see if this is the best strategy and char to use. possibly square brackets are an even better choice? what's least likely to occur in kube elements which might otherwise end up in a naive keyword search?
const QSL_TAG = '@';

export const changeQuery = str => ({
  type: types.CHANGE_QUERY,
  query: str
});

const submitQueryAction = () => ({
  type: types.SUBMIT_QUERY,
});

export function submitQuery(query) {
  return dispatch => {
    if(query !== '' && query.length >= 3) {
      dispatch(submitQueryAction());
      history.push(APP_RESULTS_ROUTE + encodeURIComponent(query));
    } else {
      dispatch(notifyActions.showNotify(QUERY_LEN_ERR));
    }
  };
}

export function submitQslQuery(query) {
  return dispatch => {
    //currently, we'll navigate to the graph without any addl validation steps
    if (validateQslQuery(query)) {
      history.push(GRAPH_ROUTE + query);
    }
  }
}

export function fetchQuery(query) {
  return dispatch => {
    let requestPromise;

    if (query.includes(QSL_TAG)) {
      requestPromise = ApiService.getQSLResult(query);
    } else {
      requestPromise = ApiService.getQueryResult(query);
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
