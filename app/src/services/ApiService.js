import {HttpService} from './HttpService';

const ALL_SERVICE_CONTEXT = '/v1';
const QUERY_SERVICE_PATH = '/query';
const ENTITY_SERVICE_PATH = '/entity/uid/';

export default class ApiService {
  static getKeyword(input) {
    const params = {
      keyword: input
    };
    //load env provided URL at query time to allow conf.js to load it in time
    //in testing
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return requestHelper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      QUERY_SERVICE_PATH, params);
  }

  static getEntity(uid) {
    //load env provided URL at query time to allow conf.js to load it in time
    //in testing
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return requestHelper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      ENTITY_SERVICE_PATH + uid);
  }
}

const requestHelper = (url, params) => {
  return HttpService.get({
    url,
    params
  }).then((response) => {
    return response;
  }).catch((error) => {
    throw error;
  });
};
