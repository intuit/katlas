import {combineReducers} from 'redux';
import query from './queryReducer';
import entity from './entityReducer';

const rootReducer = combineReducers({
  query,
  entity
});

export default rootReducer;