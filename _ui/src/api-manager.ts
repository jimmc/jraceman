export interface XhrOptions {
  method?: string;
  params?: any;
}

export interface LoginStateEvent {
  State: boolean;
  Permissions: string;
}

// ApiManager provides assistance for calling our back end services.
export class ApiManager {
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
            ApiManager.AnnounceLoginStatus()
          } else {
            reject(request);
          }
        }
      };
      const method = (options && options.method) || "GET";
      request.open(method, url);
      const params = (options && options.params) || {};
      /*
      if (params) {
        request.setRequestHeader("Content-Type", "application/json");
        request.send(JSON.stringify(params));
      }
      */
      request.send(params)
    })
  }

  // AnnounceLoginStatus gets our current login status and dispatches a LoginStateEvent.
  public static async AnnounceLoginStatus() {
    try {
      const statusUrl = "/auth/status/"
      const response = await ApiManager.xhrJson(statusUrl)
      const loggedIn = response.LoggedIn
      const permissions = response.Permissions
      ApiManager.SendLoginStateEvent(loggedIn, permissions)
    } catch (e) {
      console.error("auth status call failed")
    }
  }

  public static SendLoginStateEvent(loggedIn: boolean, permissions: string) {
    // Dispatch an event so others can take action when the login state changes.
    const event = new CustomEvent<LoginStateEvent>('jraceman-login-state-event', {
      bubbles: true,
      composed: true,
      detail: {
        State: loggedIn,
        Permissions: permissions
      } as LoginStateEvent
    })
    document.dispatchEvent(event)
  }

}
