import {HttpService} from "./HttpService";

const SERVICE_CONTEXT = "/v1";

export class ApiService {
  //TODO:DM - determine how to best re-expand use of the entity specific routes of QueryAPI (ex: /query/k8sNamespace, /query/k8sCluster), but they aren't in use currently and as-written, they weren't nearly DRY enough

  static getQueryResult(path, paramName, paramValue) {
    const params = {
      [paramName]: paramValue
    };
    const servicesURL = window.envConfig.KATLAS_API_URL;

    return HttpService.get({
      url: servicesURL + SERVICE_CONTEXT + path,
      params: params,
    }).then((response) => {
      return response;
    }).catch((error) => {
      throw error;
    });

  }
}
