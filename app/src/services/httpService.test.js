import fetchMock from 'fetch-mock';
import {HttpService} from './httpService';

it('Returns 200 for GET request, then response is as expected', () => {

    let dummyUrl = "http://katlas.com/v1/qsl";
    let dummyParams = {qslstring: 'Cluster'};
    const response = {
      body: {foo: 'bar'},
      status: 200
    };

    fetchMock.get('*', response, { overwriteRoutes: true });

    return HttpService.get({
      url: dummyUrl,
      params: dummyParams
    }).then((response) => {
      expect(fetchMock.called).toBeTruthy();
      expect(response).toEqual(response);
    });

});

it('Returns 204 for GET request, then response is null', () => {

    let dummyUrl = "http://katlas.com/v1/qsl";
    let dummyParams = {qslstring: 'Cluster'};
    const response = {
      body: {foo: 'bar'},
      status: 204
    };

    fetchMock.get('*', response, { overwriteRoutes: true });

    return HttpService.get({
      url: dummyUrl,
      params: dummyParams
    }).then((response) => {
      expect(fetchMock.called).toBeTruthy();
      expect(response).toEqual(null);
    });

});

it('Returns 400 for GET request, then response is null', () => {

    let dummyUrl = "http://katlas.com/v1/qsl";
    let dummyParams = {qslstring: 'Cluster'};
    const response = {
      body: {},
      status: 400
    };

    fetchMock.get('*', response, { overwriteRoutes: true });

    return HttpService.get({
      url: dummyUrl,
      params: dummyParams
    }).then((response) => {
      expect(fetchMock.called).toBeTruthy();
      expect(response).toEqual(null);
    });

});

