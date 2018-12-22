import * as types from './actionTypes';
import ApiService from "../services/ApiService";

export const addEntityWatch = uid => ({
  type: types.ADD_ENTITY_WATCH,
  uid
});

export const fetchEntity = uid => {
  return dispatch => {
    return ApiService.getEntity(uid)
    .then(handleResponse)
    .then(json => dispatch(receiveEntity(json[0])));//index into JSON arr to grab single entity object
  };
};

export const receiveEntity = json => ({
  type: types.RECEIVE_ENTITY,
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

//only update state if the objects fail lodash equality check AND
//     //the component is still mounted. usually, the lifecycle methods should
//     //be used directly for such things that, but in testing we're getting
//     //intermittent errors that setState is being called on unmounted
//     //components, without this check
//     if(!_.isEqual(this.state.data, json) && this._isMounted) {
//       this.setState({
//         data: json,
//         waitingOnReq: false
//       });

//logic here to transform API raw data to something usable by graph, or does that make more sense in actions?

//also trigger periodic re-requests of data here with knowledge of period, max reqs/sec, num objs tracked?
//no, I think it still makes sense to do that in a component, probably GraphContainer, where we easily have access to
//store to see num of requests needed and can compute interval ms based on limit and num entities
