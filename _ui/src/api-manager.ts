export interface XhrOptions {
  method?: string;
  params?: any;
  encoding?: string;    // 'direct' or 'json'; default is 'json'
}

export interface TimeoutSecondsEvent {
  timeoutSeconds: number;      // Number of seconds until our session times out due to inactivity.
}

// ApiManager provides assistance for calling our back end services.
export class ApiManager {
  // We don't want ApiManager to depend on any other modules, but we want
  // to allow someone else to take action when we get an authorization error,
  // as can happen when our authentication has timed out, so we add a
  // callback to allow that, and initialize it to a nop in case we happen
  // to call it before it has been set up.
  static nop(){}
  public static AuthErrorCallback = ApiManager.nop

  public static async xhrJson(url: string, options?: XhrOptions) {
    const response = await this.xhrText(url, options);
    return JSON.parse(response || 'null');
  }

  public static async xhrText(url: string, options?: XhrOptions) {
    const request = await this.xhr(url, options);
    return (request as any).responseText;
  }

  public static xhr(url: string, options?: XhrOptions) {
    const request = new XMLHttpRequest();
    return new Promise((resolve, reject) => {
      request.onreadystatechange = () => {
        if (request.readyState === 4) {
          if (request.status === 200) {
            try {
              resolve(request);
            } catch (e) {
              reject(e);
            }
          } else if (request.status == 401 && request.responseText == "Invalid token\n") {
            reject(request);
            ApiManager.AuthErrorCallback()
          } else {
            reject(request);
          }
          ApiManager.postTimeoutSecondsEvent()
        }
      };
      const method = (options && options.method) || "GET";
      request.open(method, url);
      const encoding = (options && options.encoding) || 'json';
      const params = (options && options.params) || {};
      if (params && encoding=='json') {
        request.setRequestHeader("Content-Type", "application/json");
        request.send(JSON.stringify(params));
      } else {
        request.send(params)
      }
    })
  }

  // postTimeoutSecondsEvent reads the cookie that has our session timeout
  // and dispatches an event with that many seconds as payload.
  static postTimeoutSecondsEvent() {
    const timeoutCookieValue = ApiManager.getCookieValue("JRACEMAN_TOKEN_TIMEOUT")
    const timeoutSeconds = parseInt(timeoutCookieValue)
    const event = new CustomEvent<TimeoutSecondsEvent>('jraceman-timeout-seconds-event', {
      detail: {
        timeoutSeconds: timeoutSeconds
      } as TimeoutSecondsEvent
    });
    // Dispatch the event to the document so any element can listen for it.
    document.dispatchEvent(event);
  }

  // From https://stackoverflow.com/questions/5639346/what-is-the-shortest-function-for-reading-a-cookie-by-name-in-javascript?rq=1
  static getCookieValue(name: string):string {
    return document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')?.pop() || ''
  }
}
