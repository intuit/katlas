import * as types from './actionTypes';

export const showNotify = str => ({
    type: types.SHOW_NOTIFY,
    msg: str
});