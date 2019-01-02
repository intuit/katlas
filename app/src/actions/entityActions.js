import * as types from './actionTypes';
import ApiService from "../services/ApiService";

export const setRootEntity = uid => ({
  type: types.SET_ROOT_ENTITY,
  uid
});

export const addEntityWatch = uid => ({
  type: types.ADD_ENTITY_WATCH,
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
    return Promise.all(uids.map(uid =>
      dispatch(fetchEntity(uid))
    ));
  };
};

export const receiveEntity = json => ({
  type: types.RECEIVE_ENTITY,
  results: json
});