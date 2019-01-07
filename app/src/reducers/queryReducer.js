import initialState from './initialState';
import {CHANGE_QUERY, SUBMIT_QUERY, RECEIVE_QUERY} from '../actions/actionTypes';

export default function query(state = initialState.query, action) {
  let newState;
  switch (action.type) {
    case CHANGE_QUERY:
      newState = {
        current: action.query,
        lastSubmitted: state.lastSubmitted,
        submitted: false,
        isWaiting: state.isWaiting,
        results: state.results,
      };
      return newState;
    case SUBMIT_QUERY:
      newState = {
        ...state,
        submitted: true,
        isWaiting: true,
        results: [],
      }
      return newState;
    case RECEIVE_QUERY:
      newState = {
        current: state.current,
        lastSubmitted: state.lastSubmitted,
        submitted: state.submitted,
        isWaiting: false,
        results: action.results
      };
      return newState;
    default:
      return state;
  }
}