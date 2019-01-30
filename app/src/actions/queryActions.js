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

export const requestQuery = (queryStr, isQSL, page, rowsPerPage) => ({
  type: types.REQUEST_QUERY,
  queryStr,
  isQSL,
  page,
  rowsPerPage
});

export function fetchQuery(query, page, rowsPerPage) {
  return dispatch => {
    let requestPromise;

    if (query.includes(QSL_TAG)) {
      dispatch(requestQuery(query, true, page, rowsPerPage));
      const objTypes = getQSLObjTypes(query);
      // cache metadata
      objTypes.forEach(objType => dispatch(fetchMetadata(objType)));

      requestPromise = ApiService.getQSLResult(query, page, rowsPerPage);
    } else {
      dispatch(requestQuery(query, true, page, rowsPerPage));
      requestPromise = ApiService.getQueryResult(query, page, rowsPerPage);
    }

    return requestPromise.then(json => dispatch(receiveQuery(json)));
  };
}

export const receiveQuery = json => ({
  type: types.RECEIVE_QUERY,
  json
});

export const updatePagination = (page, rowsPerPage) => {
  return (dispatch, getState) => {
    dispatch(fetchQuery(getState().query.current, page, rowsPerPage));
  };
};

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
