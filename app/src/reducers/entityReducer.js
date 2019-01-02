import _ from 'lodash';

import initialState from './initialState';
import { SET_ROOT_ENTITY, ADD_ENTITY_WATCH, FETCH_ENTITY,
  FETCH_ENTITIES, RECEIVE_ENTITY} from '../actions/actionTypes';
import { EdgeLabels } from '../config/appConfig';

export default function entity(state = initialState.entity, action) {
  let newState, now, potentialResults;
  switch (action.type) {
    case SET_ROOT_ENTITY:
      newState = {
        ...state, //start with a copy of existing state
        rootUid: action.uid, //and apply changes on top of existing state attrs
      };
      return newState;
    case ADD_ENTITY_WATCH:
      newState = {
        ...state,
        isWaiting: true,
      };
      //extend the entity obj separately since we only want to change part of it
      newState.entitiesByUid[action.uid] = {};
      return newState;
    case FETCH_ENTITY:
      return action;
    case FETCH_ENTITIES:
      return action;
    case RECEIVE_ENTITY:
      now = +new Date();
      newState = {
        ...state,
        latestTimestamp: now,
        isWaiting: false,
      };
      //extend the entity obj separately since we only want to change part of it
      newState.entitiesByUid[action.results.uid] = action.results;
      //build a new result obj but only update newState if it's different
      potentialResults = entityWalk(newState.rootUid, newState.entitiesByUid);
      if (!_.isEqual(state.results, potentialResults)) {
        newState.results = potentialResults;
      }
      return newState;
    default:
      return state;
  }
}

const entityWalk = (rootUid, entityObj) => {
  // start with root obj
  let results = entityObj[rootUid];
  //walk it (recursing into all arrs)
  //TODO:DM - this may not actually be recursing as deep as I want... more than 1 hop from root
  for (let key in results) {
    let candidate = results[key];
    //ensure that the key is an expected relationship and dealing with an array
    if ((EdgeLabels.indexOf(key) > -1) && _.isArray(candidate)) {
      for (let idx in candidate) {
        let node = candidate[idx];
        if (entityObj[node.uid]) {
          _.merge(node, entityObj[node.uid]);
        }
      }
    }
  }
  return results;
};