import * as notifyActions from '../actions/notifyActions';

export class HttpService  {

  static get({url, params, arrayParams}) {
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
                arrayParams[i][key].forEach(e => url = url + key + "=" + e + "&");
            });
        }
    }

    return fetch(url)
        .then(res => this.makeResponse(res))
  }

  static makeResponse(res) {
    if (res.status === 204) {
      return null;
    }
    if (res.ok) {
      return res.json();
    } else {
      notifyActions.notify(res.statusText);
      console.error('Error from Rest Service ' + res.statusText + ',' + res.status);
      return null;
    }
  }
}
