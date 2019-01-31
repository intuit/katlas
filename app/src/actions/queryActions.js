import * as types from './actionTypes';
import * as notifyActions from './notifyActions';
//Use of app history here so that route navigations can be occur in actions,
//which are otherwise not wrappable withRouter
import history from '../history';
import ApiService from '../services/ApiService';
import { QUERY_LEN_ERR } from '../utils/errors';
import { validateQslQuery, getQSLObjTypes } from '../utils/validate';

//TODO:DM - is there a better place to define router related consts?
const APP_RESULTS_ROUTE = '/results?query=';
const GRAPH_ROUTE = '/graph/';

export function submitQuery(query) {
  return dispatch => {
    if (query !== '' && query.length >= 3) {
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

    if (validateQslQuery(query)) {
      dispatch(requestQuery(query, true, page, rowsPerPage));
      const objTypes = getQSLObjTypes(query);
      // cache metadata
      objTypes.forEach(objType => dispatch(fetchMetadata(objType)));

      requestPromise = ApiService.getQSLResult(query, page, rowsPerPage);
    } else {
      dispatch(requestQuery(query, false, page, rowsPerPage));
      requestPromise = ApiService.getQueryResult(query, page, rowsPerPage);
    }

    return requestPromise.then(json => {
      if (json != null) {
        dispatch(receiveQuery(json.objects, json.count));
      }
    });
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
