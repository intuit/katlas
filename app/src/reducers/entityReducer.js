import _ from 'lodash';

import initialState from './initialState';
import {
  SET_ROOT_UID, ADD_WATCH_UID, FETCH_ENTITY,
  FETCH_ENTITIES, RECEIVE_ENTITY, RECEIVE_QSL_RESP, ADD_WATCH_QSL_QUERY
} from '../actions/actionTypes';
import { EdgeLabels } from '../config/appConfig';

export default function entity(state = initialState.entity, action) {
  let newState, now, potentialResults;
  switch (action.type) {
    case SET_ROOT_UID:
      newState = {
        ...state, //start with a copy of existing state
        rootUid: action.uid, //and apply changes on top of existing state attrs
      };
      return newState;
    case ADD_WATCH_UID:
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
    case ADD_WATCH_QSL_QUERY:
      newState = {
        ...state,
        isWaiting: true,
        qslQuery: action.query
      };
      return newState;
    case RECEIVE_QSL_RESP:
      now = +new Date();
      newState = {
        ...state,
        latestTimestamp: now,
        isWaiting: false,
        results: action.results
      };
      //TODO:DM - decice if there is still any value to this OR we'll just write response data directly into store as 'results'
      // if(action.results){
      //   let entityObj = action.results
      //     .reduce((results, obj) => {
      //       results[obj.uid] = obj;
      //       return results;
      //     }, {});
      //   newState.results = action.results
      //     .filter(item => item.uid)
      //     .map(item => entityWalk(item.uid, entityObj))
      //     .reduce((results, obj) => {
      //       debugger;
      //       results[obj.uid] = obj;
      //       return results;
      //     }, {});
      // }
      return newState;
    default:
      return state;
  }
}

const entityWalk = (rootUid, entityObj) => {
  // start with root obj
  let results = entityObj[rootUid];
  let encounteredUids = {};
  encounteredUids[rootUid] = true;

  entityWalkHelper(results, entityObj, encounteredUids);
  return results;
};

const entityWalkHelper = (results, entityObj, encounteredUids) => {
  _.forOwn(results, (childrenCandidate, key) => {
    if ((EdgeLabels.indexOf(key) > -1) && _.isArray(childrenCandidate)){
      _.forEach(childrenCandidate, node => {
        if (node.uid && entityObj[node.uid] && !encounteredUids[node.uid]) {
          _.assign(node, entityObj[node.uid]);
          encounteredUids[node.uid] = true;
        }
        //recurse thru object if children are present, important to do this
        //after this object is augmented so we won't later overwrite
        entityWalkHelper(node, entityObj, encounteredUids);
      });
    }
  });
};