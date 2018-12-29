import fetchMock from 'fetch-mock';
import {HttpService} from './httpService';

it('Test get method for response 200', () => {

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

it('Test get method for response 204', () => {

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

it('Test get method for error response', () => {

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

