import {combineReducers} from 'redux';
import query from './queryReducer';
import entity from './entityReducer';
import notify from './notifyReducer';

const rootReducer = combineReducers({
  query,
  entity,
  notify
});

export default rootReducer;