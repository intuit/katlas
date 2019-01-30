import initialState from './initialState';
import {
  REQUEST_QUERY,
  SUBMIT_QUERY,
  RECEIVE_QUERY,
  RECEIVE_METADATA
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
        isWaiting: true,
        isQSL: action.isQSL
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
    case RECEIVE_METADATA:
      const newMetadata = {
        ...state.metadata,
        [action.objType]: action.metadata
      };
      return {
        ...state,
        metadata: newMetadata
      };
    default:
      return state;
  }
}
