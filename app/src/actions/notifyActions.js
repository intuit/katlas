import * as types from './actionTypes';
import { openSnackbar } from '../notifier/Notifier';

export function notify(str) {
    showNotify(str);
    openSnackbar({ message: str });
}

export const showNotify = str => ({
  type: types.SHOW_NOTIFY,
  msg: str
});
