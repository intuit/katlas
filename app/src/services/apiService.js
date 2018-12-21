import {HttpService} from './httpService';

  //Used by Callers of the apiService
export const QUERY_KEYWORD_SERVICE_PATH = '/query';
export const QUERY_QSL_SERVICE_PATH = '/qsl';
export const QUERY_PARAM_NAME = 'keyword';
export const QSL_PARAM_NAME = 'qslstring';

const ALL_SERVICE_CONTEXT = '/v1';

const ENTITY_SERVICE_PATH = '/entity/uid/';

export class apiService {

  static getQueryResult(querySvcPath, paramName, paramValue) {
    const params = {
      [paramName]: paramValue
    };
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return request_helper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      querySvcPath, params);
  }

  static getEntity(uid) {
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return request_helper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      ENTITY_SERVICE_PATH + uid);
  }
}

const request_helper = (url, params) => {
  return HttpService.get({
    url,
    params
  }).then((response) => {
    return response;
  }).catch((error) => {
    throw error;
  });
};
