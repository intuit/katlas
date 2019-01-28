import * as types from './actionTypes';
import * as notifyActions from './notifyActions';
//Use of app history here so that route navigations can be occur in actions,
//which are otherwise not wrappable withRouter
import history from '../history';
import ApiService from '../services/ApiService';
import { QUERY_LEN_ERR } from '../utils/errors';
import { getQSLObjTypes } from '../utils/validate';

//TODO:DM - is there a better place to define router related consts?
const APP_RESULTS_ROUTE = '/results?query=';

//QSL requests will always include this telltale character
//TODO:SS - check with @kianjones4 to see if this is the best strategy and char to use. possibly square brackets are an even better choice? what's least likely to occur in kube elements which might otherwise end up in a naive keyword search?
const QSL_TAG = '{';

export function submitQuery(query) {
  return dispatch => {
    if (query !== '' && query.length >= 3) {
      history.push(APP_RESULTS_ROUTE + encodeURIComponent(query));
    } else {
      dispatch(notifyActions.showNotify(QUERY_LEN_ERR));
    }
  };
}

export const requestQuery = (queryStr, page, rowsPerPage) => ({
  type: types.REQUEST_QUERY,
  queryStr,
  page,
  rowsPerPage
});

export function fetchQuery(query, page, rowsPerPage) {
  return dispatch => {
    dispatch(requestQuery(query, page, rowsPerPage));
    let requestPromise;

    if (query.includes(QSL_TAG)) {
      const objTypes = getQSLObjTypes(query);
      // cache metadata
      objTypes.forEach(objType => dispatch(fetchMetadata(objType)));

      requestPromise = ApiService.getQSLResult(query, page, rowsPerPage);
    } else {
      requestPromise = ApiService.getQueryResult(query);
    }

    return requestPromise
      .then(handleResponse)
      .then(([results, count]) => dispatch(receiveQuery(results, count)));
  };
}

export const receiveQuery = (results, count) => ({
  type: types.RECEIVE_QUERY,
  results,
  count
});

export const updatePagination = (page, rowsPerPage) => {
  return (dispatch, getState) => {
    dispatch(fetchQuery(getState().query.current, page, rowsPerPage));
  };
};

function handleResponse(json) {
  let results = [];
  let existingUids = {};
  let count = 0;
  for (let objKey in json) {
    let objArr = json[objKey];
    if (objKey === 'count') {
      count = objArr;
      continue;
    }

    if (objArr.length) {
      objArr.forEach(obj => {
        //screen out duplicate UID entries
        if (!existingUids[obj.uid]) {
          results.push(obj);
          existingUids[obj.uid] = true;
        }
      });
    }
  }

  // TODO, as search query does not support pagination, there is no count returned.
  // will remove this once it supports pagination
  if (count === 0) {
    count = results.length;
  }
  return [results, count];
}

export const requestMetadata = objType => ({
  type: types.REQUEST_METADATA,
  objType
});

export const receiveMetadata = (objType, metadata) => ({
  type: types.RECEIVE_METADATA,
  objType,
  metadata
});

export function fetchMetadata(objType) {
  return (dispatch, getState) => {
    // only fetch metadata not cached before
    if (objType in getState().query.metadata) {
      return;
    }
    dispatch(requestMetadata(objType));
    ApiService.getMetadata(objType).then(json =>
      dispatch(receiveMetadata(objType, json))
    );
  };
}
