import initialState from './initialState';
import {SHOW_NOTIFY} from '../actions/actionTypes';

export default function notify(state = initialState.notify, action) {
  let newState;
  switch (action.type) {
    case SHOW_NOTIFY:
      newState = {
        msg: action.msg,
        timestamp: +new Date()
      };
      return newState;
    default:
      return state;
  }
}