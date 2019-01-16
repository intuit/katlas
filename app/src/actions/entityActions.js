import * as types from './actionTypes';
import ApiService from "../services/ApiService";

export const setRootUid = uid => ({
  type: types.SET_ROOT_UID,
  uid
});

export const addWatchUid = uid => ({
  type: types.ADD_WATCH_UID,
  uid
});

export const fetchEntity = uid => {
  return dispatch => {
    return ApiService.getEntity(uid)
      //idx 0 used because Entity API is known to return 1 object, wrapped in an array
      .then(json => dispatch(receiveEntity(json.objects[0])));
  };
};

export const fetchEntities = uids => {
  return dispatch => {
    uids.map(uid => dispatch(fetchEntity(uid)));
  };
};

export const receiveEntity = results => ({
  type: types.RECEIVE_ENTITY,
  results
});

export const addWatchQslQuery = query => ({
  type: types.ADD_WATCH_QSL_QUERY,
  query
});

export const fetchQslQuery = query => {
  return dispatch => {
    return ApiService.getQSLResult(query)
      .then(json => dispatch(receiveQslResp(json.objects)));
  };
};

export const receiveQslResp = results => ({
  type: types.RECEIVE_QSL_RESP,
  results
});