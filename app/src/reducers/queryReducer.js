import initialState from './initialState';
import {
  REQUEST_QUERY,
  SUBMIT_QUERY,
  RECEIVE_QUERY
} from '../actions/actionTypes';

export default function query(state = initialState.query, action) {
  let newState;
  switch (action.type) {
    case REQUEST_QUERY:
      newState = {
        ...state,
        current: action.queryStr,
        page: action.page,
        rowsPerPage: action.rowsPerPage,
        isWaiting: true
      };
      return newState;
    case SUBMIT_QUERY:
      newState = {
        ...state,
        results: []
      };
      return newState;
    case RECEIVE_QUERY:
      newState = {
        ...state,
        results: action.results,
        count: action.count,
        isWaiting: false
      };
      return newState;
    default:
      return state;
  }
}
