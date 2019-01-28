export class HttpService {
  static get({ url, params, arrayParams }) {
    if (params) {
      url = url + '?';
      let paramCnt = 1;

      Object.keys(params).forEach(key => {
        const paramValue = params[key];
        if (paramValue !== undefined) {
          url = url + key + '=' + params[key];
          if (paramCnt < Object.keys(params).length) {
            url = url + '&';
          }
          paramCnt++;
        }
      });
    }

    if (arrayParams) {
      if (!params) {
        url = url + '?';
      }

      for (let i = 0; i < arrayParams.length; i++) {
        //TODO:DM - clean up fn creation in loop
        Object.keys(arrayParams[i]).forEach(key => {
          arrayParams[i][key].forEach(e => (url = url + key + '=' + e + '&'));
        });
      }
    }

    return fetch(url);
  }
}
