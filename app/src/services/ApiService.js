import { HttpService } from './HttpService';

const ALL_SERVICE_CONTEXT = '/v1';

const QUERY_KEYWORD_SERVICE_PATH = '/query';
const QUERY_KEYWORD_PARAM_NAME = 'keyword';
const QUERY_QSL_SERVICE_PATH = '/qsl';
const QUERY_QSL_PARAM_NAME = 'qslstring';

const ENTITY_SERVICE_PATH = '/entity/uid/';

export default class ApiService {
  static getQueryResult(query) {
    const params = {
      [QUERY_KEYWORD_PARAM_NAME]: query
    };
    //load env provided URL at query time to allow conf.js to load it in time
    //in testing
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return requestHelper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      QUERY_KEYWORD_SERVICE_PATH, params);
  }

  static getQSLResult(query) {
    const params = {
      [QUERY_QSL_PARAM_NAME]: query
    };
    //load env provided URL at query time to allow conf.js to load it in time
    //in testing
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return requestHelper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      QUERY_QSL_SERVICE_PATH, params);
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
