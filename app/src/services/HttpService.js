import * as notifyActions from "../actions/notifyActions";
import store from "../store.js";

export class HttpService {
  static get({ url, params, arrayParams }) {
    if (params) {
      url = url + "?";
      let paramCnt = 1;

      Object.keys(params).forEach(key => {
        const paramValue = params[key];
        if (paramValue !== undefined) {
          url = url + key + "=" + params[key];
          if (paramCnt < Object.keys(params).length) {
            url = url + "&";
          }
          paramCnt++;
        }
      });
    }

    if (arrayParams) {
      if (!params) {
        url = url + "?";
      }

      for (let i = 0; i < arrayParams.length; i++) {
        //TODO:DM - clean up fn creation in loop
        Object.keys(arrayParams[i]).forEach(key => {
          arrayParams[i][key].forEach(e => (url = url + key + "=" + e + "&"));
        });
      }
    }

    return fetch(url).then(res => this.makeResponse(res));
  }

  static makeResponse(res) {
    if (res.status === 204) {
      return null;
    }
    if (res.ok) {
      return res.json();
    } else {
      res.text().then(txt => store.dispatch(notifyActions.showNotify(txt)));
      return null;
    }
  }
}
