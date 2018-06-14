/// <reference path="../../components/api-manager/api-manager.ts" />

describe('api-manager', () => {
  let testFixture: ApiManager;

  let requests: sinon.SinonFakeXMLHttpRequest[];

  beforeEach(() => {
    const xhrFake = sinon.useFakeXMLHttpRequest();
    requests = [];
    xhrFake.onCreate = function(request: sinon.SinonFakeXMLHttpRequest):void {
      requests.push(request);
    };
  });

  it('resolves a success response', (done: ()=>void) => {
    const promise = ApiManager.xhr("http://localhost").then((req: XMLHttpRequest) => {
      chai.expect(requests[0].url).to.equal("http://localhost");
      chai.expect(req.status).to.equal(200);
      chai.expect(req.responseText).to.equal("test body");
      done();
    });
    chai.expect(requests.length).to.equal(1);
    const headers = { "content-type": "application/text" };
    const body = 'test body';
    requests[0].respond(200, headers, body);
  });

  it('rejects an error response', (done: ()=>void) => {
    const promise = ApiManager.xhr("http://localhost").catch((req: XMLHttpRequest) => {
      chai.expect(requests[0].url).to.equal("http://localhost");
      chai.expect(req.status).to.equal(400);
      chai.expect(req.responseText).to.equal("error body");
      done();
    });
    chai.expect(requests.length).to.equal(1);
    const headers = { "content-type": "application/text" };
    const body = 'error body';
    requests[0].respond(400, headers, body);
  });
  
  it('returns body text from xhrText', (done: ()=>void) => {
    const promise = ApiManager.xhrText("http://localhost").then((text: string) => {
      chai.expect(requests[0].url).to.equal("http://localhost");
      chai.expect(text).to.equal("test body");
      done();
    });
    chai.expect(requests.length).to.equal(1);
    const headers = { "content-type": "application/text" };
    const body = 'test body';
    requests[0].respond(200, headers, body);
  });

  it('returns JSON from xhrJson', (done: ()=>void) => {
    const promise = ApiManager.xhrJson("http://localhost").then((data: any) => {
      chai.expect(requests[0].url).to.equal("http://localhost");
      chai.expect(data).to.deep.equal({'a': 1, 'b': 'x'});
      done();
    });
    chai.expect(requests.length).to.equal(1);
    const headers = { "content-type": "application/json" };
    const body = '{"a": 1, "b": "x"}';
    requests[0].respond(200, headers, body);
  });

  it('returns null from xhrJson for empty body text', (done: ()=>void) => {
    const promise = ApiManager.xhrJson("http://localhost").then((data: any) => {
      chai.expect(requests[0].url).to.equal("http://localhost");
      chai.expect(data).to.equal(null);
      done();
    });
    chai.expect(requests.length).to.equal(1);
    const headers = { "content-type": "application/json" };
    const body = '';
    requests[0].respond(200, headers, body);
  });

});
