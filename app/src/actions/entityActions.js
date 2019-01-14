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
      //TODO:DM - can this hard-coded idx be cleaned up?
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

export const fetchQslResp = str => {
  return dispatch => {
    return ApiService.getQSLResult(str)
      .then(json => dispatch(receiveQslResp(json.objects[56])));
  };
};

export const receiveQslResp = results => ({
  type: types.RECEIVE_QSL_RESP,
  results
});