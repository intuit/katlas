import * as types from './actionTypes';

export const showNotify = str => {
  //alert('In showNotify action ' + str);
  return {
    type: types.SHOW_NOTIFY,
    msg: str
  };
};