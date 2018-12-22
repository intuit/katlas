import initialState from './initialState';
import {ADD_ENTITY_WATCH, FETCH_ENTITY, RECEIVE_ENTITY} from '../actions/actionTypes';

export default function entity(state = initialState.entity, action) {
  let newState, now;
  switch (action.type) {
    case ADD_ENTITY_WATCH:
      //start with a copy of existing state
      newState = { ...state };
      //extend the uidsObj which is otherwise difficult to with variable key
      //name action.uid in statement above
      newState.uidsObj[action.uid] = true;
      return newState;
    case FETCH_ENTITY:
      return action;
    case RECEIVE_ENTITY:
      now = +new Date();
      //project changes on top of existing state attrs
      newState = { ...state,
        //results: action.results, // how to add to existing results?
        latestTimestamp: now,
      };
      //TODO:DM - technically I could use just 1 obj to maintain UIDs to watch as well as results, just make null until new obj arrives {val: null, timestamp: ...}
      newState.results[action.results.uid] = action.results;
      return newState;
    default:
      return state;
  }
}

