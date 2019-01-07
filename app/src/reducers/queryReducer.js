import initialState from './initialState';
import {CHANGE_QUERY, SUBMIT_QUERY, FETCH_QUERY, RECEIVE_QUERY} from '../actions/actionTypes';

export default function query(state = initialState.query, action) {
  let newState;
  switch (action.type) {
    case CHANGE_QUERY:
      newState = {
        ...state,
        current: action.query,
        submitted: false,
      };
      return newState;
    case SUBMIT_QUERY:
      newState = {
        ...state,
        submitted: true,
        isWaiting: true,
        results: [], //new array to clear out old results upon new submission
      };
      return newState;
    case FETCH_QUERY:
      return action;
    case RECEIVE_QUERY:
      newState = {
        ...state,
        isWaiting: false,
        results: action.results
      };
      return newState;
    default:
      return state;
  }
}